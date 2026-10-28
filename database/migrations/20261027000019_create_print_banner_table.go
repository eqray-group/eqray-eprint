package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000019CreatePrintBannerTable 建 print_banner 门店轮播图表。
// print_shop_id = 0 表示全平台通用，大于 0 表示仅该门店展示。
type M20261027000019CreatePrintBannerTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000019CreatePrintBannerTable) Signature() string {
	return "20261027000019_create_print_banner_table"
}

// Up 建轮播图表，按「门店 + 排序」建索引支撑前台取图。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000019CreatePrintBannerTable) Up() error {
	if facades.Schema().HasTable("print_banner") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_banner (
  id            bigint       NOT NULL COMMENT '主键',
  print_shop_id bigint       NOT NULL DEFAULT 0 COMMENT '门店 ID，0 表示全平台通用',
  title         varchar(128) NOT NULL COMMENT '标题',
  image_url     varchar(512) NOT NULL COMMENT '图片地址',
  link_type     varchar(16)  NOT NULL DEFAULT 'none' COMMENT '跳转类型：none 不跳 / url 外部链接 / page 站内页面',
  link_url      varchar(512) NOT NULL DEFAULT '' COMMENT '跳转地址',
  sort          int          NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  status        tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 停用 / 1 启用',
  creator       varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at    datetime     NULL COMMENT '创建时间',
  updater       varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at    datetime     NULL COMMENT '更新时间',
  deleted_at    datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  KEY idx_shop_sort (print_shop_id, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='门店轮播图表'`)
}

// Down 删除轮播图表，无参数，返回删除结果。
func (r *M20261027000019CreatePrintBannerTable) Down() error {
	return facades.Schema().DropIfExists("print_banner")
}
