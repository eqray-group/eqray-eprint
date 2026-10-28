// Package admin 注册 pay 模块的后台管理接口。
// 会员支付本身在 client 端，这里只做渠道配置、支付单与退款查询、人工退款。
package admin

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/pay/admin/controller"
)

// Prefix 是后台支付管理接口前缀。
const Prefix = "/api/v1/pay"

// Register 注册后台支付接口，由 pay.Register() 统一调用。
func Register() {
	channelController := controller.NewChannelConfigController()
	refundController := controller.NewRefundController()
	paymentController := controller.NewPaymentController()
	walletController := controller.NewWalletController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		auth := r.Middleware(&baseMiddleware.TokenAuth{})

		perm := func(code string) route.Router {
			return auth.Middleware(&baseMiddleware.Permission{Code: code})
		}
		operate := func(typ, sub string) *baseMiddleware.OperateLog {
			return &baseMiddleware.OperateLog{Type: typ, SubType: sub}
		}

		// ---------- 支付渠道配置 ----------
		// 列表不回显密钥，详情接口才回显，供编辑页使用
		perm("pay:channel:list").Get("/channel-configs", channelController.Index)
		perm("pay:channel:list").Get("/channel-configs/{id}", channelController.Show)
		perm("pay:channel:create").Middleware(operate("支付管理", "新增支付渠道")).
			Post("/channel-configs", channelController.Store)
		perm("pay:channel:update").Middleware(operate("支付管理", "编辑支付渠道")).
			Put("/channel-configs/{id}", channelController.Update)
		perm("pay:channel:delete").Middleware(operate("支付管理", "删除支付渠道")).
			Delete("/channel-configs/{id}", channelController.Destroy)

		// ---------- 支付流水 ----------
		// 只读查询，支付单由会员端创建，后台不代下单
		perm("pay:payment:list").Get("/payments", paymentController.Index)

		// ---------- 会员钱包 ----------
		// 余额与流水只读查询；调账单独授权
		perm("pay:wallet:list").Get("/wallets/{memberId}", walletController.Show)
		perm("pay:wallet:list").Get("/wallet-transactions", walletController.Index)
		// 人工调账直接改余额，单独授权并留痕
		perm("pay:wallet:adjust").Middleware(operate("支付管理", "钱包人工调账")).
			Post("/wallets/{memberId}/adjust", walletController.Adjust)

		// ---------- 退款 ----------
		// 退款直接动钱，必须单独授权；发起人取自后台令牌，不接受前端传入
		perm("pay:refund:list").Get("/refunds", refundController.Index)
		perm("pay:refund:create").Middleware(operate("支付管理", "发起退款")).
			Post("/refunds", refundController.Store)
	})
}
