// Package client 注册 pay 模块的会员端接口。
// 会员端只操作自己的钱包与支付单，身份取自会员令牌，入参不出现 memberId。
package client

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/pay/client/controller"
)

// Prefix 是会员端支付接口前缀。
const Prefix = "/api/v1/client/pay"

// Register 注册会员端钱包与支付接口。
func Register() {
	walletController := controller.NewWalletController()
	paymentController := controller.NewPaymentController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		memberAuth := r.Middleware(&baseMiddleware.MemberAuth{})

		// ---------- 钱包 ----------
		memberAuth.Get("/wallet", walletController.Show)
		memberAuth.Get("/wallet/transactions", walletController.Transactions)

		// ---------- 支付单 ----------
		// 创建支付单：渠道支付返回调起参数，余额支付直接完成扣款
		memberAuth.Post("/payments", paymentController.Store)
		memberAuth.Get("/payments", paymentController.List)
		memberAuth.Get("/payments/{paymentNo}", paymentController.Show)
		// 主动查单，供前端支付后轮询兜底
		memberAuth.Post("/payments/{paymentNo}/sync", paymentController.Sync)

		// ---------- 充值 ----------
		memberAuth.Post("/recharge", paymentController.Recharge)
	})
}
