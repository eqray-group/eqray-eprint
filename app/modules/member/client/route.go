// Package client 注册会员模块的 C 端接口（小程序 / H5 / App）。
// C 端访问的是「会员本人」的数据，只做登录态校验，不进入后台 RBAC 权限体系。
package client

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/member/client/controller"
)

// Prefix 是 C 端接口前缀。
const Prefix = "/api/v1/client/member"

// Register 注册 C 端接口，由 member.Register() 统一调用。
// 无参数无返回；签到、积分、等级相关接口已随体系下线移除。
func Register() {
	authController := controller.NewMemberAuthController()
	profileController := controller.NewMemberProfileController()
	addressController := controller.NewMemberAddressController()
	cardController := controller.NewMemberCardController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		// ---------- 免登录 ----------
		// 登录、续期与授权地址不需要令牌
		r.Get("/auth/social/{type}/authorize-url", authController.AuthorizeURL)
		r.Post("/auth/social/login", authController.SocialLogin)
		r.Post("/auth/sms-code", authController.SmsCode)
		r.Post("/auth/sms-login", authController.SmsLogin)
		r.Post("/auth/social-bind-login", authController.SocialBindLogin)
		r.Post("/auth/refresh", authController.Refresh)

		// ---------- 需要会员令牌 ----------
		// 与后台令牌作用域隔离，互不通用
		memberAuth := r.Middleware(&baseMiddleware.MemberAuth{})
		memberAuth.Get("/auth/me", authController.Me)
		memberAuth.Post("/auth/logout", authController.Logout)
		memberAuth.Post("/auth/social/bind", authController.Bind)
		memberAuth.Get("/auth/social/bindings", authController.Bindings)
		memberAuth.Delete("/auth/social/bindings/{id}", authController.Unbind)

		// ---------- 我的资料 ----------
		memberAuth.Get("/profile", profileController.Show)
		memberAuth.Put("/profile", profileController.Update)

		// ---------- 我的会员卡 ----------
		// 折扣随卡走，与会员等级无关；未发卡返回业务错误码
		memberAuth.Get("/card", cardController.Show)

		// ---------- 我的收货地址 ----------
		memberAuth.Get("/addresses", addressController.Index)
		memberAuth.Post("/addresses", addressController.Store)
		memberAuth.Put("/addresses/{id}", addressController.Update)
		memberAuth.Delete("/addresses/{id}", addressController.Destroy)
		memberAuth.Put("/addresses/{id}/default", addressController.SetDefault)

		// 后续 C 端业务接口统一下沉到 client/controller，会员身份一律用
		// support.MemberID(ctx.Context()) 取，禁止信任前端传入的会员 ID。
	})
}
