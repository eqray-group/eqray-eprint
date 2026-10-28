package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000011CreatePrintShopTable 建 print_shop 门店表。
// 门店编号已取消，门店以主键标识；营业状态与审核状态分开存，
// 避免「停业」和「待审核」混在一个字段里。
type M20261027000011CreatePrintShopTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000011CreatePrintShopTable) Signature() string {
	return "20261027000011_create_print_shop_table"
}

// Up 建门店表；按状态、审核状态与经纬度建索引，支撑列表筛选与附近门店查询。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000011CreatePrintShopTable) Up() error {
	if facades.Schema().HasTable("print_shop") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_shop (
  id             bigint         NOT NULL COMMENT '主键',
  name           varchar(64)    NOT NULL COMMENT '门店名称',
  contact_name   varchar(32)    NOT NULL DEFAULT '' COMMENT '门店联系人',
  contact_mobile varchar(20)    NOT NULL DEFAULT '' COMMENT '门店联系电话',
  province       varchar(32)    NOT NULL DEFAULT '' COMMENT '省',
  city           varchar(32)    NOT NULL DEFAULT '' COMMENT '市',
  district       varchar(32)    NOT NULL DEFAULT '' COMMENT '区/县',
  address        varchar(255)   NOT NULL DEFAULT '' COMMENT '详细地址',
  longitude      decimal(10, 7) NOT NULL DEFAULT 0.0000000 COMMENT '经度',
  latitude       decimal(10, 7) NOT NULL DEFAULT 0.0000000 COMMENT '纬度',
  avatar         varchar(512)   NOT NULL DEFAULT '' COMMENT '门头照地址',
  business_hours json           NULL COMMENT '营业时间 JSON，如 [{"start":"09:00","end":"21:00"}]，NULL 表示全天',
  status         tinyint        NOT NULL DEFAULT 1 COMMENT '营业状态：0 停业 / 1 营业',
  audit_status   tinyint        NOT NULL DEFAULT 0 COMMENT '审核状态：0 待审核 / 1 审核通过 / 2 审核拒绝',
  audit_remark   varchar(255)   NOT NULL DEFAULT '' COMMENT '审核意见',
  remark         varchar(255)   NOT NULL DEFAULT '' COMMENT '备注',
  creator        varchar(64)    NOT NULL DEFAULT '' COMMENT '创建人',
  created_at     datetime       NULL COMMENT '创建时间',
  updater        varchar(64)    NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at     datetime       NULL COMMENT '更新时间',
  deleted_at     datetime       NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  KEY idx_status (status),
  KEY idx_audit_status (audit_status),
  KEY idx_geo (latitude, longitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打印门店表'`)
}

// Down 删除门店表，无参数，返回删除结果。
func (r *M20261027000011CreatePrintShopTable) Down() error {
	return facades.Schema().DropIfExists("print_shop")
}
