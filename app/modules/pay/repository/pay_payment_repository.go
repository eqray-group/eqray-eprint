package repository

import (
	"context"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	payEnums "eqray-eprint/app/modules/pay/enums"
	"eqray-eprint/app/modules/pay/model"
)

const paymentTable = model.PayPaymentTable

// PayPaymentRepository 是 pay_payment 表的仓储。
type PayPaymentRepository struct{}

// NewPayPaymentRepository 构造支付单仓储。
func NewPayPaymentRepository() *PayPaymentRepository {
	return &PayPaymentRepository{}
}

// Create 插入支付单。
func (r *PayPaymentRepository) Create(ctx context.Context, m *model.PayPaymentModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(paymentTable).Insert(map[string]any{
		"id":               m.ID,
		"payment_no":       m.PaymentNo,
		"biz_type":         m.BizType,
		"biz_id":           m.BizID,
		"member_user_id":   m.MemberUserID,
		"channel":          m.Channel,
		"channel_code":     m.ChannelCode,
		"scene":            m.Scene,
		"client_ip":        m.ClientIP,
		"amount":           m.Amount,
		"paid_amount":      m.PaidAmount,
		"refund_amount":    m.RefundAmount,
		"status":           m.Status,
		"channel_trade_no": m.ChannelTradeNo,
		"open_id":          m.OpenID,
		"expire_time":      m.ExpireTime,
		"success_time":     m.SuccessTime,
		"close_time":       m.CloseTime,
		"remark":           m.Remark,
		"creator":          m.Creator,
		"created_at":       m.CreatedAt,
		"updater":          m.Updater,
		"updated_at":       m.UpdatedAt,
	})

	return err
}

