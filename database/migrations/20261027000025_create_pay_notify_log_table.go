package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000025CreatePayNotifyLogTable 建 pay_notify_log 支付渠道回调日志表。
// 微信 / 支付宝回调无法重放，必须落原始报文（脱敏）便于排查。
// 日志表只留 created_at，不需要审计四字段。
type M20261027000025CreatePayNotifyLogTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000025CreatePayNotifyLogTable) Signature() string {
	return "20261027000025_create_pay_notify_log_table"
}

// Up 建回调日志表；按商户单号与创建时间建索引，便于定位单笔回调。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000025CreatePayNotifyLogTable) Up() error {
	if facades.Schema().HasTable("pay_notify_log") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE pay_notify_log (
  id            bigint       NOT NULL COMMENT '主键',
  channel       varchar(16)  NOT NULL COMMENT '渠道：wechat / alipay',
  notify_type   varchar(16)  NOT NULL COMMENT '回调类型：payment 支付 / refund 退款',
  out_trade_no  varchar(64)  NOT NULL DEFAULT '' COMMENT '商户单号（支付回调用 out_trade_no，退款回调用 out_refund_no）',
  channel_no    varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道单号（transaction_id / refund_id）',
  body          text         NULL COMMENT '原始报文（脱敏后存储）',
  verified      tinyint      NOT NULL DEFAULT 0 COMMENT '验签结果：0 失败 / 1 通过',
  handle_status tinyint      NOT NULL DEFAULT 0 COMMENT '处理结果：0 待处理 / 1 处理成功 / 2 幂等忽略 / 3 处理失败',
  handle_msg    varchar(255) NOT NULL DEFAULT '' COMMENT '处理说明',
  created_at    datetime     NULL COMMENT '创建时间',
  PRIMARY KEY (id),
  KEY idx_out_trade_no (out_trade_no),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付渠道回调日志表'`)
}

// Down 删除回调日志表，无参数，返回删除结果。
func (r *M20261027000025CreatePayNotifyLogTable) Down() error {
	return facades.Schema().DropIfExists("pay_notify_log")
}
