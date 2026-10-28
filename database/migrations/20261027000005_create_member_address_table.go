package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000005CreateMemberAddressTable 建 member_address 会员收货地址表。
// 当前业务只做到店打印，地址表先保留以支撑后续快递发货；
// 写入由 C 端会员自己维护，后台只读展示。
type M20261027000005CreateMemberAddressTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000005CreateMemberAddressTable) Signature() string {
	return "20261027000005_create_member_address_table"
}

// Up 建 member_address 表，按会员维度查询（idx_user）。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000005CreateMemberAddressTable) Up() error {
	if facades.Schema().HasTable("member_address") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE member_address (
  id             bigint       NOT NULL COMMENT '主键',
  user_id        bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id',
  name           varchar(30)  NOT NULL DEFAULT '' COMMENT '收件人',
  mobile         varchar(20)  NOT NULL DEFAULT '' COMMENT '手机号',
  area_id        int          NOT NULL DEFAULT 0 COMMENT '地区编号',
  detail_address varchar(255) NOT NULL DEFAULT '' COMMENT '详细地址',
  default_status tinyint      NOT NULL DEFAULT 0 COMMENT '是否默认：0 否 / 1 是',
  creator        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at     datetime     NULL COMMENT '创建时间',
  updater        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at     datetime     NULL COMMENT '更新时间',
  deleted_at     datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  KEY idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员收货地址表'`)
}

// Down 删除会员收货地址表，无参数，返回删除结果。
func (r *M20261027000005CreateMemberAddressTable) Down() error {
	return facades.Schema().DropIfExists("member_address")
}
