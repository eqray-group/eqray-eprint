package controller

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/pay/service"
)

// PaymentController 提供后台支付流水查询。
type PaymentController struct {
	service *service.PaymentService
}

// NewPaymentController 构造后台支付单控制器。
func NewPaymentController() *PaymentController {
	return &PaymentController{service: service.NewPaymentService()}
}

// Index 分页查询支付单，支持支付单号、会员、渠道与状态筛选。
// 入参：paymentNo、memberId、channel、status、pageNo、pageSize；返回支付单列表与总数。
func (c *PaymentController) Index(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	q := service.AdminPaymentQuery{
		PaymentNo: ctx.Request().Query("paymentNo"),
		Channel:   ctx.Request().Query("channel"),
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

	list, total, err := c.service.AdminPage(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}
