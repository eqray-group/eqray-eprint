package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintDeviceTable 是打印设备表名常量。
const PrintDeviceTable = "print_device"

// PrintDeviceModel 对应 print_device，安装在门店的 Windows 客户端实例。
// 客户端用后台账号登录，设备归属门店取该账号在 print_shop_staff 的绑定。
type PrintDeviceModel struct {
	baseModel.BaseModel
	PrintShopID     types.ID        `db:"print_shop_id" json:"printShopId"`
	Code            string          `db:"code" json:"code"`
	Name            string          `db:"name" json:"name"`
	ClientVersion   string          `db:"client_version" json:"clientVersion"`
	Status          int8            `db:"status" json:"status"`
	OnlineStatus    int8            `db:"online_status" json:"onlineStatus"`
	IP              string          `db:"ip" json:"ip"`
	LastHeartbeatAt *types.DateTime `db:"last_heartbeat_at" json:"lastHeartbeatAt"`
	LastLoginUserID types.ID        `db:"last_login_user_id" json:"lastLoginUserId"`
	LastErrorInfo   string          `db:"last_error_info" json:"lastErrorInfo"`
	Remark          string          `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PrintDeviceModel) TableName() string {
	return PrintDeviceTable
}
