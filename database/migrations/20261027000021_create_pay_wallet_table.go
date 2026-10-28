package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000021CreatePayWalletTable 建 pay_wallet 会员钱包表。
// 只有会员有余额账户，后台账号不参与支付。
// 余额只允许与「钱包流水」同事务变更；扣减必须带 WHERE balance >= amount，
// 防止并发把余额扣成负数。
type M20261027000021CreatePayWalletTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000021CreatePayWalletTable) Signature() string {
	return "20261027000021_create_pay_wallet_table"
}

// Up 建钱包表，一个会员一个钱包（uk_member）。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000021CreatePayWalletTable) Up() error {
	if facades.Schema().HasTable("pay_wallet") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE pay_wallet (
  id             bigint      NOT NULL COMMENT '主键',
  member_user_id bigint      NOT NULL COMMENT '会员 ID，对应 member_user.id，一个会员一个钱包',
  balance        bigint      NOT NULL DEFAULT 0 COMMENT '可用余额（分），不允许为负',
  total_income   bigint      NOT NULL DEFAULT 0 COMMENT '累计收入（分）',
  total_expense  bigint      NOT NULL DEFAULT 0 COMMENT '累计支出（分）',
  status         tinyint     NOT NULL DEFAULT 1 COMMENT '状态：0 冻结（禁止支付）/ 1 正常',
  version        int         NOT NULL DEFAULT 0 COMMENT '乐观锁版本号，每次余额变动 +1',
  creator        varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  created_at     datetime    NULL COMMENT '创建时间',
  updater        varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at     datetime    NULL COMMENT '更新时间',
  deleted_at     datetime    NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_member (member_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员钱包表'`)
}

// Down 删除钱包表，无参数，返回删除结果。
func (r *M20261027000021CreatePayWalletTable) Down() error {
	return facades.Schema().DropIfExists("pay_wallet")
}
