package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintOrderItemTable 是打印订单明细表名常量。
const PrintOrderItemTable = "print_order_item"

// PrintOrderItemModel 对应 print_order_item，一条明细 = 一个文件 + 一套参数 + 一台打印机，
// 同时承担「打印任务」职责。task_no 全局唯一，客户端据此防重复打印。
type PrintOrderItemModel struct {
	baseModel.BaseModel
	PrintShopID   types.ID `db:"print_shop_id" json:"printShopId"`
	OrderID       types.ID `db:"order_id" json:"orderId"`
	TaskNo        string   `db:"task_no" json:"taskNo"`
	FileID        types.ID `db:"file_id" json:"fileId"`
	FileName      string   `db:"file_name" json:"fileName"`
	FilePageCount int      `db:"file_page_count" json:"filePageCount"`
	PaperType     string   `db:"paper_type" json:"paperType"`
	ColorType     string   `db:"color_type" json:"colorType"`
	PrintMode     string   `db:"print_mode" json:"printMode"`
	Orientation   string   `db:"orientation" json:"orientation"`
	PageRange     string   `db:"page_range" json:"pageRange"`
	PageCount     int      `db:"page_count" json:"pageCount"`
	Copies        int      `db:"copies" json:"copies"`
	SheetCount    int      `db:"sheet_count" json:"sheetCount"`
	// PrintConfig 为打印参数全量快照，服务端下发、客户端只读
	PrintConfig  *string         `db:"print_config" json:"printConfig"`
	Amount       int64           `db:"amount" json:"amount"`
	Status       int8            `db:"status" json:"status"`
	DeviceID     types.ID        `db:"device_id" json:"deviceId"`
	PrinterID    types.ID        `db:"printer_id" json:"printerId"`
	RetryCount   int             `db:"retry_count" json:"retryCount"`
	ErrorCode    string          `db:"error_code" json:"errorCode"`
	FailReason   string          `db:"fail_reason" json:"failReason"`
	PrintedPages int             `db:"printed_pages" json:"printedPages"`
	StartedAt    *types.DateTime `db:"started_at" json:"startedAt"`
	FinishedAt   *types.DateTime `db:"finished_at" json:"finishedAt"`
}

// TableName 返回表名。
func (m *PrintOrderItemModel) TableName() string {
	return PrintOrderItemTable
}
