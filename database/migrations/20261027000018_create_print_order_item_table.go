package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000018CreatePrintOrderItemTable 建 print_order_item 订单明细表。
// 一条明细同时承担「打印任务」：一个文件 + 一套参数 + 一台打印机。
// task_no 全局唯一，客户端据此去重防重复打印；
// 打印参数以本表为唯一来源，print_config 只是全量快照，客户端与门店后台不得覆盖。
type M20261027000018CreatePrintOrderItemTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000018CreatePrintOrderItemTable) Signature() string {
	return "20261027000018_create_print_order_item_table"
}

// Up 建订单明细表；按订单、设备状态、门店状态与文件建索引。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000018CreatePrintOrderItemTable) Up() error {
	if facades.Schema().HasTable("print_order_item") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_order_item (
  id              bigint       NOT NULL COMMENT '主键',
  print_shop_id   bigint       NOT NULL COMMENT '门店 ID（冗余，便于按店统计）',
  order_id        bigint       NOT NULL COMMENT '订单 ID，对应 print_order.id',
  task_no         varchar(64)  NOT NULL COMMENT '任务号，全局唯一，客户端防重复打印依据',
  file_id         bigint       NOT NULL COMMENT '打印文件 ID，对应 print_file.id',
  file_name       varchar(255) NOT NULL COMMENT '文件名快照',
  file_page_count int          NOT NULL DEFAULT 0 COMMENT '文件总页数快照',
  paper_type      varchar(16)  NOT NULL DEFAULT 'A4' COMMENT '纸张类型：A4 / A3',
  color_type      varchar(16)  NOT NULL DEFAULT 'bw' COMMENT '色彩：bw 黑白 / color 彩色',
  print_mode      varchar(16)  NOT NULL DEFAULT 'single' COMMENT '单双面：single 单面 / double 双面',
  orientation     varchar(16)  NOT NULL DEFAULT 'portrait' COMMENT '纸张方向：portrait 纵向 / landscape 横向',
  page_range      varchar(128) NOT NULL DEFAULT '' COMMENT '页码范围，如 1-5,8,10-12；空表示全部',
  page_count      int          NOT NULL DEFAULT 0 COMMENT '本明细计费页数（按页码范围计算）',
  copies          int          NOT NULL DEFAULT 1 COMMENT '打印份数',
  sheet_count     int          NOT NULL DEFAULT 0 COMMENT '印张数 = page_count × copies',
  print_config    json         NULL COMMENT '打印参数全量快照，服务端下发、客户端只读',
  amount          bigint       NOT NULL DEFAULT 0 COMMENT '小计金额（分）',
  status          tinyint      NOT NULL DEFAULT 0 COMMENT '任务状态：0 待分配 / 1 已分配 / 2 打印中 / 3 已完成 / 4 失败 / 5 已取消',
  device_id       bigint       NOT NULL DEFAULT 0 COMMENT '执行设备 ID，0 表示未分配',
  printer_id      bigint       NOT NULL DEFAULT 0 COMMENT '执行打印机 ID，0 表示未分配；创建后锁定，换机需服务端重建任务',
  retry_count     int          NOT NULL DEFAULT 0 COMMENT '重试次数，上限 3，超限转人工并触发订单退款',
  error_code      varchar(32)  NOT NULL DEFAULT '' COMMENT '失败错误码，如 PAPER_JAM / COLOR_UNSUPPORTED / FILE_BROKEN',
  fail_reason     varchar(255) NOT NULL DEFAULT '' COMMENT '失败原因',
  printed_pages   int          NOT NULL DEFAULT 0 COMMENT '实际打印页数（客户端回执）',
  started_at      datetime     NULL COMMENT '开始打印时间',
  finished_at     datetime     NULL COMMENT '打印结束时间',
  creator         varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at      datetime     NULL COMMENT '创建时间',
  updater         varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at      datetime     NULL COMMENT '更新时间',
  deleted_at      datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_task_no (task_no),
  KEY idx_order (order_id),
  KEY idx_device_status (device_id, status),
  KEY idx_shop_status (print_shop_id, status),
  KEY idx_file (file_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打印订单明细表'`)
}

// Down 删除订单明细表，无参数，返回删除结果。
func (r *M20261027000018CreatePrintOrderItemTable) Down() error {
	return facades.Schema().DropIfExists("print_order_item")
}
