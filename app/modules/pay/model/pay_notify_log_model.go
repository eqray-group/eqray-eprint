package model

import (
	"eqray-eprint/app/modules/base/types"
)

// PayNotifyLogTable 是支付回调日志表名常量。
const PayNotifyLogTable = "pay_notify_log"

// PayNotifyLogModel 对应 pay_notify_log，记录渠道回调的原始报文与处理结果。
// 日志表只有 created_at，不参与软删除与审计字段维护。
type PayNotifyLogModel struct {
	ID           types.ID        `db:"id" json:"id"`
	Channel      string          `db:"channel" json:"channel"`
	NotifyType   string          `db:"notify_type" json:"notifyType"`
	OutTradeNo   string          `db:"out_trade_no" json:"outTradeNo"`
	ChannelNo    string          `db:"channel_no" json:"channelNo"`
	Body         string          `db:"body" json:"-"`
	Verified     int8            `db:"verified" json:"verified"`
	HandleStatus int8            `db:"handle_status" json:"handleStatus"`
	HandleMsg    string          `db:"handle_msg" json:"handleMsg"`
	CreatedAt    *types.DateTime `db:"created_at" json:"createdAt"`
}

// TableName 返回表名。
func (m *PayNotifyLogModel) TableName() string {
	return PayNotifyLogTable
}
