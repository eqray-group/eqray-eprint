package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PayRefundTable 是退款单表名常量。
const PayRefundTable = "pay_refund"

// PayRefundModel 对应 pay_refund，一次支付可以有多笔退款。
// 累计退款金额不得超过支付单的实付金额；微信渠道异步回调置成功，
// 余额渠道直接入账并立即成功。
type PayRefundModel struct {
	baseModel.BaseModel
	RefundNo        string          `db:"refund_no" json:"refundNo"`
	PaymentID       types.ID        `db:"payment_id" json:"paymentId"`
	PaymentNo       string          `db:"payment_no" json:"paymentNo"`
	BizType         string          `db:"biz_type" json:"bizType"`
	BizID           types.ID        `db:"biz_id" json:"bizId"`
	MemberUserID    types.ID        `db:"member_user_id" json:"memberUserId"`
	Channel         string          `db:"channel" json:"channel"`
	ChannelCode     string          `db:"channel_code" json:"channelCode"`
	Amount          int64           `db:"amount" json:"amount"`
	Reason          string          `db:"reason" json:"reason"`
	Status          int8            `db:"status" json:"status"`
	ChannelRefundNo string          `db:"channel_refund_no" json:"channelRefundNo"`
	SuccessTime     *types.DateTime `db:"success_time" json:"successTime"`
	FailReason      string          `db:"fail_reason" json:"failReason"`
	OperatorID      types.ID        `db:"operator_id" json:"operatorId"`
}

// TableName 返回表名。
func (m *PayRefundModel) TableName() string {
	return PayRefundTable
}
