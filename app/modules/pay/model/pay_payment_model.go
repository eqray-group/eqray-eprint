package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PayPaymentTable 是支付单表名常量。
const PayPaymentTable = "pay_payment"

// PayPaymentModel 对应 pay_payment，打印订单支付与钱包充值共用。
// PaymentNo 即微信 out_trade_no；ChannelCode 决定用哪套商户配置，
// Channel 只用于统计口径（wechat / alipay / balance）。
type PayPaymentModel struct {
	baseModel.BaseModel
	PaymentNo      string          `db:"payment_no" json:"paymentNo"`
	BizType        string          `db:"biz_type" json:"bizType"`
	BizID          types.ID        `db:"biz_id" json:"bizId"`
	MemberUserID   types.ID        `db:"member_user_id" json:"memberUserId"`
	Channel        string          `db:"channel" json:"channel"`
	ChannelCode    string          `db:"channel_code" json:"channelCode"`
	Scene          string          `db:"scene" json:"scene"`
	ClientIP       string          `db:"client_ip" json:"clientIp"`
	Amount         int64           `db:"amount" json:"amount"`
	PaidAmount     int64           `db:"paid_amount" json:"paidAmount"`
	RefundAmount   int64           `db:"refund_amount" json:"refundAmount"`
	Status         int8            `db:"status" json:"status"`
	ChannelTradeNo string          `db:"channel_trade_no" json:"channelTradeNo"`
	OpenID         string          `db:"open_id" json:"openId"`
	ExpireTime     *types.DateTime `db:"expire_time" json:"expireTime"`
	SuccessTime    *types.DateTime `db:"success_time" json:"successTime"`
	CloseTime      *types.DateTime `db:"close_time" json:"closeTime"`
	Remark         string          `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PayPaymentModel) TableName() string {
	return PayPaymentTable
}

// IsPending 判断支付单是否仍可复用（待支付或支付中）。
// 无参数；返回 true 表示同一业务单再次下单时应复用本单而不是新建。
func (m *PayPaymentModel) IsPending() bool {
	if m == nil {
		return false
	}

	return m.Status == PayPaymentStatusPending || m.Status == PayPaymentStatusPaying
}

// 支付单状态，与 pay/enums 保持一致，放在 model 侧便于 IsPending 直接引用。
const (
	PayPaymentStatusPending  int8 = 0
	PayPaymentStatusPaying   int8 = 1
	PayPaymentStatusPaid     int8 = 2
	PayPaymentStatusFailed   int8 = 3
	PayPaymentStatusClosed   int8 = 4
	PayPaymentStatusRefunded int8 = 5
)
