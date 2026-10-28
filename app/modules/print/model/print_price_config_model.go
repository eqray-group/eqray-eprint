package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintPriceConfigTable 是打印价格配置表名常量。
const PrintPriceConfigTable = "print_price_config"

// PrintPriceConfigModel 对应 print_price_config。
// 门店 × 纸张 × 色彩 × 单双面 唯一确定一个单价（分/页）；
// 表无 deleted_at，规则删除走物理删除。
type PrintPriceConfigModel struct {
	baseModel.BaseModel
	PrintShopID types.ID `db:"print_shop_id" json:"printShopId"`
	PaperType   string   `db:"paper_type" json:"paperType"`
	ColorType   string   `db:"color_type" json:"colorType"`
	PrintMode   string   `db:"print_mode" json:"printMode"`
	Price       int64    `db:"price" json:"price"`
	Status      int8     `db:"status" json:"status"`
}

// TableName 返回表名。
func (m *PrintPriceConfigModel) TableName() string {
	return PrintPriceConfigTable
}
