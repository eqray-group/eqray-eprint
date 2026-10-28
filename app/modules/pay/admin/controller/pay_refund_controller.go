package controller

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/pay/admin/request"
	"eqray-eprint/app/modules/pay/service"
)

// RefundController 提供后台退款单的查询与发起退款接口。
type RefundController struct {
	service *service.RefundService
}

// NewRefundController 构造退款控制器。
func NewRefundController() *RefundController {
	return &RefundController{service: service.NewRefundService()}
}

// Index 分页查询退款单，支持按支付单号、会员与状态筛选。
// 入参：paymentNo、memberId、status、pageNo、pageSize；返回退款单列表与总数。
func (c *RefundController) Index(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	q := service.RefundListQuery{
		PaymentNo: ctx.Request().Query("paymentNo"),
		PageNo:    page.PageNo,
		PageSize:  page.PageSize,
	}
	if raw := ctx.Request().Query("memberId"); raw != "" {
		if id, err := strconv.ParseInt(raw, 10, 64); err == nil && id > 0 {
			q.MemberID = types.ID(id)
		}
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := int8(status)
		q.Status = &v
	}

	list, total, err := c.service.List(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Store 后台发起退款：余额支付立即到账，微信支付受理后等回调。
// 入参见 RefundCreateRequest；返回退款单。
func (c *RefundController) Store(ctx http.Context) http.Response {
	var req request.RefundCreateRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if req.PaymentNo == "" {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "支付单号不能为空"))
	}

	refund, err := c.service.Create(ctx.Context(), service.RefundInput{
		PaymentNo:  req.PaymentNo,
		Amount:     req.Amount,
		Reason:     req.Reason,
		OperatorID: types.ID(support.UserID(ctx.Context())),
	})
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, refund)
}
