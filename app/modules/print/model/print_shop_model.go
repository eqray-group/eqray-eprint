package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
)

// PrintShopTable 是打印门店表名常量。
const PrintShopTable = "print_shop"

// PrintShopModel 对应 print_shop，到店打印的营业门店。
// status 只表示营业/停业，审核状态单独用 audit_status，二者互不覆盖。
type PrintShopModel struct {
	baseModel.BaseModel
	Name      string  `db:"name" json:"name"`
	Address   string  `db:"address" json:"address"`
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Avatar    string  `db:"avatar" json:"avatar"`
	// BusinessHours 为营业时段 JSON，NULL 表示全天营业
	BusinessHours *string `db:"business_hours" json:"businessHours"`
	Status        int8    `db:"status" json:"status"`
	AuditStatus   int8    `db:"audit_status" json:"auditStatus"`
	AuditRemark   string  `db:"audit_remark" json:"auditRemark"`
	Remark        string  `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PrintShopModel) TableName() string {
	return PrintShopTable
}
