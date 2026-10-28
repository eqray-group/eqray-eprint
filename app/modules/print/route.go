// Package print 注册打印模块的全部路由：后台管理、会员端、打印客户端。
package print

import (
	"eqray-eprint/app/modules/print/admin"
	"eqray-eprint/app/modules/print/client"
	"eqray-eprint/app/modules/print/device"
)

// Register 注册三端路由，由 routes/modules.go 调用。
func Register() {
	admin.Register()
	client.Register()
	device.Register()
}
