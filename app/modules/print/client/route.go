// Package client 注册 print 模块的会员端接口。
// 会员身份取自会员令牌，入参不出现 memberId。
package client

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/print/client/controller"
)

// Prefix 是会员端打印接口前缀。
const Prefix = "/api/v1/client/print"

// Register 注册会员端门店、下单与订单查询接口。
func Register() {
	shops := controller.NewShopController()
	orders := controller.NewOrderController()
	files := controller.NewFileController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		memberAuth := r.Middleware(&baseMiddleware.MemberAuth{})

		// ---------- 文件 ----------
		// 上传即同步转换；DOCX 依赖服务器安装 LibreOffice，失败原因写在文件记录上
		memberAuth.Post("/files", files.Upload)
		memberAuth.Get("/files/{id}", files.Show)
		memberAuth.Post("/files/{id}/convert", files.Convert)

		// ---------- 门店 ----------
		memberAuth.Get("/shops", shops.Index)
		// 可售组合：前端据此渲染纸张/色彩/单双面与单价
		memberAuth.Get("/shops/{id}/sellable", shops.Sellable)
		// 轮播图：门店自有 + 全平台，仅启用
		memberAuth.Get("/shops/{id}/banners", shops.Banners)

		// ---------- 订单 ----------
		memberAuth.Post("/orders", orders.Store)
		memberAuth.Get("/orders", orders.Index)
		memberAuth.Get("/orders/{orderNo}", orders.Show)
		// 未支付订单重新发起支付
		memberAuth.Post("/orders/{orderNo}/pay", orders.Pay)
	})
}
