package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000001CreateMemberUserTable 建 member_user 会员档案表。
// 会员只保留身份与资料：经验值、积分、标签、分组、等级体系全部下线，
// 支付折扣由 member_card 独立承担，本表不再持有任何等级字段。
type M20261027000001CreateMemberUserTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000001CreateMemberUserTable) Signature() string {
	return "20261027000001_create_member_user_table"
}

// Up 建 member_user 表；手机号建唯一键，会员身份以手机号为准，软删后号码不复用。
// 无参数；表已存在时直接返回 nil，保证可重复执行。
func (r *M20261027000001CreateMemberUserTable) Up() error {
	if facades.Schema().HasTable("member_user") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE member_user (
  id          bigint       NOT NULL COMMENT '主键',
  nickname    varchar(64)  NOT NULL COMMENT '昵称',
  mobile      varchar(20)  NOT NULL COMMENT '手机号，唯一，会员身份标识',
  password    varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
  avatar      varchar(512) NOT NULL DEFAULT '' COMMENT '头像地址',
  name        varchar(30)  NOT NULL DEFAULT '' COMMENT '真实姓名',
  sex         tinyint      NOT NULL DEFAULT 0 COMMENT '性别：0 未知 / 1 男 / 2 女',
  birthday    datetime     NULL COMMENT '出生日期',
  area_id     int          NOT NULL DEFAULT 0 COMMENT '所在地',
  mark        varchar(255) NOT NULL DEFAULT '' COMMENT '后台备注标记',
  register_ip varchar(32)  NOT NULL DEFAULT '' COMMENT '注册 IP',
  login_ip    varchar(32)  NOT NULL DEFAULT '' COMMENT '最后登录 IP',
  login_date  datetime     NULL COMMENT '最后登录时间',
  status      tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 禁用 / 1 启用',
  creator     varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at  datetime     NULL COMMENT '创建时间',
  updater     varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at  datetime     NULL COMMENT '更新时间',
  deleted_at  datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_mobile (mobile),
  KEY idx_status (status),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员档案表'`)
}

// Down 删除会员档案表，无参数，返回删除结果。
func (r *M20261027000001CreateMemberUserTable) Down() error {
	return facades.Schema().DropIfExists("member_user")
}
