package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000014CreatePrintPrinterTable 建 print_printer 打印机表。
// 一条记录 = 某台设备上的一个 Windows 打印队列：队列装在具体机器上，
// 同型号打印机在两台电脑各装一份就是两条记录。
// 只收录 Windows 能识别的打印机，客户端按 system_name 指定队列静默打印，
// 不直连 IPP / 厂商 SDK；不设默认打印机，路由只按配置的参数能力匹配。
type M20261027000014CreatePrintPrinterTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000014CreatePrintPrinterTable) Signature() string {
	return "20261027000014_create_print_printer_table"
}

// Up 建打印机表，uk_device_system_name 保证同一设备内队列名唯一。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000014CreatePrintPrinterTable) Up() error {
	if facades.Schema().HasTable("print_printer") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_printer (
  id             bigint       NOT NULL COMMENT '主键',
  print_shop_id  bigint       NOT NULL COMMENT '门店 ID（冗余，便于按店查询）',
  device_id      bigint       NOT NULL COMMENT '所属设备 ID，对应 print_device.id；打印队列装在这台 Windows 机器上',
  name           varchar(64)  NOT NULL COMMENT '打印机别名，如 A4 黑白 1 号机（后台展示用）',
  system_name    varchar(128) NOT NULL COMMENT 'Windows 打印队列名称（spooler 的 PrinterName），同一设备内唯一，静默打印按此指定',
  model          varchar(64)  NOT NULL DEFAULT '' COMMENT '型号',
  paper_types    varchar(64)  NOT NULL DEFAULT 'A4' COMMENT '本机支持的纸张，逗号分隔：A4,A3，由门店后台勾选配置',
  color_support  tinyint      NOT NULL DEFAULT 0 COMMENT '是否支持彩色：0 否 / 1 是，彩色任务遇 0 必须失败',
  duplex_support tinyint      NOT NULL DEFAULT 0 COMMENT '是否支持双面：0 否 / 1 是',
  status         tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 停用 / 1 启用',
  remark         varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  creator        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at     datetime     NULL COMMENT '创建时间',
  updater        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at     datetime     NULL COMMENT '更新时间',
  deleted_at     datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_device_system_name (device_id, system_name),
  KEY idx_shop (print_shop_id),
  KEY idx_device (device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备打印机表'`)
}

// Down 删除打印机表，无参数，返回删除结果。
func (r *M20261027000014CreatePrintPrinterTable) Down() error {
	return facades.Schema().DropIfExists("print_printer")
}
