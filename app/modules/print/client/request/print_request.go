// Package request 定义 print 模块会员端接口入参。
package request

import (
	"github.com/goravel/framework/contracts/http"
)

// OrderItemRequest 是下单明细：一个文件 + 一套打印参数。
type OrderItemRequest struct {
	FileID      int64  `form:"fileId" json:"fileId"`
	PaperType   string `form:"paperType" json:"paperType"`
	ColorType   string `form:"colorType" json:"colorType"`
	PrintMode   string `form:"printMode" json:"printMode"`
	Orientation string `form:"orientation" json:"orientation"`
	PageRange   string `form:"pageRange" json:"pageRange"`
	Copies      int    `form:"copies" json:"copies"`
}

// CreateOrderRequest 是会员下单入参；金额由服务端计算，不接受前端传值。
type CreateOrderRequest struct {
	PrintShopID int64              `form:"printShopId" json:"printShopId"`
	ChannelCode string             `form:"channelCode" json:"channelCode"`
	OpenID      string             `form:"openid" json:"openid"`
	Remark      string             `form:"remark" json:"remark"`
	Items       []OrderItemRequest `form:"items" json:"items"`
}

// Rules 返回校验规则。
func (r *CreateOrderRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"printShopId": "required|integer|min:1",
		"channelCode": "required|in:balance,wechat_mini,wechat_mp,wechat_h5,alipay_wap",
		"remark":      "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *CreateOrderRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"printShopId.required": "请选择打印门店",
		"channelCode.required": "请选择支付方式",
		"channelCode.in":       "支付方式不支持",
	}
}

// PayOrderRequest 是订单重新发起支付的入参。
type PayOrderRequest struct {
	ChannelCode string `form:"channelCode" json:"channelCode"`
	OpenID      string `form:"openid" json:"openid"`
}

// Rules 返回校验规则。
func (r *PayOrderRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"channelCode": "required|in:balance,wechat_mini,wechat_mp,wechat_h5,alipay_wap",
	}
}

// Messages 覆盖默认提示。
func (r *PayOrderRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{"channelCode.in": "支付方式不支持"}
}
