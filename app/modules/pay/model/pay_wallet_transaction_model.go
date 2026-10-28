package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PayWalletTransactionTable 是钱包流水表名常量。
const PayWalletTransactionTable = "pay_wallet_transaction"

// PayWalletTransactionModel 对应 pay_wallet_transaction。
// (BizType, BizID, Type) 是幂等键：同一业务重复入账会撞唯一键，
// 服务层捕获重复键即视为「已处理」，不得重复加钱。
type PayWalletTransactionModel struct {
	baseModel.BaseModel
	WalletID     types.ID `db:"wallet_id" json:"walletId"`
	MemberUserID types.ID `db:"member_user_id" json:"memberUserId"`
	BizType      string   `db:"biz_type" json:"bizType"`
	BizID        types.ID `db:"biz_id" json:"bizId"`
	Type         int8     `db:"type" json:"type"`
	Amount       int64    `db:"amount" json:"amount"`
	BalanceAfter int64    `db:"balance_after" json:"balanceAfter"`
	Title        string   `db:"title" json:"title"`
	Remark       string   `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PayWalletTransactionModel) TableName() string {
	return PayWalletTransactionTable
}
