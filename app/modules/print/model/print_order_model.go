package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintOrderTable 是打印订单表名常量。
const PrintOrderTable = "print_order"

// PrintOrderModel 对应 print_order，到店打印订单。
// 支付成功即自动打印，打印完成转待取件，取件码核销后完成。
type PrintOrderModel struct {
	baseModel.BaseModel
	OrderNo            string          `db:"order_no" json:"orderNo"`
	PrintShopID        types.ID        `db:"print_shop_id" json:"printShopId"`
	MemberUserID       types.ID        `db:"member_user_id" json:"memberUserId"`
	DeviceID           types.ID        `db:"device_id" json:"deviceId"`
	PickupCode         string          `db:"pickup_code" json:"pickupCode"`
	PickupExpireTime   *types.DateTime `db:"pickup_expire_time" json:"pickupExpireTime"`
	PickupVerifyAt     *types.DateTime `db:"pickup_verify_at" json:"pickupVerifyAt"`
	PickupVerifyUserID types.ID        `db:"pickup_verify_user_id" json:"pickupVerifyUserId"`
	ItemCount          int             `db:"item_count" json:"itemCount"`
	TotalPages         int             `db:"total_pages" json:"totalPages"`
	TotalSheets        int             `db:"total_sheets" json:"totalSheets"`
	TotalAmount        int64           `db:"total_amount" json:"totalAmount"`
	DiscountPercent    int             `db:"discount_percent" json:"discountPercent"`
	DiscountAmount     int64           `db:"discount_amount" json:"discountAmount"`
	PayAmount          int64           `db:"pay_amount" json:"payAmount"`
	RefundAmount       int64           `db:"refund_amount" json:"refundAmount"`
	PayStatus          int8            `db:"pay_status" json:"payStatus"`
	PayType            string          `db:"pay_type" json:"payType"`
	TradeNo            string          `db:"trade_no" json:"tradeNo"`
	PayTime            *types.DateTime `db:"pay_time" json:"payTime"`
	OrderStatus        int8            `db:"order_status" json:"orderStatus"`
	PayExpireTime      *types.DateTime `db:"pay_expire_time" json:"payExpireTime"`
	FinishTime         *types.DateTime `db:"finish_time" json:"finishTime"`
	CancelTime         *types.DateTime `db:"cancel_time" json:"cancelTime"`
	CancelReason       string          `db:"cancel_reason" json:"cancelReason"`
	CloseTime          *types.DateTime `db:"close_time" json:"closeTime"`
	Remark             string          `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PrintOrderModel) TableName() string {
	return PrintOrderTable
}
