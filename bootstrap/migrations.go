package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"eqray-eprint/database/migrations"
)

// Migrations 返回全部迁移，顺序即执行顺序。
// 约定：一个表对应一个 create_*_table 迁移；数据填充（菜单、初始账号、默认配置等）
// 单独放 seed_* 迁移，不混在建表迁移里。
func Migrations() []schema.Migration {
	return []schema.Migration{
		// ---------- system / base ----------
		&migrations.M20261017000001CreateSystemUsersTable{},
		&migrations.M20261017000002CreateSystemMenuTable{},
		&migrations.M20261017000003CreateSystemRoleTable{},
		&migrations.M20261017000004CreateSystemUserRoleTable{},
		&migrations.M20261017000005CreateSystemRoleMenuTable{},
		&migrations.M20261017000006SeedSystemMenu{},
		&migrations.M20261017000007SeedSystemUsers{},
		&migrations.M20261017000008SeedSystemRole{},
		&migrations.M20261017000009SeedMenuIconAndDashboard{},
		&migrations.M20261017000010SeedSystemUserRole{},
		&migrations.M20261019000001CreateAuthLoginLogTable{},
		&migrations.M20261019000002CreateSystemDictTypeTable{},
		&migrations.M20261019000003CreateSystemDictDataTable{},
		&migrations.M20261019000004SeedDictData{},
		&migrations.M20261019000005SeedDictAndLoginLogMenu{},

		// ---------- infra ----------
		&migrations.M20261020000001CreateInfraFileTable{},
		&migrations.M20261020000002CreateSystemOperateLogTable{},
		&migrations.M20261020000003SeedOperateLogMenu{},
		&migrations.M20261020000004CreateInfraConfigTable{},
		&migrations.M20261020000005SeedInfraConfigMenu{},
		&migrations.M20261020000006SeedInfraFileMenu{},
		&migrations.M20261020000007CreateInfraFileConfigTable{},
		&migrations.M20261020000008CreateInfraApiAccessLogTable{},
		&migrations.M20261020000009CreateInfraApiErrorLogTable{},
		&migrations.M20261020000012SeedInfraOpsMenus{},
		&migrations.M20261020000013SeedInfraFileConfig{},
		&migrations.M20261020000030CreateInfraSmsChannelTable{},
		&migrations.M20261020000031CreateInfraSmsTemplateTable{},
		&migrations.M20261020000032SeedInfraSms{},
		&migrations.M20261020000033SeedInfraSmsMenu{},
		// 高德地图 Key：门店经纬度选点使用，值在参数配置里填写
		&migrations.M20261028000001SeedInfraConfigAmap{},

		// ---------- member ----------
		// 会员模块重建：只保留会员档案、会员卡、收货地址与第三方登录；
		// 等级 / 签到 / 积分 / 分组 / 标签 / 经验值体系已全部下线
		&migrations.M20261027000001CreateMemberUserTable{},
		&migrations.M20261027000004CreateMemberCardTable{},
		&migrations.M20261027000005CreateMemberAddressTable{},
		&migrations.M20261027000006CreateMemberSocialClientTable{},
		&migrations.M20261027000007CreateMemberSocialUserTable{},
		&migrations.M20261027000008CreateMemberSocialBindingTable{},
		&migrations.M20261027000009SeedMemberMenus{},

		// ---------- print ----------
		// 到店打印：门店 → 成员 → 设备 → 打印机 → 价格 → 文件 → 订单 → 明细 → 轮播
		// 表结构见 docs/p.sql；菜单种子等前端页面就绪后再补
		&migrations.M20261027000011CreatePrintShopTable{},
		&migrations.M20261027000012CreatePrintShopStaffTable{},
		&migrations.M20261027000013CreatePrintDeviceTable{},
		&migrations.M20261027000014CreatePrintPrinterTable{},
		&migrations.M20261027000015CreatePrintPriceConfigTable{},
		&migrations.M20261027000016CreatePrintFileTable{},
		&migrations.M20261027000017CreatePrintOrderTable{},
		&migrations.M20261027000018CreatePrintOrderItemTable{},
		&migrations.M20261027000019CreatePrintBannerTable{},

		// ---------- pay ----------
		// 支付渠道配置 → 会员钱包 → 流水 → 支付单 → 退款单 → 回调日志
		&migrations.M20261027000020CreatePayChannelConfigTable{},
		&migrations.M20261027000021CreatePayWalletTable{},
		&migrations.M20261027000022CreatePayWalletTransactionTable{},
		&migrations.M20261027000023CreatePayPaymentTable{},
		&migrations.M20261027000024CreatePayRefundTable{},
		&migrations.M20261027000025CreatePayNotifyLogTable{},
		&migrations.M20261027000026SeedPayMenus{},
		&migrations.M20261027000027SeedPrintMenus{},
		&migrations.M20261027000028SeedMissingPerms{},
	}
}
