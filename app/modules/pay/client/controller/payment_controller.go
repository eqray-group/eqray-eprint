package controller

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/pay/client/request"
	"eqray-eprint/app/modules/pay/service"
)

// PaymentController 提供会员支付单的创建、查询与充值接口。
type PaymentController struct {
	service *service.PaymentService
}

// NewPaymentController 构造支付单控制器。
func NewPaymentController() *PaymentController {
	return &PaymentController{service: service.NewPaymentService()}
}

// Store 创建支付单并返回调起参数；余额支付会直接扣款成功。
// 入参见 CreatePaymentRequest；返回 paymentNo、paid、params / mwebUrl。
func (c *PaymentController) Store(ctx http.Context) http.Response {
	var req request.CreatePaymentRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	bizID, err := strconv.ParseInt(req.BizID, 10, 64)
	if err != nil || bizID <= 0 {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "业务单号不合法"))
	}

	out, err := c.service.Create(ctx.Context(), service.CreatePaymentInput{
		MemberID:    support.MemberID(ctx.Context()),
		BizType:     req.BizType,
		BizID:       types.ID(bizID),
		Amount:      req.Amount,
		ChannelCode: req.ChannelCode,
		Description: req.Description,
		OpenID:      req.OpenID,
		ClientIP:    req.ClientIP,
	})
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, out)
}

// Recharge 创建钱包充值支付单。
// 入参见 RechargeRequest；返回调起参数。
func (c *PaymentController) Recharge(ctx http.Context) http.Response {
	var req request.RechargeRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	out, err := c.service.Recharge(ctx.Context(), support.MemberID(ctx.Context()), req.Amount, req.ChannelCode, req.OpenID, req.ClientIP)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, out)
}

// Show 查询本人支付单。
// 入参：路径 paymentNo；返回支付单详情。
func (c *PaymentController) Show(ctx http.Context) http.Response {
	payment, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()), ctx.Request().Input("paymentNo"))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, payment)
}

// List 分页查询本人支付单。
// 入参：pageNo、pageSize；返回支付单列表与总数。
func (c *PaymentController) List(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	list, total, err := c.service.List(ctx.Context(), support.MemberID(ctx.Context()), page.PageNo, page.PageSize)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Sync 主动查单，用于前端支付后轮询兜底。
// 入参：路径 paymentNo；返回最新支付单。
func (c *PaymentController) Sync(ctx http.Context) http.Response {
	payment, err := c.service.Sync(ctx.Context(), support.MemberID(ctx.Context()), ctx.Request().Input("paymentNo"))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, payment)
}
