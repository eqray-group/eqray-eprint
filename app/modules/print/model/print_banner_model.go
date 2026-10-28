package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintBannerTable 是门店轮播图表名常量。
const PrintBannerTable = "print_banner"

// PrintBannerModel 对应 print_banner，print_shop_id = 0 表示全平台通用。
type PrintBannerModel struct {
	baseModel.BaseModel
	PrintShopID types.ID `db:"print_shop_id" json:"printShopId"`
	Title       string   `db:"title" json:"title"`
	ImageURL    string   `db:"image_url" json:"imageUrl"`
	LinkType    string   `db:"link_type" json:"linkType"`
	LinkURL     string   `db:"link_url" json:"linkUrl"`
	Sort        int      `db:"sort" json:"sort"`
	Status      int8     `db:"status" json:"status"`
}

// TableName 返回表名。
func (m *PrintBannerModel) TableName() string {
	return PrintBannerTable
}
