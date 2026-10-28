// Package admin 注册 print 模块的后台管理接口。
package admin

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/print/admin/controller"
)

// Prefix 是后台打印管理接口前缀。
const Prefix = "/api/v1/print"

// Register 注册后台打印接口，由 print.Register() 统一调用。
func Register() {
	shops := controller.NewShopController()
	devices := controller.NewDeviceController()
	printers := controller.NewPrinterController()
	prices := controller.NewPriceController()
	orders := controller.NewOrderController()
	banners := controller.NewBannerController()
	files := controller.NewFileController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		auth := r.Middleware(&baseMiddleware.TokenAuth{})

		perm := func(code string) route.Router {
			return auth.Middleware(&baseMiddleware.Permission{Code: code})
		}
		operate := func(typ, sub string) *baseMiddleware.OperateLog {
			return &baseMiddleware.OperateLog{Type: "打印管理", SubType: sub}
		}

		// ---------- 门店 ----------
		perm("print:shop:list").Get("/shops", shops.Index)
		perm("print:shop:list").Get("/shops/{id}", shops.Show)
		perm("print:shop:create").Middleware(operate("打印管理", "新增打印门店")).
			Post("/shops", shops.Store)
		perm("print:shop:update").Middleware(operate("打印管理", "编辑打印门店")).
			Put("/shops/{id}", shops.Update)
		perm("print:shop:delete").Middleware(operate("打印管理", "删除打印门店")).
			Delete("/shops/{id}", shops.Destroy)
		// 营业状态开关频次高，单列权限点，避免给了编辑权就能直接切换营业
		perm("print:shop:update").Middleware(operate("打印管理", "切换门店营业状态")).
			Put("/shops/{id}/status", shops.ChangeStatus)
		perm("print:shop:audit").Middleware(operate("打印管理", "审核打印门店")).
			Put("/shops/{id}/audit", shops.Audit)

		// ---------- 门店成员 ----------
		perm("print:staff:list").Get("/shops/{id}/staffs", shops.StaffIndex)
		perm("print:staff:create").Middleware(operate("打印管理", "绑定门店成员")).
			Post("/shops/{id}/staffs", shops.StaffStore)
		perm("print:staff:delete").Middleware(operate("打印管理", "解绑门店成员")).
			Delete("/shops/{id}/staffs/{staffId}", shops.StaffDestroy)

		// ---------- 设备 ----------
		perm("print:device:list").Get("/devices", devices.Index)
		perm("print:device:list").Get("/devices/{id}", devices.Show)
		perm("print:device:create").Middleware(operate("打印管理", "新增打印设备")).
			Post("/devices", devices.Store)
		perm("print:device:update").Middleware(operate("打印管理", "编辑打印设备")).
			Put("/devices/{id}", devices.Update)
		perm("print:device:delete").Middleware(operate("打印管理", "删除打印设备")).
			Delete("/devices/{id}", devices.Destroy)

		// ---------- 打印机 ----------
		perm("print:printer:list").Get("/printers", printers.Index)
		perm("print:printer:list").Get("/printers/{id}", printers.Show)
		perm("print:printer:create").Middleware(operate("打印管理", "新增打印机")).
			Post("/printers", printers.Store)
		perm("print:printer:update").Middleware(operate("打印管理", "编辑打印机")).
			Put("/printers/{id}", printers.Update)
		perm("print:printer:delete").Middleware(operate("打印管理", "停用打印机")).
			Put("/printers/{id}/disable", printers.Disable)

		// ---------- 价格 ----------
		perm("print:price:list").Get("/price-configs", prices.Index)
		perm("print:price:list").Get("/price-configs/warnings", prices.Warnings)
		perm("print:price:create").Middleware(operate("打印管理", "新增价格规则")).
			Post("/price-configs", prices.Store)
		perm("print:price:update").Middleware(operate("打印管理", "修改价格规则")).
			Put("/price-configs/{id}", prices.Update)
		perm("print:price:delete").Middleware(operate("打印管理", "删除价格规则")).
			Delete("/price-configs/{id}", prices.Destroy)

		// ---------- 打印文件 ----------
		// 只读：文件内容属于会员，后台不提供下载与删除
		perm("print:file:list").Get("/files", files.Index)
		// 控制台概览与订单同权限：能看到订单就能看到汇总
		perm("print:order:list").Get("/statistics/overview", files.Statistics)

		// ---------- 轮播图 ----------
		perm("print:banner:list").Get("/banners", banners.Index)
		perm("print:banner:list").Get("/banners/{id}", banners.Show)
		perm("print:banner:create").Middleware(operate("打印管理", "新增轮播图")).
			Post("/banners", banners.Store)
		perm("print:banner:update").Middleware(operate("打印管理", "编辑轮播图")).
			Put("/banners/{id}", banners.Update)
		perm("print:banner:delete").Middleware(operate("打印管理", "删除轮播图")).
			Delete("/banners/{id}", banners.Destroy)

		// ---------- 订单 ----------
		perm("print:order:list").Get("/orders", orders.Index)
		// 详情含明细与任务状态，用于排查打印失败
		perm("print:order:list").Get("/orders/{id}", orders.Show)
		// 核销会改订单状态并放行取件，单独授权并留痕
		perm("print:order:pickup").Middleware(operate("打印管理", "取件核销")).
			Post("/orders/pickup", orders.PickupVerify)
		// 退款直接动钱，单独授权；金额取订单实付，不接受前端传值
		perm("print:order:refund").Middleware(operate("打印管理", "订单退款")).
			Post("/orders/{id}/refund", orders.Refund)
	})
}
