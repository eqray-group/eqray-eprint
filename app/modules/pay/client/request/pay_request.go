package request

import (
	"github.com/goravel/framework/contracts/http"
)

// CreatePaymentRequest 是会员发起支付的入参。
// 金额由服务端按订单重新校验，前端传入只用于快速失败提示。
type CreatePaymentRequest struct {
	BizType     string `form:"bizType" json:"bizType"`
	BizID       string `form:"bizId" json:"bizId"`
	Amount      int64  `form:"amount" json:"amount"`
	ChannelCode string `form:"channelCode" json:"channelCode"`
	OpenID      string `form:"openId" json:"openId"`
	ClientIP    string `form:"clientIp" json:"clientIp"`
	Description string `form:"description" json:"description"`
}

// Rules 返回校验规则。
func (r *CreatePaymentRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"bizType":     "required|in:print_order",
		"bizId":       "required",
		"amount":      "required|integer|min:1",
		"channelCode": "required|max:32",
		"openId":      "max:64",
		"clientIp":    "max:50",
		"description": "required|max:127",
	}
}

// Messages 覆盖默认提示。
func (r *CreatePaymentRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"bizType.required":     "业务类型不能为空",
		"bizType.in":           "不支持的业务类型",
		"bizId.required":       "业务单号不能为空",
		"amount.required":      "支付金额不能为空",
		"amount.min":           "支付金额必须大于 0",
		"amount.integer":       "支付金额必须是整数分",
		"channelCode.required": "请选择支付方式",
		"description.required": "商品描述不能为空",
		"description.max":      "商品描述不能超过 127 个字符",
	}
}

// RechargeRequest 是钱包充值入参，只允许微信渠道。
type RechargeRequest struct {
	Amount      int64  `form:"amount" json:"amount"`
	ChannelCode string `form:"channelCode" json:"channelCode"`
	OpenID      string `form:"openId" json:"openId"`
	ClientIP    string `form:"clientIp" json:"clientIp"`
}

// Rules 返回校验规则。
func (r *RechargeRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"amount":      "required|integer|min:1",
		"channelCode": "required|in:wechat_mini,wechat_mp,wechat_h5",
		"openId":      "max:64",
		"clientIp":    "max:50",
	}
}

// Messages 覆盖默认提示。
func (r *RechargeRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"amount.required":      "充值金额不能为空",
		"amount.min":           "充值金额必须大于 0",
		"amount.integer":       "充值金额必须是整数分",
		"channelCode.required": "请选择支付方式",
		"channelCode.in":       "充值只支持微信支付",
	}
}
