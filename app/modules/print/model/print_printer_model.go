package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintPrinterTable 是设备打印机表名常量。
const PrintPrinterTable = "print_printer"

// PrintPrinterModel 对应 print_printer，一条记录 = 某台设备上的一个 Windows 打印队列。
// 能力字段（paperTypes / colorSupport / duplexSupport）由门店后台配置，
// 路由只按这些字段匹配，不设默认打印机、不降级打印。
type PrintPrinterModel struct {
	baseModel.BaseModel
	PrintShopID   types.ID `db:"print_shop_id" json:"printShopId"`
	DeviceID      types.ID `db:"device_id" json:"deviceId"`
	Name          string   `db:"name" json:"name"`
	SystemName    string   `db:"system_name" json:"systemName"`
	Model         string   `db:"model" json:"model"`
	PaperTypes    string   `db:"paper_types" json:"paperTypes"`
	ColorSupport  int8     `db:"color_support" json:"colorSupport"`
	DuplexSupport int8     `db:"duplex_support" json:"duplexSupport"`
	Status        int8     `db:"status" json:"status"`
	Remark        string   `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PrintPrinterModel) TableName() string {
	return PrintPrinterTable
}
