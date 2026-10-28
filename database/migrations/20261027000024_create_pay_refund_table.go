package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000024CreatePayRefundTable 建 pay_refund 退款单表。
// 一次支付可以有多笔退款，累计退款金额不得超过实付金额；
// 退款必须使用下单时那套商户配置，因此冗余 channel_code。
// 微信渠道走异步回调，钱包渠道直接入账并立即置成功。
type M20261027000024CreatePayRefundTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000024CreatePayRefundTable) Signature() string {
	return "20261027000024_create_pay_refund_table"
}

// Up 建退款单表；退款单号全局唯一，按支付单、业务单与会员建索引。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000024CreatePayRefundTable) Up() error {
	if facades.Schema().HasTable("pay_refund") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE pay_refund (
  id                bigint       NOT NULL COMMENT '主键',
  refund_no         varchar(32)  NOT NULL COMMENT '退款单号（雪花），全局唯一',
  payment_id        bigint       NOT NULL COMMENT '支付单 ID，对应 pay_payment.id',
  payment_no        varchar(32)  NOT NULL COMMENT '支付单号（冗余，便于对账）',
  biz_type          varchar(32)  NOT NULL COMMENT '业务类型，与支付单一致',
  biz_id            bigint       NOT NULL COMMENT '业务单 ID，与支付单一致',
  member_user_id    bigint       NOT NULL COMMENT '收款会员 ID，对应 member_user.id',
  channel           varchar(16)  NOT NULL COMMENT '渠道大类：与支付单一致，wechat / alipay / balance',
  channel_code      varchar(32)  NOT NULL DEFAULT '' COMMENT '渠道配置编码，与支付单一致；balance 为空',
  amount            bigint       NOT NULL COMMENT '退款金额（分）',
  reason            varchar(255) NOT NULL DEFAULT '' COMMENT '退款原因',
  status            tinyint      NOT NULL DEFAULT 0 COMMENT '状态：0 待退款 / 1 退款中 / 2 退款成功 / 3 退款失败 / 4 已取消',
  channel_refund_no varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道退款单号（微信 refund_id）',
  success_time      datetime     NULL COMMENT '退款成功时间',
  fail_reason       varchar(255) NOT NULL DEFAULT '' COMMENT '失败原因',
  operator_id       bigint       NOT NULL DEFAULT 0 COMMENT '发起退款的后台操作人，对应 system_users.id；0 表示系统自动',
  creator           varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at        datetime     NULL COMMENT '创建时间',
  updater           varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at        datetime     NULL COMMENT '更新时间',
  deleted_at        datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_refund_no (refund_no),
  KEY idx_payment (payment_id),
  KEY idx_biz (biz_type, biz_id),
  KEY idx_member (member_user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员退款单表'`)
}

// Down 删除退款单表，无参数，返回删除结果。
func (r *M20261027000024CreatePayRefundTable) Down() error {
	return facades.Schema().DropIfExists("pay_refund")
}
