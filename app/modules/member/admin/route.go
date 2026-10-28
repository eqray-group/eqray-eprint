// Package admin 注册会员模块的后台管理接口。
// 管理端访问的是「任意会员」的数据，必须挂后台令牌与权限点。
package admin

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/member/admin/controller"
)

// Prefix 是后台管理端接口前缀。
const Prefix = "/api/v1/member"

// Register 注册后台管理接口，由 member.Register() 统一调用。
// 无参数无返回，直接向路由注册会员、会员卡与第三方登录配置三组接口。
func Register() {
	memberController := controller.NewMemberController()
	cardController := controller.NewMemberCardController()
	socialClientController := controller.NewMemberSocialClientController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		auth := r.Middleware(&baseMiddleware.TokenAuth{})

		// 权限中间件在 base，member 只声明自己的权限点
		perm := func(code string) route.Router {
			return auth.Middleware(&baseMiddleware.Permission{Code: code})
		}

		// 写操作留痕，中间件注册在 Permission 之内
		operate := func(typ, sub string) *baseMiddleware.OperateLog {
			return &baseMiddleware.OperateLog{Type: typ, SubType: sub}
		}

		// ---------- 会员 ----------
		perm("member:user:list").Get("/users", memberController.Index)
		perm("member:user:list").Get("/users/{id}", memberController.Show)
		perm("member:user:create").Middleware(operate("会员管理", "新增会员")).
			Post("/users", memberController.Store)
		perm("member:user:update").Middleware(operate("会员管理", "编辑会员")).
			Put("/users/{id}", memberController.Update)
		perm("member:user:update").Middleware(operate("会员管理", "变更会员状态")).
			Put("/users/{id}/status", memberController.UpdateStatus)
		perm("member:user:delete").Middleware(operate("会员管理", "删除会员")).
			Delete("/users/{id}", memberController.Destroy)

		// 收货地址与第三方绑定：后台只读
		perm("member:user:list").Get("/users/{id}/addresses", memberController.Addresses)
		perm("member:user:list").Get("/users/{id}/social-bindings", memberController.SocialBindings)

		// ---------- 会员卡 ----------
		// 一个会员一张卡，按会员 ID 定位；卡号由服务端生成
		perm("member:card:list").Get("/member-cards", cardController.Index)
		perm("member:card:list").Get("/users/{id}/card", cardController.Show)
		perm("member:card:update").Middleware(operate("会员管理", "发放或修改会员卡")).
			Put("/users/{id}/card", cardController.Save)
		perm("member:card:delete").Middleware(operate("会员管理", "停用会员卡")).
			Put("/users/{id}/card/disable", cardController.Disable)

		// ---------- 第三方登录渠道配置 ----------
		perm("member:social-client:list").Get("/social-clients", socialClientController.Index)
		perm("member:social-client:list").Get("/social-clients/{id}", socialClientController.Show)
		perm("member:social-client:create").Middleware(operate("会员管理", "新增社交客户端")).
			Post("/social-clients", socialClientController.Store)
		perm("member:social-client:update").Middleware(operate("会员管理", "编辑社交客户端")).
			Put("/social-clients/{id}", socialClientController.Update)
		perm("member:social-client:delete").Middleware(operate("会员管理", "删除社交客户端")).
			Delete("/social-clients/{id}", socialClientController.Destroy)
	})
}
