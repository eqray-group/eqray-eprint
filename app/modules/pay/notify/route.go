// Package notify 注册支付渠道回调路由。
// 回调由渠道服务器发起，没有会员令牌，因此只做验签与幂等，不挂鉴权中间件。
package notify

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/pay/notify/controller"
)

// Prefix 是支付回调前缀，需与 pay_channel_config.notify_url 保持一致。
const Prefix = "/api/v1/pay/notify"

// Register 注册微信支付回调；退款回调待退款服务落地后补上。
func Register() {
	wechatNotify := controller.NewWechatNotifyController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		r.Post("/wechat/payment", wechatNotify.Payment)
		r.Post("/wechat/refund", wechatNotify.Refund)
		// 支付宝异步通知为 form 表单，验签后回纯文本 success
		r.Post("/alipay/payment", wechatNotify.AlipayPayment)
	})
}
