package repository

import (
	"context"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/model"
)

const orderTable = model.PrintOrderTable

// OrderFilter 是订单列表的筛选条件，零值表示不限制。
type OrderFilter struct {
	ShopID      types.ID
	MemberID    types.ID
	OrderNo     string
	PayStatus   *int8
	OrderStatus *int8
}

// PrintOrderRepository 是 print_order 表的仓储。
type PrintOrderRepository struct{}

// NewPrintOrderRepository 构造订单仓储。
func NewPrintOrderRepository() *PrintOrderRepository {
	return &PrintOrderRepository{}
}

// CreateInTx 在事务内创建订单，与明细同生共死。
func (r *PrintOrderRepository) CreateInTx(ctx context.Context, tx db.Tx, m *model.PrintOrderModel) error {
	m.FillCreate(ctx)

	_, err := tx.Table(orderTable).Insert(map[string]any{
		"id":                 m.ID,
		"order_no":           m.OrderNo,
		"print_shop_id":      m.PrintShopID,
		"member_user_id":     m.MemberUserID,
		"device_id":          m.DeviceID,
		"pickup_code":        m.PickupCode,
		"pickup_expire_time": m.PickupExpireTime,
		"item_count":         m.ItemCount,
		"total_pages":        m.TotalPages,
		"total_sheets":       m.TotalSheets,
		"total_amount":       m.TotalAmount,
		"discount_percent":   m.DiscountPercent,
		"discount_amount":    m.DiscountAmount,
		"pay_amount":         m.PayAmount,
		"refund_amount":      m.RefundAmount,
		"pay_status":         m.PayStatus,
		"pay_type":           m.PayType,
		"order_status":       m.OrderStatus,
		"pay_expire_time":    m.PayExpireTime,
		"remark":             m.Remark,
		"creator":            m.Creator,
		"created_at":         m.CreatedAt,
		"updater":            m.Updater,
		"updated_at":         m.UpdatedAt,
	})

	return err
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintOrderRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintOrderModel, error) {
	var m model.PrintOrderModel

	if err := facades.DB().WithContext(ctx).Table(orderTable).
		Where("id", id).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindByNo 按订单号查询；未命中返回 (nil, nil)。
func (r *PrintOrderRepository) FindByNo(ctx context.Context, orderNo string) (*model.PrintOrderModel, error) {
	var m model.PrintOrderModel

	if err := facades.DB().WithContext(ctx).Table(orderTable).
		Where("order_no", orderNo).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// LockByID 在事务内锁定订单，用于状态推进。
func (r *PrintOrderRepository) LockByID(ctx context.Context, tx db.Tx, id types.ID) (*model.PrintOrderModel, error) {
	var m model.PrintOrderModel

	if err := tx.Table(orderTable).
		Where("id", id).
		WhereNull("deleted_at").
		LockForUpdate().
		First(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// MarkPaid 在事务内把待支付订单置为已支付并转「待打印」。
// 条件更新保证重复回调只生效一次，返回影响行数。
func (r *PrintOrderRepository) MarkPaid(ctx context.Context, tx db.Tx, id types.ID, payType, tradeNo string) (int64, error) {
	now := types.Now()

	res, err := tx.Table(orderTable).
		Where("id", id).
		Where("order_status", printEnums.OrderStatusUnpaid).
		Update(map[string]any{
			"pay_status":   printEnums.PayStatusPaid,
			"order_status": printEnums.OrderStatusPending,
			"pay_type":     payType,
			"trade_no":     tradeNo,
			"pay_time":     now,
			"updated_at":   now,
		})
	if err != nil {
		return 0, err
	}

	return res.RowsAffected, nil
}

// UpdateStatus 更新订单状态与附加字段。
func (r *PrintOrderRepository) UpdateStatus(ctx context.Context, id types.ID, status int8, values map[string]any) error {
	update := map[string]any{
		"order_status": status,
		"updated_at":   types.Now(),
	}
	for k, v := range values {
		update[k] = v
	}

	_, err := facades.DB().WithContext(ctx).Table(orderTable).
		Where("id", id).
		Update(update)

	return err
}

// ExistsPickupCode 判断门店内未完成订单是否已占用该取件码。
func (r *PrintOrderRepository) ExistsPickupCode(ctx context.Context, shopID types.ID, code string) (bool, error) {
	total, err := facades.DB().WithContext(ctx).Table(orderTable).
		Where("print_shop_id", shopID).
		Where("pickup_code", code).
		WhereIn("order_status", []any{
			printEnums.OrderStatusPending,
			printEnums.OrderStatusPrinting,
			printEnums.OrderStatusPickup,
		}).
		WhereNull("deleted_at").
		Count()
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// FindByPickup 按「门店 + 取件码」查询待取件订单，核销用。
func (r *PrintOrderRepository) FindByPickup(ctx context.Context, shopID types.ID, code string) (*model.PrintOrderModel, error) {
	var m model.PrintOrderModel

	if err := facades.DB().WithContext(ctx).Table(orderTable).
		Where("print_shop_id", shopID).
		Where("pickup_code", code).
		Where("order_status", printEnums.OrderStatusPickup).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindExpiredUnpaid 查询已过支付截止时间仍未支付的订单，供关单任务使用。
// 参数 limit 为单批上限，避免一次拉太多导致长事务。
func (r *PrintOrderRepository) FindExpiredUnpaid(ctx context.Context, limit int) ([]model.PrintOrderModel, error) {
	list := make([]model.PrintOrderModel, 0)

	if err := facades.DB().WithContext(ctx).Table(orderTable).
		Where("order_status", printEnums.OrderStatusUnpaid).
		Where("pay_expire_time", "<", types.Now()).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Limit(uint64(limit)).
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// Paginate 按条件分页查询订单，会员端与后台共用。
func (r *PrintOrderRepository) Paginate(ctx context.Context, filter OrderFilter, pageNo, pageSize int) ([]model.PrintOrderModel, int64, error) {
	list := make([]model.PrintOrderModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(orderTable).WhereNull("deleted_at")
	if filter.ShopID > 0 {
		query = query.Where("print_shop_id", filter.ShopID)
	}
	if filter.MemberID > 0 {
		query = query.Where("member_user_id", filter.MemberID)
	}
	if filter.OrderNo != "" {
		query = query.WhereLike("order_no", "%"+filter.OrderNo+"%")
	}
	if filter.PayStatus != nil {
		query = query.Where("pay_status", *filter.PayStatus)
	}
	if filter.OrderStatus != nil {
		query = query.Where("order_status", *filter.OrderStatus)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
