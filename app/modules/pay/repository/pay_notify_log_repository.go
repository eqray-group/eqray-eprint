package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/pay/model"
)

const notifyLogTable = model.PayNotifyLogTable

// PayNotifyLogRepository 是 pay_notify_log 表的仓储。
// 日志只追加不修改，用于回调排查与对账取证。
type PayNotifyLogRepository struct{}

// NewPayNotifyLogRepository 构造回调日志仓储。
func NewPayNotifyLogRepository() *PayNotifyLogRepository {
	return &PayNotifyLogRepository{}
}

// Create 写入一条回调日志；日志失败不阻断主流程，调用方可忽略返回值。
func (r *PayNotifyLogRepository) Create(ctx context.Context, m *model.PayNotifyLogModel) error {
	now := types.Now()
	m.CreatedAt = &now
	if m.ID == 0 {
		m.ID = newID()
	}

	_, err := facades.DB().WithContext(ctx).Table(notifyLogTable).Insert(map[string]any{
		"id":            m.ID,
		"channel":       m.Channel,
		"notify_type":   m.NotifyType,
		"out_trade_no":  m.OutTradeNo,
		"channel_no":    m.ChannelNo,
		"body":          m.Body,
		"verified":      m.Verified,
		"handle_status": m.HandleStatus,
		"handle_msg":    m.HandleMsg,
		"created_at":    m.CreatedAt,
	})

	return err
}
