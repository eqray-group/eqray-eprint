package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/admin/request"
	"eqray-eprint/app/modules/print/service"
)

// OrderController 提供打印订单的后台查询与核销接口。
type OrderController struct {
	service *service.OrderService
}

// NewOrderController 构造订单控制器。
func NewOrderController() *OrderController {
	return &OrderController{service: service.NewOrderService()}
}

// Index 分页查询订单，支持按门店、订单号、支付状态、订单状态筛选。
func (c *OrderController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.OrderListQuery{
		OrderNo:  ctx.Request().Query("orderNo"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if shopID := ctx.Request().QueryInt("printShopId", 0); shopID > 0 {
		q.ShopID = types.ID(shopID)
	}
	if status := ctx.Request().QueryInt("orderStatus", -1); status >= 0 {
		v := int8(status)
		q.OrderStatus = &v
	}
	if payStatus := ctx.Request().QueryInt("payStatus", -1); payStatus >= 0 {
		v := int8(payStatus)
		q.PayStatus = &v
	}

	list, total, err := c.service.List(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询订单详情，返回订单与每条明细（明细即打印任务）。
func (c *OrderController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	order, items, err := c.service.AdminDetail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"order": order, "items": items})
}

// Refund 对已支付订单发起全额退款。
// 入参：路径订单 ID + 退款原因；退款走 pay 模块，余额立即到账、微信等回调。
func (c *OrderController) Refund(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.OrderRefundRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Refund(ctx.Context(), id, req.Reason, types.ID(support.UserID(ctx.Context()))); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// PickupVerify 门店核销取件码。
// 入参：printShopId + code；核销人取自后台令牌。
func (c *OrderController) PickupVerify(ctx http.Context) http.Response {
	var req request.PickupVerifyRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	order, err := c.service.PickupVerify(ctx.Context(), types.ID(req.PrintShopID), req.Code, types.ID(support.UserID(ctx.Context())))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, order)
}
