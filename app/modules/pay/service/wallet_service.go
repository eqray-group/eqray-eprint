package service

import (
	"context"
	"fmt"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	payEnums "eqray-eprint/app/modules/pay/enums"
	payErrs "eqray-eprint/app/modules/pay/errs"
	"eqray-eprint/app/modules/pay/model"
	"eqray-eprint/app/modules/pay/repository"
)

// WalletService 负责会员钱包的查询与余额变更。
// 余额变更必须在外层事务里调用 Credit / Debit，保证与业务单状态同生共死。
type WalletService struct {
	wallets *repository.PayWalletRepository
}

// NewWalletService 构造钱包服务。
func NewWalletService() *WalletService {
	return &WalletService{wallets: repository.NewPayWalletRepository()}
}

// Ensure 取会员钱包，不存在则创建（懒创建，注册流程不需要专门建钱包）。
// 并发重复创建由 uk_member 兜底，撞键后回查即可。
func (s *WalletService) Ensure(ctx context.Context, memberID types.ID) (*model.PayWalletModel, error) {
	wallet, err := s.wallets.FindByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		return wallet, nil
	}

	wallet = &model.PayWalletModel{
		MemberUserID: memberID,
		Status:       payEnums.WalletStatusNormal,
	}
	wallet.ID = support.NewID()

	if err := s.wallets.Create(ctx, wallet); err != nil {
		if support.IsDuplicate(err) {
			return s.wallets.FindByMember(ctx, memberID)
		}

		return nil, err
	}

	return wallet, nil
}

// AdminTransactionQuery 是后台钱包流水列表的查询条件。
type AdminTransactionQuery struct {
	MemberID types.ID
	Type     *int8
	BizType  string
	PageNo   int
	PageSize int
}

// AdminDetail 后台查询会员钱包。
// 与 Detail 不同，这里不懒创建：查询动作不应该产生数据，没有钱包直接返回空余额。
func (s *WalletService) AdminDetail(ctx context.Context, memberID types.ID) (*model.PayWalletModel, error) {
	wallet, err := s.wallets.FindByMember(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		return wallet, nil
	}

	return &model.PayWalletModel{MemberUserID: memberID}, nil
}

// AdminTransactions 后台分页查询钱包流水。
func (s *WalletService) AdminTransactions(ctx context.Context, q AdminTransactionQuery) ([]model.PayWalletTransactionModel, int64, error) {
	return s.wallets.PaginateTransactions(ctx, repository.WalletTransactionFilter{
		MemberID: q.MemberID,
		Type:     q.Type,
		BizType:  q.BizType,
	}, q.PageNo, q.PageSize)
}

// Detail 返回会员钱包，没有则创建一个空钱包。
func (s *WalletService) Detail(ctx context.Context, memberID types.ID) (*model.PayWalletModel, error) {
	return s.Ensure(ctx, memberID)
}

// Transactions 分页查询会员钱包流水。
func (s *WalletService) Transactions(ctx context.Context, memberID types.ID, pageNo, pageSize int) ([]model.PayWalletTransactionModel, int64, error) {
	return s.wallets.Transactions(ctx, memberID, pageNo, pageSize)
}

// Adjust 后台人工调账：amount 为正是加款、为负是扣款，返回调账后余额。
// 每次调账用新的雪花 ID 作为 biz_id，因此可以重复调账而不会被幂等键拦住。
func (s *WalletService) Adjust(ctx context.Context, memberID types.ID, amount int64, reason string) (int64, error) {
	if amount == 0 {
		return 0, baseErrs.New(payErrs.PayWalletAmountInvalid)
	}

	title := "后台调账"
	if reason != "" {
		title = "后台调账：" + reason
	}

	var balance int64
	err := facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		var err error
		if amount > 0 {
			balance, err = s.Credit(ctx, tx, memberID, payEnums.BizTypeManualAdjust, support.NewID(), amount, title)
		} else {
			balance, err = s.Debit(ctx, tx, memberID, payEnums.BizTypeManualAdjust, support.NewID(), -amount, title)
		}

		return err
	})

	return balance, err
}

// Credit 在事务内给会员钱包入账并写流水，返回入账后余额。
// 已存在同业务同方向的流水时直接返回，保证回调重试不会重复加钱。
func (s *WalletService) Credit(ctx context.Context, tx db.Tx, memberID types.ID, bizType string, bizID types.ID, amount int64, title string) (int64, error) {
	return s.change(ctx, tx, memberID, bizType, bizID, amount, payEnums.WalletTxIncome, title)
}

// Debit 在事务内从会员钱包扣款并写流水，返回扣款后余额。
// 余额不足由仓储层的条件更新拦住，直接返回业务错误。
func (s *WalletService) Debit(ctx context.Context, tx db.Tx, memberID types.ID, bizType string, bizID types.ID, amount int64, title string) (int64, error) {
	return s.change(ctx, tx, memberID, bizType, bizID, amount, payEnums.WalletTxExpense, title)
}

// change 是入账与扣款的公共实现：幂等预检 → 改余额 → 写流水。
func (s *WalletService) change(ctx context.Context, tx db.Tx, memberID types.ID, bizType string, bizID types.ID, amount int64, direction int8, title string) (int64, error) {
	wallet, err := s.Ensure(ctx, memberID)
	if err != nil {
		return 0, err
	}

	exists, err := s.wallets.ExistsTransaction(ctx, tx, bizType, bizID, direction)
	if err != nil {
		return 0, err
	}
	if exists {
		// 已入账，直接返回当前余额，不重复变更
		return wallet.Balance, nil
	}

	var balance int64
	if direction == payEnums.WalletTxIncome {
		balance, err = s.wallets.Increase(ctx, tx, wallet.ID, amount)
	} else {
		balance, err = s.wallets.Decrease(ctx, tx, wallet.ID, amount)
	}
	if err != nil {
		return 0, err
	}

	record := &model.PayWalletTransactionModel{
		WalletID:     wallet.ID,
		MemberUserID: memberID,
		BizType:      bizType,
		BizID:        bizID,
		Type:         direction,
		Amount:       amount,
		BalanceAfter: balance,
		Title:        title,
	}
	record.ID = support.NewID()

	if err := s.wallets.CreateTransaction(ctx, tx, record); err != nil {
		return 0, err
	}

	return balance, nil
}

// balanceTitle 生成流水展示文案，便于 C 端账单直接展示。
func balanceTitle(bizType string, amount int64) string {
	switch bizType {
	case payEnums.BizTypeWalletRecharge:
		return fmt.Sprintf("钱包充值 %.2f 元", float64(amount)/100)
	case payEnums.BizTypePrintRefund:
		return fmt.Sprintf("打印订单退款 %.2f 元", float64(amount)/100)
	default:
		return fmt.Sprintf("打印消费 %.2f 元", float64(amount)/100)
	}
}
