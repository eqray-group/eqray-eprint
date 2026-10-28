// Package controller 提供 print 模块的后台管理接口。
package controller

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/admin/request"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/service"
)

// ShopController 提供门店与门店成员的管理接口。
type ShopController struct {
	service *service.ShopService
}

// NewShopController 构造门店控制器。
func NewShopController() *ShopController {
	return &ShopController{service: service.NewShopService()}
}

// Index 分页查询门店。
func (c *ShopController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.ShopListQuery{
		Name:     ctx.Request().Query("name"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := int8(status)
		q.Status = &v
	}
	if audit := ctx.Request().QueryInt("auditStatus", -1); audit >= 0 {
		v := int8(audit)
		q.AuditStatus = &v
	}

	list, total, err := c.service.Page(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询门店详情。
func (c *ShopController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	shop, err := c.service.Detail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, shop)
}

// Store 新增门店，审核状态固定为待审核。
func (c *ShopController) Store(ctx http.Context) http.Response {
	var req request.ShopSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	id, err := c.service.Create(ctx.Context(), toShopInput(req))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 编辑门店基础信息。
func (c *ShopController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.ShopSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, toShopInput(req)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Destroy 删除门店。
func (c *ShopController) Destroy(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Delete(ctx.Context(), id, operatorName(ctx)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// ChangeStatus 切换营业状态。
func (c *ShopController) ChangeStatus(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.ShopStatusRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.ChangeStatus(ctx.Context(), id, req.Status, operatorName(ctx)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Audit 审核门店。
func (c *ShopController) Audit(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.ShopAuditRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Audit(ctx.Context(), id, req.AuditStatus, req.Remark, operatorName(ctx)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// StaffIndex 查询门店成员。
func (c *ShopController) StaffIndex(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	list, err := c.service.ListStaff(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, list)
}

// StaffStore 绑定后台账号到门店。
func (c *ShopController) StaffStore(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.StaffSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}
	if req.UserID <= 0 {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "后台账号 ID 不合法"))
	}

	if err := c.service.BindStaff(ctx.Context(), types.ID(req.UserID), id, req.Role, types.ID(support.UserID(ctx.Context()))); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// StaffDestroy 解除门店成员绑定。
func (c *ShopController) StaffDestroy(ctx http.Context) http.Response {
	staffID, err := strconv.ParseInt(ctx.Request().Route("staffId"), 10, 64)
	if err != nil || staffID <= 0 {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "成员 ID 不合法"))
	}

	if err := c.service.UnbindStaff(ctx.Context(), types.ID(staffID)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// toShopInput 把请求体转成服务入参，状态缺省按停业处理。
func toShopInput(req request.ShopSaveRequest) service.ShopInput {
	status := printEnums.ShopStatusClosed
	if req.Status != nil {
		status = *req.Status
	}

	return service.ShopInput{
		Name:          req.Name,
		Address:       req.Address,
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
		Avatar:        req.Avatar,
		BusinessHours: req.BusinessHours,
		Status:        status,
		Remark:        req.Remark,
	}
}
