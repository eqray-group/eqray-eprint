package repository

import (
	"context"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/pay/model"
)

const refundTable = model.PayRefundTable

// RefundFilter 是退款单列表的筛选条件，零值表示不限制。
type RefundFilter struct {
	PaymentNo string
	MemberID  types.ID
	Status    *int8
}

// PayRefundRepository 是 pay_refund 表的仓储。
type PayRefundRepository struct{}

// NewPayRefundRepository 构造退款单仓储。
func NewPayRefundRepository() *PayRefundRepository {
	return &PayRefundRepository{}
}

// Create 插入退款单。
func (r *PayRefundRepository) Create(ctx context.Context, m *model.PayRefundModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(refundTable).Insert(map[string]any{
		"id":                m.ID,
		"refund_no":         m.RefundNo,
		"payment_id":        m.PaymentID,
		"payment_no":        m.PaymentNo,
		"biz_type":          m.BizType,
		"biz_id":            m.BizID,
		"member_user_id":    m.MemberUserID,
		"channel":           m.Channel,
		"channel_code":      m.ChannelCode,
		"amount":            m.Amount,
		"reason":            m.Reason,
		"status":            m.Status,
		"channel_refund_no": m.ChannelRefundNo,
		"success_time":      m.SuccessTime,
		"fail_reason":       m.FailReason,
		"operator_id":       m.OperatorID,
		"creator":           m.Creator,
		"created_at":        m.CreatedAt,
		"updater":           m.Updater,
		"updated_at":        m.UpdatedAt,
	})

	return err
}

// FindByNo 按退款单号查询；未命中返回 (nil, nil)。
func (r *PayRefundRepository) FindByNo(ctx context.Context, refundNo string) (*model.PayRefundModel, error) {
	var m model.PayRefundModel

	if err := facades.DB().WithContext(ctx).Table(refundTable).
		Where("refund_no", refundNo).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// UpdateStatus 在事务内更新退款单状态与附加信息。
func (r *PayRefundRepository) UpdateStatus(ctx context.Context, tx db.Tx, id types.ID, status int8, values map[string]any) error {
	update := map[string]any{
		"status":     status,
		"updated_at": types.Now(),
	}
	for k, v := range values {
		update[k] = v
	}

	_, err := tx.Table(refundTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(update)

	return err
}

// Paginate 按条件分页查询退款单。
func (r *PayRefundRepository) Paginate(ctx context.Context, filter RefundFilter, pageNo, pageSize int) ([]model.PayRefundModel, int64, error) {
	list := make([]model.PayRefundModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(refundTable).WhereNull("deleted_at")
	if filter.PaymentNo != "" {
		query = query.Where("payment_no", filter.PaymentNo)
	}
	if filter.MemberID > 0 {
		query = query.Where("member_user_id", filter.MemberID)
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
