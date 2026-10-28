package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000004CreateMemberCardTable 建 member_card 会员卡表。
// 会员卡是独立的折扣载体：卡自己带等级与折扣，不引用 member_level，
// 与会员等级体系互不影响；支付折扣只读本表，读取顺序为
// 「卡存在且 status=1 且未过期 → discount_percent」，否则按 100（无折扣）处理。
type M20261027000004CreateMemberCardTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000004CreateMemberCardTable) Signature() string {
	return "20261027000004_create_member_card_table"
}

// Up 建 member_card 表：一个会员一张卡（uk_member），卡号全局唯一（uk_card_no）。
// 卡记录不做软删，停用与过期都改 status，避免唯一键被删除记录占用。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000004CreateMemberCardTable) Up() error {
	if facades.Schema().HasTable("member_card") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE member_card (
  id               bigint       NOT NULL COMMENT '主键',
  member_user_id   bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id，一个会员一张卡',
  card_no          varchar(32)  NOT NULL COMMENT '卡号，服务端生成（C + 雪花后 10 位），全局唯一',
  card_level       varchar(32)  NOT NULL DEFAULT 'normal' COMMENT '卡等级编码：normal 普通 / silver 银卡 / gold 金卡 / black 黑金，可按业务自定义',
  card_name        varchar(64)  NOT NULL COMMENT '卡名称，展示用，如「金卡会员」',
  discount_percent int          NOT NULL DEFAULT 100 COMMENT '卡自带折扣百分比：100 无折扣 / 95 表示 95 折，支付折扣的唯一来源',
  card_type        varchar(16)  NOT NULL DEFAULT 'level' COMMENT '卡类型：level 折扣卡（预留 stored 储值卡）',
  source           varchar(16)  NOT NULL DEFAULT 'admin' COMMENT '开卡来源：register 注册赠送 / admin 后台发放 / purchase 购买（预留）',
  status           tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 停用 / 1 正常 / 2 已过期',
  start_time       datetime     NULL COMMENT '开卡时间',
  expire_time      datetime     NULL COMMENT '到期时间，NULL 表示永久有效；到期由定时任务置 2，折扣随之失效',
  remark           varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  creator          varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at       datetime     NULL COMMENT '创建时间',
  updater          varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at       datetime     NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_card_no (card_no),
  UNIQUE KEY uk_member (member_user_id),
  KEY idx_card_level (card_level),
  KEY idx_status_expire (status, expire_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员卡表'`)
}

// Down 删除会员卡表，无参数，返回删除结果。
func (r *M20261027000004CreateMemberCardTable) Down() error {
	return facades.Schema().DropIfExists("member_card")
}