// FindByNo 按商户支付单号查询；未命中返回 (nil, nil)。
func (r *PayPaymentRepository) FindByNo(ctx context.Context, paymentNo string) (*model.PayPaymentModel, error) {
	var m model.PayPaymentModel

	if err := facades.DB().WithContext(ctx).Table(paymentTable).
		Where("payment_no", paymentNo).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindPending 查询同一业务单未完成的支付单，用于重复下单时复用而不是新建。
// 返回最近一条；过期判断由服务层负责，避免把「已超时但未关单」当成可复用。
func (r *PayPaymentRepository) FindPending(ctx context.Context, bizType string, bizID types.ID, channelCode string) (*model.PayPaymentModel, error) {
	var m model.PayPaymentModel

	if err := facades.DB().WithContext(ctx).Table(paymentTable).
		Where("biz_type", bizType).
		Where("biz_id", bizID).
		Where("channel_code", channelCode).
		WhereIn("status", []any{payEnums.PaymentStatusPending, payEnums.PaymentStatusPaying}).
		WhereNull("deleted_at").
		OrderByDesc("id").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// UpdateStatus 在事务内更新支付单状态，供回调与关单使用。
// 参数 values 为需要一并写入的字段（如 channel_trade_no、success_time）。
func (r *PayPaymentRepository) UpdateStatus(ctx context.Context, tx db.Tx, id types.ID, status int8, values map[string]any) error {
	update := map[string]any{
		"status":     status,
		"updated_at": types.Now(),
	}
	for k, v := range values {
		update[k] = v
	}

	_, err := tx.Table(paymentTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(update)

	return err
}

// MarkPaying 把支付单置为支付中，并记录调起渠道后的交易信息。
func (r *PayPaymentRepository) MarkPaying(ctx context.Context, id types.ID, openID string) error {
	_, err := facades.DB().WithContext(ctx).Table(paymentTable).
		Where("id", id).
		Where("status", payEnums.PaymentStatusPending).
		Update(map[string]any{
			"status":     payEnums.PaymentStatusPaying,
			"open_id":    openID,
			"updated_at": types.Now(),
		})

	return err
}

// PaymentFilter 是后台支付单列表的筛选条件，零值表示不限制。
type PaymentFilter struct {
	PaymentNo string
	MemberID  types.ID
	Channel   string
	Status    *int8
}

// Paginate 按条件分页查询支付单，供后台支付流水页使用。
func (r *PayPaymentRepository) Paginate(ctx context.Context, filter PaymentFilter, pageNo, pageSize int) ([]model.PayPaymentModel, int64, error) {
	list := make([]model.PayPaymentModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(paymentTable).WhereNull("deleted_at")
	if filter.PaymentNo != "" {
		query = query.WhereLike("payment_no", "%"+filter.PaymentNo+"%")
	}
	if filter.MemberID > 0 {
		query = query.Where("member_user_id", filter.MemberID)
	}
	if filter.Channel != "" {
		query = query.Where("channel", filter.Channel)
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// FindExpired 查询已超时但仍处于待支付 / 支付中的支付单，供关单任务使用。
func (r *PayPaymentRepository) FindExpired(ctx context.Context, limit int) ([]model.PayPaymentModel, error) {
	list := make([]model.PayPaymentModel, 0)

	if err := facades.DB().WithContext(ctx).Table(paymentTable).
		WhereIn("status", []any{payEnums.PaymentStatusPending, payEnums.PaymentStatusPaying}).
		Where("expire_time", "<", types.Now()).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Limit(uint64(limit)).
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// MarkClosed 把支付单置为已关闭，仅在待支付 / 支付中时才生效，返回影响行数。
func (r *PayPaymentRepository) MarkClosed(ctx context.Context, id types.ID) (int64, error) {
	now := types.Now()

	res, err := facades.DB().WithContext(ctx).Table(paymentTable).
		Where("id", id).
		WhereIn("status", []any{payEnums.PaymentStatusPending, payEnums.PaymentStatusPaying}).
		WhereNull("deleted_at").
		Update(map[string]any{
			"status":     payEnums.PaymentStatusClosed,
			"close_time": now,
			"updated_at": now,
		})
	if err != nil {
		return 0, err
	}

	return res.RowsAffected, nil
}

// FindPaidByBiz 按业务定位已支付单，业务方（如打印订单）退款时用它拿 payment_no。
// 参数 bizType/bizID 与支付单创建时一致；只取已支付或部分退款的记录。
func (r *PayPaymentRepository) FindPaidByBiz(ctx context.Context, bizType string, bizID types.ID) (*model.PayPaymentModel, error) {
	var m model.PayPaymentModel

	if err := facades.DB().WithContext(ctx).Table(paymentTable).
		Where("biz_type", bizType).
		Where("biz_id", bizID).
		Where("status", payEnums.PaymentStatusPaid).
		WhereNull("deleted_at").
		OrderByDesc("id").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// LockByID 在事务内锁定支付单，退款时用于串行化同一支付单的并发退款。
func (r *PayPaymentRepository) LockByID(ctx context.Context, tx db.Tx, id types.ID) (*model.PayPaymentModel, error) {
	var m model.PayPaymentModel

	if err := tx.Table(paymentTable).
		Where("id", id).
		WhereNull("deleted_at").
		LockForUpdate().
		First(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// UpdateRefund 在事务内写回累计退款金额与支付单状态。
func (r *PayPaymentRepository) UpdateRefund(ctx context.Context, tx db.Tx, m *model.PayPaymentModel) error {
	_, err := tx.Table(paymentTable).
		Where("id", m.ID).
		Update(map[string]any{
			"refund_amount": m.RefundAmount,
			"status":        m.Status,
			"updated_at":    types.Now(),
		})

	return err
}

// PaginateByMember 按会员分页查询支付单，供 C 端账单列表使用。
func (r *PayPaymentRepository) PaginateByMember(ctx context.Context, memberID types.ID, pageNo, pageSize int) ([]model.PayPaymentModel, int64, error) {
	list := make([]model.PayPaymentModel, 0)
	var total int64

	if err := facades.DB().WithContext(ctx).Table(paymentTable).
		Where("member_user_id", memberID).
		WhereNull("deleted_at").
		OrderByDesc("id").
		Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
