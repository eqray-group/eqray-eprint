package gateway

import (
	payEnums "eqray-eprint/app/modules/pay/enums"
)

// registry 保存 channelCode → Gateway 的映射。
// 新增渠道只在这里注册，业务层按配置编码取实例，不写渠道分支。
var registry = map[string]Gateway{
	payEnums.CodeWechatMini: NewWechatMiniGateway(),
	payEnums.CodeWechatMp:   NewWechatMpGateway(),
	payEnums.CodeWechatH5:   NewWechatH5Gateway(),
	payEnums.CodeAlipayWap:  NewAlipayWapGateway(),
}

// Get 按渠道配置编码取网关实现，第二个返回值表示是否支持该编码。
// 余额支付不走网关（直接扣钱包），调用方需自行分支。
func Get(code string) (Gateway, bool) {
	gw, ok := registry[code]

	return gw, ok
}

// Codes 返回当前已注册的渠道编码，供后台下拉与日志排查使用。
func Codes() []string {
	codes := make([]string, 0, len(registry))
	for code := range registry {
		codes = append(codes, code)
	}

	return codes
}
