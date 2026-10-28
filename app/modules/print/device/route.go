// Package device 注册打印客户端（Windows 端）接口。
// 客户端使用后台账号登录，因此走后台令牌鉴权，不走会员令牌。
package device

import (
	"github.com/goravel/framework/contracts/route"

	"eqray-eprint/app/facades"
	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/print/device/controller"
)

// Prefix 是打印客户端接口前缀。
const Prefix = "/api/v1/print/device"

// Register 注册客户端登记、心跳、拉任务、回执与核销接口。
func Register() {
	deviceController := controller.NewDeviceController()

	facades.Route().Prefix(Prefix).Group(func(r route.Router) {
		auth := r.Middleware(&baseMiddleware.TokenAuth{})

		// 登记设备：首次登录创建，之后更新归属并刷新心跳
		auth.Post("/devices/register", deviceController.Register)
		// 登录后选门店：只返回当前账号绑定的门店
		auth.Get("/shops", deviceController.MyShops)
		// 打印机上报：客户端勾选本机队列并把纸张 / 彩色 / 双面能力写到服务端
		auth.Get("/devices/{id}/printers", deviceController.Printers)
		auth.Post("/devices/{id}/printers/sync", deviceController.SyncPrinters)
		auth.Post("/devices/{id}/heartbeat", deviceController.Heartbeat)
		// 拉任务会顺带刷新心跳，客户端不必额外发心跳
		auth.Get("/devices/{id}/tasks", deviceController.Tasks)
		auth.Post("/devices/{id}/tasks/report", deviceController.Report)
		auth.Post("/shops/{id}/pickup", deviceController.Pickup)
	})
}
