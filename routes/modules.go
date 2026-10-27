package routes

import (
	"eqray-eprint/app/modules/infra"
	"eqray-eprint/app/modules/member"
	"eqray-eprint/app/modules/system"
)

// Modules 聚合各业务模块的路由，在 bootstrap/app.go 的 WithRouting 中调用。
func Modules() {
	system.Register()
	infra.Register()
	member.Register()
}
