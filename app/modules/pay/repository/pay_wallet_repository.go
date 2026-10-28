package repository

import (
	"context"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/types"
	payEnums "eqray-eprint/app/modules/pay/enums"
	payErrs "eqray-eprint/app/modules/pay/errs"
	"eqray-eprint/app/modules/pay/model"
)

const walletTable = model.PayWalletTable

// PayWalletRepository 是 pay_wallet 与 pay_wallet_transaction 的仓储。
// 余额变更与流水写入必须在同一事务内完成，因此变动方法都接收 tx。
type PayWalletRepository struct{}

// NewPayWalletRepository 构造钱包仓储。
func NewPayWalletRepository() *PayWalletRepository {
	return &PayWalletRepository{}
}

// Create 新建钱包账户，会员首次访问钱包时懒创建。
func (r *PayWalletRepository) Create(ctx context.Context, m *model.PayWalletModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(walletTable).Insert(map[string]any{
		"id":             m.ID,
		"member_user_id": m.MemberUserID,
		"balance":        m.Balance,
		"total_income":   m.TotalIncome,
		"total_expense":  m.TotalExpense,
		"status":         m.Status,
		"version":        m.Version,
		"creator":        m.Creator,
		"created_at":     m.CreatedAt,
		"updater":        m.Updater,
		"updated_at":     m.UpdatedAt,
	})

	return err
}

// FindByMember 按会员查询钱包；未命中返回 (nil, nil)。
func (r *PayWalletRepository) FindByMember(ctx context.Context, memberID types.ID) (*model.PayWalletModel, error) {
	var m model.PayWalletModel

	if err := facades.DB().WithContext(ctx).Table(walletTable).
		Where("member_user_id", memberID).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// Increase 在事务内增加余额，返回变动后的余额。
// 先 SELECT ... FOR UPDATE 锁行再写绝对值，避免「读-改-写」并发丢更新。
func (r *PayWalletRepository) Increase(ctx context.Context, tx db.Tx, walletID types.ID, amount int64) (int64, error) {
	wallet, err := r.lock(ctx, tx, walletID)
	if err != nil {
		return 0, err
	}
	if wallet.Status != payEnums.WalletStatusNormal {
		return 0, baseErrs.New(payErrs.PayWalletDisabled)
	}

	balance := wallet.Balance + amount
	now := types.Now()

	_, err = tx.Table(walletTable).
		Where("id", walletID).
		Update(map[string]any{
			"balance":      balance,
			"total_income": wallet.TotalIncome + amount,
			"version":      wallet.Version + 1,
			"updated_at":   now,
		})
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// Decrease 在事务内扣减余额，返回变动后的余额。
// 余额不足或钱包冻结时返回业务错误，调用方据此回滚整个事务。
func (r *PayWalletRepository) Decrease(ctx context.Context, tx db.Tx, walletID types.ID, amount int64) (int64, error) {
	wallet, err := r.lock(ctx, tx, walletID)
	if err != nil {
		return 0, err
	}
	if wallet.Status != payEnums.WalletStatusNormal {
		return 0, baseErrs.New(payErrs.PayWalletDisabled)
	}
	if wallet.Balance < amount {
		return 0, baseErrs.New(payErrs.PayWalletBalanceNotEnough)
	}

	balance := wallet.Balance - amount
	now := types.Now()

	_, err = tx.Table(walletTable).
		Where("id", walletID).
		Update(map[string]any{
			"balance":       balance,
			"total_expense": wallet.TotalExpense + amount,
			"version":       wallet.Version + 1,
			"updated_at":    now,
		})
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// CreateTransaction 写入一条钱包流水；唯一键冲突说明业务已入账，由上层判定幂等。
func (r *PayWalletRepository) CreateTransaction(ctx context.Context, tx db.Tx, m *model.PayWalletTransactionModel) error {
	m.FillCreate(ctx)

	_, err := tx.Table(model.PayWalletTransactionTable).Insert(map[string]any{
		"id":             m.ID,
		"wallet_id":      m.WalletID,
		"member_user_id": m.MemberUserID,
		"biz_type":       m.BizType,
		"biz_id":         m.BizID,
		"type":           m.Type,
		"amount":         m.Amount,
		"balance_after":  m.BalanceAfter,
		"title":          m.Title,
		"remark":         m.Remark,
		"creator":        m.Creator,
		"created_at":     m.CreatedAt,
		"updater":        m.Updater,
		"updated_at":     m.UpdatedAt,
	})

	return err
}

// ExistsTransaction 在事务内判断业务流水是否已存在，用于入账前幂等预检。
// 参数 bizType/bizID/direction 与 (biz_type, biz_id, type) 唯一键一致。
func (r *PayWalletRepository) ExistsTransaction(ctx context.Context, tx db.Tx, bizType string, bizID types.ID, direction int8) (bool, error) {
	total, err := tx.Table(model.PayWalletTransactionTable).
		Where("biz_type", bizType).
		Where("biz_id", bizID).
		Where("type", direction).
		Count()
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// Transactions 按会员分页查询流水，最新的在前。
func (r *PayWalletRepository) Transactions(ctx context.Context, memberID types.ID, pageNo, pageSize int) ([]model.PayWalletTransactionModel, int64, error) {
	list := make([]model.PayWalletTransactionModel, 0)
	var total int64

	if err := facades.DB().WithContext(ctx).Table(model.PayWalletTransactionTable).
		Where("member_user_id", memberID).
		OrderByDesc("id").
		Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// WalletTransactionFilter 是后台钱包流水的筛选条件，零值表示不限制。
type WalletTransactionFilter struct {
	MemberID types.ID
	// Type 为 1 收入 / 2 支出，nil 表示不限
	Type    *int8
	BizType string
}

// PaginateTransactions 后台分页查询钱包流水，支持按会员、收支方向与业务类型筛选。
func (r *PayWalletRepository) PaginateTransactions(ctx context.Context, filter WalletTransactionFilter, pageNo, pageSize int) ([]model.PayWalletTransactionModel, int64, error) {
	list := make([]model.PayWalletTransactionModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(model.PayWalletTransactionTable).WhereNull("deleted_at")
	if filter.MemberID > 0 {
		query = query.Where("member_user_id", filter.MemberID)
	}
	if filter.Type != nil {
		query = query.Where("type", *filter.Type)
	}
	if filter.BizType != "" {
		query = query.Where("biz_type", filter.BizType)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// lock 在事务内锁定钱包行，未找到返回业务错误。
func (r *PayWalletRepository) lock(ctx context.Context, tx db.Tx, walletID types.ID) (*model.PayWalletModel, error) {
	var m model.PayWalletModel

	if err := tx.Table(walletTable).
		Where("id", walletID).
		WhereNull("deleted_at").
		LockForUpdate().
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, baseErrs.New(payErrs.PayWalletNotExists, walletID)
	}

	return &m, nil
}
