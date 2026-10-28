package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000013CreatePrintDeviceTable 建 print_device 打印设备表。
// 只支持 Windows 客户端：客户端由店主 / 店员用 system_users 账号登录，
// 没有独立设备令牌，设备归属门店取自 print_shop_staff 绑定。
// 在线状态由心跳推导，online_status 是冗余字段便于统计在线率。
type M20261027000013CreatePrintDeviceTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000013CreatePrintDeviceTable) Signature() string {
	return "20261027000013_create_print_device_table"
}

// Up 建打印设备表，code（机器指纹）全局唯一，一台机器只登记一次。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000013CreatePrintDeviceTable) Up() error {
	if facades.Schema().HasTable("print_device") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_device (
  id                 bigint       NOT NULL COMMENT '主键',
  print_shop_id      bigint       NOT NULL COMMENT '门店 ID',
  code               varchar(64)  NOT NULL COMMENT '设备唯一编码（客户端机器指纹），全局唯一',
  name               varchar(128) NOT NULL COMMENT '设备名称',
  client_version     varchar(32)  NOT NULL DEFAULT '' COMMENT 'Windows 客户端版本号，用于强制升级校验',
  status             tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 禁用 / 1 启用',
  online_status      tinyint      NOT NULL DEFAULT 0 COMMENT '在线状态：0 离线 / 1 在线（心跳超过 90 秒由服务端置 0）',
  ip                 varchar(50)  NOT NULL DEFAULT '' COMMENT '最近上报 IP',
  last_heartbeat_at  datetime     NULL COMMENT '最后心跳时间',
  last_login_user_id bigint       NOT NULL DEFAULT 0 COMMENT '最近登录该设备的后台账号，对应 system_users.id',
  last_error_info    varchar(512) NOT NULL DEFAULT '' COMMENT '最后错误信息',
  remark             varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  creator            varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at         datetime     NULL COMMENT '创建时间',
  updater            varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at         datetime     NULL COMMENT '更新时间',
  deleted_at         datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_code (code),
  KEY idx_shop (print_shop_id),
  KEY idx_online_status (online_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打印设备表'`)
}

// Down 删除打印设备表，无参数，返回删除结果。
func (r *M20261027000013CreatePrintDeviceTable) Down() error {
	return facades.Schema().DropIfExists("print_device")
}
