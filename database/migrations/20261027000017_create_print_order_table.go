package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000017CreatePrintOrderTable 建 print_order 打印订单表（到店场景）。
// 流程：提交订单(待支付) → 支付成功(待打印) → 打印中 → 待取件 → 取件码核销完成。
// 只做到店打印，不涉及快递发货，因此没有收货信息与运费字段。
// pay_status 只表达钱，order_status 表达履约，两者组合由 service 约束。
type M20261027000017CreatePrintOrderTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000017CreatePrintOrderTable) Signature() string {
	return "20261027000017_create_print_order_table"
}

// Up 建订单表；订单号全局唯一，取件码按「门店 + 码」建索引便于核销定位。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000017CreatePrintOrderTable) Up() error {
	if facades.Schema().HasTable("print_order") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_order (
  id                     bigint       NOT NULL COMMENT '主键',
  order_no               varchar(32)  NOT NULL COMMENT '订单号，全局唯一',
  print_shop_id          bigint       NOT NULL COMMENT '门店 ID',
  member_user_id         bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id',
  device_id              bigint       NOT NULL DEFAULT 0 COMMENT '执行设备 ID，0 表示尚未路由（支付成功后由服务端分配）',
  pickup_code            char(6)      NOT NULL DEFAULT '' COMMENT '取件码，6 位数字，门店内未完成订单唯一',
  pickup_expire_time     datetime     NULL COMMENT '取件保留到期时间，默认支付后 7 天，到期转 8 超期关闭',
  pickup_verify_at       datetime     NULL COMMENT '取件核销时间',
  pickup_verify_user_id  bigint       NOT NULL DEFAULT 0 COMMENT '核销操作人，对应 system_users.id',
  item_count             int          NOT NULL DEFAULT 0 COMMENT '明细数量（文件数）',
  total_pages            int          NOT NULL DEFAULT 0 COMMENT '总页数（不含份数）',
  total_sheets           int          NOT NULL DEFAULT 0 COMMENT '总印张数（页数 × 份数）',
  total_amount           bigint       NOT NULL DEFAULT 0 COMMENT '明细金额合计（分）',
  discount_percent       int          NOT NULL DEFAULT 100 COMMENT '下单时的会员卡折扣快照：100 无折扣 / 95 表示 95 折',
  discount_amount        bigint       NOT NULL DEFAULT 0 COMMENT '优惠金额（分）= 明细合计 × (100 - discount_percent) / 100，向下取整',
  pay_amount             bigint       NOT NULL DEFAULT 0 COMMENT '应付金额（分）= 明细合计 - 优惠',
  refund_amount          bigint       NOT NULL DEFAULT 0 COMMENT '已退款金额（分）',
  pay_status             tinyint      NOT NULL DEFAULT 0 COMMENT '支付状态：0 未支付 / 1 已支付 / 2 退款中 / 3 已退款',
  pay_type               varchar(16)  NOT NULL DEFAULT '' COMMENT '支付方式：wechat 微信 / alipay 支付宝 / wallet 余额',
  trade_no               varchar(64)  NOT NULL DEFAULT '' COMMENT '第三方支付交易号',
  pay_time               datetime     NULL COMMENT '支付时间',
  order_status           tinyint      NOT NULL DEFAULT 0 COMMENT '订单状态：0 待支付 / 1 待打印 / 2 打印中 / 3 待取件 / 4 已完成 / 5 已取消 / 6 退款中 / 7 已退款 / 8 超期关闭',
  pay_expire_time        datetime     NULL COMMENT '未支付关单时间，默认创建后 30 分钟',
  finish_time            datetime     NULL COMMENT '完成时间',
  cancel_time            datetime     NULL COMMENT '取消时间',
  cancel_reason          varchar(255) NOT NULL DEFAULT '' COMMENT '取消原因',
  close_time             datetime     NULL COMMENT '关闭时间（超时关单 / 超期未取）',
  remark                 varchar(255) NOT NULL DEFAULT '' COMMENT '用户备注',
  creator                varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at             datetime     NULL COMMENT '创建时间',
  updater                varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at             datetime     NULL COMMENT '更新时间',
  deleted_at             datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_order_no (order_no),
  KEY idx_member (member_user_id, created_at),
  KEY idx_shop_status (print_shop_id, order_status, created_at),
  KEY idx_shop_pickup (print_shop_id, pickup_code),
  KEY idx_trade_no (trade_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打印订单表'`)
}

// Down 删除订单表，无参数，返回删除结果。
func (r *M20261027000017CreatePrintOrderTable) Down() error {
	return facades.Schema().DropIfExists("print_order")
}
