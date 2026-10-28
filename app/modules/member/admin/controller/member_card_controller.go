package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/member/admin/request"
	"eqray-eprint/app/modules/member/service"
)

// MemberCardController 提供后台会员卡的查询、发放、修改与停用接口。
type MemberCardController struct {
	service *service.MemberCardService
}

// NewMemberCardController 构造会员卡控制器。
func NewMemberCardController() *MemberCardController {
	return &MemberCardController{service: service.NewMemberCardService()}
}

// Index 分页查询会员卡，支持卡号与状态筛选。
// 入参：cardNo、status、pageNo、pageSize（查询串）；返回卡列表与总数。
func (c *MemberCardController) Index(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	q := service.MemberCardListQuery{
		CardNo:   ctx.Request().Query("cardNo"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := int8(status)
		q.Status = &v
	}

	list, total, err := c.service.Page(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询指定会员的会员卡。
// 入参：路径 userID；返回卡详情，未发卡返回业务错误码。
func (c *MemberCardController) Show(ctx http.Context) http.Response {
	memberID, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	card, err := c.service.Detail(ctx.Context(), memberID)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, card)
}

// Save 给会员发卡或修改卡（已有卡时更新，没有卡时新建）。
// 入参：路径 userID + MemberCardSaveRequest；返回会员卡 ID。
func (c *MemberCardController) Save(ctx http.Context) http.Response {
	memberID, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.MemberCardSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	// 会员身份以路径为准，忽略请求体里的 memberUserId，避免越权给其他会员发卡
	id, err := c.service.Save(ctx.Context(), service.MemberCardInput{
		MemberUserID:    memberID,
		CardLevel:       req.CardLevel,
		CardName:        req.CardName,
		DiscountPercent: req.DiscountPercent,
		Status:          req.Status,
		StartTime:       req.StartTime,
		ExpireTime:      req.ExpireTime,
		Remark:          req.Remark,
	})
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Disable 停用会员卡，停用后折扣立即失效。
// 入参：路径 userID；返回 null。
func (c *MemberCardController) Disable(ctx http.Context) http.Response {
	memberID, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Disable(ctx.Context(), memberID); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}
