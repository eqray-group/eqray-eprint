package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000022CreatePayWalletTransactionTable 建 pay_wallet_transaction 钱包流水表。
// 每一次余额变动必须有一条对应流水，与余额更新同事务。
// uk_biz 是幂等键：同一业务单重复入账或重复扣款会直接撞键，业务层据此判定「已处理」。
type M20261027000022CreatePayWalletTransactionTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000022CreatePayWalletTransactionTable) Signature() string {
	return "20261027000022_create_pay_wallet_transaction_table"
}

// Up 建流水表；按钱包、会员与时间建索引，供账单页分页查询。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000022CreatePayWalletTransactionTable) Up() error {
	if facades.Schema().HasTable("pay_wallet_transaction") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE pay_wallet_transaction (
  id             bigint       NOT NULL COMMENT '主键',
  wallet_id      bigint       NOT NULL COMMENT '钱包 ID，对应 pay_wallet.id',
  member_user_id bigint       NOT NULL COMMENT '会员 ID（冗余，便于按人查流水）',
  biz_type       varchar(32)  NOT NULL COMMENT '业务类型：recharge 充值 / print_consume 打印消费 / print_refund 打印退款 / manual_adjust 人工调账',
  biz_id         bigint       NOT NULL COMMENT '业务单 ID：充值为 pay_payment.id，消费与退款为 print_order.id',
  type           tinyint      NOT NULL COMMENT '资金方向：1 收入 / 2 支出',
  amount         bigint       NOT NULL COMMENT '变动金额（分，正数，方向由 type 表示）',
  balance_after  bigint       NOT NULL COMMENT '变动后可用余额（分）',
  title          varchar(64)  NOT NULL DEFAULT '' COMMENT '展示文案，如「打印订单 P2026... 消费」',
  remark         varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  creator        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at     datetime     NULL COMMENT '创建时间',
  updater        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at     datetime     NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_biz (biz_type, biz_id, type),
  KEY idx_wallet_created (wallet_id, created_at),
  KEY idx_member (member_user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钱包流水表'`)
}

// Down 删除钱包流水表，无参数，返回删除结果。
func (r *M20261027000022CreatePayWalletTransactionTable) Down() error {
	return facades.Schema().DropIfExists("pay_wallet_transaction")
}
