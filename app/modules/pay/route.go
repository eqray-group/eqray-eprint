// Package pay 聚合支付与钱包模块的路由。
// 后台维护渠道配置与退款，会员端下单/充值/查单，回调独立且免鉴权。
package pay

import (
	"eqray-eprint/app/modules/pay/admin"
	"eqray-eprint/app/modules/pay/client"
	"eqray-eprint/app/modules/pay/notify"
)

// Register 注册支付模块的全部路由，由 routes.Modules() 统一调用。
func Register() {
	admin.Register()
	client.Register()
	notify.Register()
}
