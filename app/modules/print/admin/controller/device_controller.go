package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/admin/request"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/service"
)

// DeviceController 提供打印设备的后台管理接口。
type DeviceController struct {
	service *service.DeviceService
}

// NewDeviceController 构建设备控制器。
func NewDeviceController() *DeviceController {
	return &DeviceController{service: service.NewDeviceService()}
}

// Index 分页查询设备，支持按门店、关键字、状态与在线状态筛选。
func (c *DeviceController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.DeviceListQuery{
		Keyword:  ctx.Request().Query("keyword"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if raw := ctx.Request().QueryInt("printShopId", 0); raw > 0 {
		q.ShopID = types.ID(raw)
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := int8(status)
		q.Status = &v
	}
	if online := ctx.Request().QueryInt("onlineStatus", -1); online >= 0 {
		v := int8(online)
		q.OnlineStatus = &v
	}

	list, total, err := c.service.Page(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询设备详情。
func (c *DeviceController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	device, err := c.service.Detail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, device)
}

// Store 预录设备（客户端首次登录也会自动登记）。
func (c *DeviceController) Store(ctx http.Context) http.Response {
	var req request.DeviceSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	id, err := c.service.Create(ctx.Context(), req.Code, toDeviceInput(req))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 编辑设备。
func (c *DeviceController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.DeviceSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, toDeviceInput(req)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Destroy 删除设备；其下打印机需先停用。
func (c *DeviceController) Destroy(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Delete(ctx.Context(), id, operatorName(ctx)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// toDeviceInput 把请求体转成服务入参，状态缺省按启用处理。
func toDeviceInput(req request.DeviceSaveRequest) service.DeviceInput {
	status := printEnums.StatusEnabled
	if req.Status != nil {
		status = *req.Status
	}

	return service.DeviceInput{
		PrintShopID: types.ID(req.PrintShopID),
		Name:        req.Name,
		Status:      status,
		Remark:      req.Remark,
	}
}
