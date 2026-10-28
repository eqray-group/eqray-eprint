package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintShopStaffTable 是门店成员绑定表名常量。
const PrintShopStaffTable = "print_shop_staff"

// PrintShopStaffModel 对应 print_shop_staff，后台账号与门店的多对多绑定。
// 表无 deleted_at：解绑走物理删除，避免唯一键挡住重新绑定。
type PrintShopStaffModel struct {
	baseModel.BaseModel
	UserID types.ID `db:"user_id" json:"userId"`
	ShopID types.ID `db:"shop_id" json:"shopId"`
	Role   string   `db:"role" json:"role"`
	Status int8     `db:"status" json:"status"`
}

// TableName 返回表名。
func (m *PrintShopStaffModel) TableName() string {
	return PrintShopStaffTable
}
