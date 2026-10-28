package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000015CreatePrintPriceConfigTable 建 print_price_config 打印价格配置表。
// 门店 × 纸张 × 色彩 × 单双面 唯一确定一个单价（分/页）；
// 价格只是「卖什么」，能不能真印出来由 print_printer 的能力决定，
// 两者在用户选参数与下单时实时求交集，本表不关联打印机。
type M20261027000015CreatePrintPriceConfigTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000015CreatePrintPriceConfigTable) Signature() string {
	return "20261027000015_create_print_price_config_table"
}

// Up 建价格配置表，uk_rule 防止同一参数组合配出两条价格。
// 价格规则属于可重建数据，删除走物理删除，不建 deleted_at。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000015CreatePrintPriceConfigTable) Up() error {
	if facades.Schema().HasTable("print_price_config") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_price_config (
  id            bigint      NOT NULL COMMENT '主键',
  print_shop_id bigint      NOT NULL COMMENT '门店 ID',
  paper_type    varchar(16) NOT NULL COMMENT '纸张类型：A4 / A3',
  color_type    varchar(16) NOT NULL COMMENT '色彩：bw 黑白 / color 彩色',
  print_mode    varchar(16) NOT NULL COMMENT '单双面：single 单面 / double 双面',
  price         bigint      NOT NULL COMMENT '单价（分/页）',
  status        tinyint     NOT NULL DEFAULT 1 COMMENT '状态：0 停用 / 1 启用',
  creator       varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  created_at    datetime    NULL COMMENT '创建时间',
  updater       varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at    datetime    NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_rule (print_shop_id, paper_type, color_type, print_mode),
  KEY idx_shop (print_shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打印价格配置表'`)
}

// Down 删除价格配置表，无参数，返回删除结果。
func (r *M20261027000015CreatePrintPriceConfigTable) Down() error {
	return facades.Schema().DropIfExists("print_price_config")
}
