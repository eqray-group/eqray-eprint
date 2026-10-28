package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000012CreatePrintShopStaffTable 建 print_shop_staff 门店成员绑定表。
// 后台账号与门店是多对多：一个账号可管多家店，角色按 (user_id, shop_id) 存，
// 因此同一员工在 A 店是店主、B 店是店员可以表达。
// 解绑走物理删除，不建 deleted_at，避免唯一键挡住重新绑定。
type M20261027000012CreatePrintShopStaffTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000012CreatePrintShopStaffTable) Signature() string {
	return "20261027000012_create_print_shop_staff_table"
}

// Up 建门店成员表，uk_user_shop 保证同一账号在同一门店只有一条绑定。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000012CreatePrintShopStaffTable) Up() error {
	if facades.Schema().HasTable("print_shop_staff") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_shop_staff (
  id         bigint      NOT NULL COMMENT '主键',
  user_id    bigint      NOT NULL COMMENT '后台账号 ID，对应 system_users.id',
  shop_id    bigint      NOT NULL COMMENT '门店 ID，对应 print_shop.id',
  role       varchar(16) NOT NULL DEFAULT 'staff' COMMENT '门店角色：owner 店主 / staff 店员',
  status     tinyint     NOT NULL DEFAULT 1 COMMENT '状态：0 停用 / 1 启用',
  creator    varchar(64) NOT NULL DEFAULT '' COMMENT '创建人（操作人）',
  created_at datetime    NULL COMMENT '创建时间',
  updater    varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at datetime    NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_shop (user_id, shop_id),
  KEY idx_shop (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='门店成员绑定表'`)
}

// Down 删除门店成员表，无参数，返回删除结果。
func (r *M20261027000012CreatePrintShopStaffTable) Down() error {
	return facades.Schema().DropIfExists("print_shop_staff")
}
