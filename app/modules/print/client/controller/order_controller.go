package controller

import (
	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/client/request"
	"eqray-eprint/app/modules/print/service"
)

// OrderController 提供会员端下单、支付与订单查询。
type OrderController struct {
	service *service.OrderService
}

// NewOrderController 构造会员端订单控制器。
func NewOrderController() *OrderController {
	return &OrderController{service: service.NewOrderService()}
}

// Store 会员下单：服务端计价并返回支付调起参数。
func (c *OrderController) Store(ctx http.Context) http.Response {
	var req request.CreateOrderRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}
	if len(req.Items) == 0 {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "请至少选择一个文件"))
	}

	items := make([]service.CreateItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		if item.FileID <= 0 {
			return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "文件 ID 不合法"))
		}

		items = append(items, service.CreateItemInput{
			FileID:      types.ID(item.FileID),
			PaperType:   item.PaperType,
			ColorType:   item.ColorType,
			PrintMode:   item.PrintMode,
			Orientation: item.Orientation,
			PageRange:   item.PageRange,
			Copies:      item.Copies,
		})
	}

	out, err := c.service.Create(ctx.Context(), service.CreateOrderInput{
		MemberID:    support.MemberID(ctx.Context()),
		ShopID:      types.ID(req.PrintShopID),
		ChannelCode: req.ChannelCode,
		OpenID:      req.OpenID,
		ClientIP:    ctx.Request().Origin().RemoteAddr,
		Remark:      req.Remark,
		Items:       items,
	})
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, out)
}

// Pay 为未支付订单重新发起支付。
func (c *OrderController) Pay(ctx http.Context) http.Response {
	orderNo := ctx.Request().Route("orderNo")
	if orderNo == "" {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "订单号不能为空"))
	}

	var req request.PayOrderRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	out, err := c.service.PayOrder(
		ctx.Context(),
		support.MemberID(ctx.Context()),
		orderNo,
		req.ChannelCode,
		req.OpenID,
		ctx.Request().Origin().RemoteAddr,
	)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, out)
}

// Index 分页查询本人订单。
func (c *OrderController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.OrderListQuery{
		MemberID: support.MemberID(ctx.Context()),
		OrderNo:  ctx.Request().Query("orderNo"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if status := ctx.Request().QueryInt("orderStatus", -1); status >= 0 {
		v := int8(status)
		q.OrderStatus = &v
	}

	list, total, err := c.service.List(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询本人订单详情（含明细）。
func (c *OrderController) Show(ctx http.Context) http.Response {
	orderNo := ctx.Request().Route("orderNo")
	if orderNo == "" {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "订单号不能为空"))
	}

	order, items, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()), orderNo)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"order": order, "items": items})
}
