// Package controller 处理支付渠道的异步回调。
// 回调不走会员令牌，安全完全依赖验签：验签失败一律拒绝，不得当作有效支付。
package controller

import (
	"io"

	"github.com/goravel/framework/contracts/http"

	payEnums "eqray-eprint/app/modules/pay/enums"
	"eqray-eprint/app/modules/pay/service"
)

// wechatNotifyMaxBody 限制回调报文大小，防止超大 body 打满内存。
const wechatNotifyMaxBody = 1 << 20

// WechatNotifyController 处理微信支付与退款回调。
type WechatNotifyController struct {
	service *service.PaymentService
	refunds *service.RefundService
}

// NewWechatNotifyController 构造微信回调控制器。
func NewWechatNotifyController() *WechatNotifyController {
	return &WechatNotifyController{
		service: service.NewPaymentService(),
		refunds: service.NewRefundService(),
	}
}

// Payment 接收微信支付结果通知。
// 返回 200 + code=SUCCESS 表示处理完成，微信不再重试；
// 返回 500 表示处理失败，微信会按策略重试，因此业务必须幂等。
func (c *WechatNotifyController) Payment(ctx http.Context) http.Response {
	req := ctx.Request().Origin()

	body, err := io.ReadAll(io.LimitReader(req.Body, wechatNotifyMaxBody))
	if err != nil {
		return ctx.Response().Json(400, map[string]string{
			"code":    "FAIL",
			"message": "读取报文失败",
		})
	}
	defer req.Body.Close()

	if err := c.service.HandleWechatNotify(ctx.Context(), req.Header, body); err != nil {
		return ctx.Response().Json(500, map[string]string{
			"code":    "FAIL",
			"message": "处理失败",
		})
	}

	return ctx.Response().Json(200, map[string]string{
		"code":    "SUCCESS",
		"message": "成功",
	})
}

// AlipayPayment 接收支付宝异步通知（form 表单）。
// 支付宝只认纯文本 success/failure，返回 JSON 会被当成失败并持续重试。
func (c *WechatNotifyController) AlipayPayment(ctx http.Context) http.Response {
	req := ctx.Request().Origin()

	body, err := io.ReadAll(io.LimitReader(req.Body, wechatNotifyMaxBody))
	if err != nil {
		return ctx.Response().String(200, "failure")
	}
	defer req.Body.Close()

	if err := c.service.HandleNotify(ctx.Context(), payEnums.ChannelAlipay, req.Header, body); err != nil {
		return ctx.Response().String(200, "failure")
	}

	return ctx.Response().String(200, "success")
}

// Refund 接收微信退款结果通知。
// 处理约定与支付回调一致：200 表示已处理，500 让微信重试；退款单已成功时直接回 200。
func (c *WechatNotifyController) Refund(ctx http.Context) http.Response {
	req := ctx.Request().Origin()

	body, err := io.ReadAll(io.LimitReader(req.Body, wechatNotifyMaxBody))
	if err != nil {
		return ctx.Response().Json(400, map[string]string{
			"code":    "FAIL",
			"message": "读取报文失败",
		})
	}
	defer req.Body.Close()

	if err := c.refunds.HandleWechatNotify(ctx.Context(), req.Header, body); err != nil {
		return ctx.Response().Json(500, map[string]string{
			"code":    "FAIL",
			"message": "处理失败",
		})
	}

	return ctx.Response().Json(200, map[string]string{
		"code":    "SUCCESS",
		"message": "成功",
	})
}
