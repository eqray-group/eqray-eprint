package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000023CreatePayPaymentTable 建 pay_payment 支付单表。
// 打印订单支付与钱包充值共用一张表，用 biz_type + biz_id 区分业务；
// payment_no 即微信 out_trade_no，回调与对账都按它定位。
// channel 只用于统计口径，channel_code 才决定用哪套商户配置。
type M20261027000023CreatePayPaymentTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000023CreatePayPaymentTable) Signature() string {
	return "20261027000023_create_pay_payment_table"
}

// Up 建支付单表；同一业务单同时只允许一个进行中的支付单，由服务层保证。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000023CreatePayPaymentTable) Up() error {
	if facades.Schema().HasTable("pay_payment") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE pay_payment (
  id               bigint       NOT NULL COMMENT '主键',
  payment_no       varchar(32)  NOT NULL COMMENT '商户支付单号（雪花），全局唯一，即微信 out_trade_no',
  biz_type         varchar(32)  NOT NULL COMMENT '业务类型：print_order 打印订单 / wallet_recharge 钱包充值',
  biz_id           bigint       NOT NULL COMMENT '业务单 ID：打印订单为 print_order.id，充值为 pay_wallet.id',
  member_user_id   bigint       NOT NULL COMMENT '付款会员 ID，对应 member_user.id',
  channel          varchar(16)  NOT NULL COMMENT '渠道大类：wechat 微信 / alipay 支付宝 / balance 钱包余额',
  channel_code     varchar(32)  NOT NULL DEFAULT '' COMMENT '渠道配置编码，对应 pay_channel_config.code；balance 为空',
  scene            varchar(16)  NOT NULL DEFAULT '' COMMENT '支付场景：mini 小程序 / jsapi 公众号 / h5 手机浏览器 / wap 支付宝；balance 为空',
  client_ip        varchar(50)  NOT NULL DEFAULT '' COMMENT '下单客户端 IP，微信 H5 支付必填',
  amount           bigint       NOT NULL COMMENT '应付金额（分）',
  paid_amount      bigint       NOT NULL DEFAULT 0 COMMENT '实付金额（分）',
  refund_amount    bigint       NOT NULL DEFAULT 0 COMMENT '已退款金额（分），不得超过 paid_amount',
  status           tinyint      NOT NULL DEFAULT 0 COMMENT '状态：0 待支付 / 1 支付中 / 2 已支付 / 3 支付失败 / 4 已关闭 / 5 部分退款 / 6 全额退款',
  channel_trade_no varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道交易号（微信 transaction_id）',
  open_id          varchar(64)  NOT NULL DEFAULT '' COMMENT '微信 openid，JSAPI（小程序 / 公众号）支付用',
  expire_time      datetime     NULL COMMENT '支付超时时间，默认创建后 30 分钟',
  success_time     datetime     NULL COMMENT '支付成功时间',
  close_time       datetime     NULL COMMENT '关闭时间（超时关单）',
  remark           varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  creator          varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at       datetime     NULL COMMENT '创建时间',
  updater          varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at       datetime     NULL COMMENT '更新时间',
  deleted_at       datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_payment_no (payment_no),
  KEY idx_biz (biz_type, biz_id),
  KEY idx_channel_code (channel_code),
  KEY idx_channel_trade_no (channel_trade_no),
  KEY idx_member (member_user_id, created_at),
  KEY idx_status_expire (status, expire_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员支付单表'`)
}

// Down 删除支付单表，无参数，返回删除结果。
func (r *M20261027000023CreatePayPaymentTable) Down() error {
	return facades.Schema().DropIfExists("pay_payment")
}
