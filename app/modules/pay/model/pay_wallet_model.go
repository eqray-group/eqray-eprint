package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PayWalletTable 是会员钱包表名常量。
const PayWalletTable = "pay_wallet"

// PayWalletModel 对应 pay_wallet，一个会员一个钱包。
// Balance 只能与流水同事务变更，扣减时必须带 balance >= amount 条件。
type PayWalletModel struct {
	baseModel.BaseModel
	MemberUserID types.ID `db:"member_user_id" json:"memberUserId"`
	Balance      int64    `db:"balance" json:"balance"`
	TotalIncome  int64    `db:"total_income" json:"totalIncome"`
	TotalExpense int64    `db:"total_expense" json:"totalExpense"`
	Status       int8     `db:"status" json:"status"`
	Version      int      `db:"version" json:"version"`
}

// TableName 返回表名。
func (m *PayWalletModel) TableName() string {
	return PayWalletTable
}
