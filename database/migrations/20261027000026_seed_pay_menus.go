package migrations

import (
	"strconv"
	"time"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/system/model"
)

// M20261027000026SeedPayMenus 写入支付管理菜单与权限点。
// 三个页面：支付渠道、支付流水、退款管理；钱包由会员端自助使用，后台不建页面。
type M20261027000026SeedPayMenus struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000026SeedPayMenus) Signature() string {
	return "20261027000026_seed_pay_menus"
}

// Up 依次写入目录、页面、按钮权限点，并给内置角色补齐授权。
// 无参数；先查后插，可重复执行；菜单表不存在时直接返回 nil。
func (r *M20261027000026SeedPayMenus) Up() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	now := time.Now()

	// 目录：支付管理
	payMenuID, err := r.ensureMenu(0, "支付管理", 1, 40, "/pay", "", "ri:money-cny-circle-line", "", now)
	if err != nil {
		return err
	}

	pages := []struct {
		name      string
		path      string
		component string
		icon      string
		sort      int
		buttons   [][2]string
	}{
		{
			name: "支付渠道", path: "channel-config", component: "/pay/channel-config", icon: "ri:secure-payment-line", sort: 1,
			buttons: [][2]string{
				{"渠道查询", "pay:channel:list"},
				{"渠道新增", "pay:channel:create"},
				{"渠道修改", "pay:channel:update"},
				{"渠道删除", "pay:channel:delete"},
			},
		},
		{
			name: "支付流水", path: "payment", component: "/pay/payment", icon: "ri:file-list-3-line", sort: 2,
			buttons: [][2]string{
				{"支付查询", "pay:payment:list"},
				{"钱包调账", "pay:wallet:adjust"},
			},
		},
		{
			name: "会员钱包", path: "wallet", component: "/pay/wallet", icon: "ri:wallet-3-line", sort: 3,
			buttons: [][2]string{
				{"钱包查询", "pay:wallet:list"},
				{"钱包调账", "pay:wallet:adjust"},
			},
		},
		{
			name: "退款管理", path: "refund", component: "/pay/refund", icon: "ri:refund-2-line", sort: 4,
			buttons: [][2]string{
				{"退款查询", "pay:refund:list"},
				{"发起退款", "pay:refund:create"},
			},
		},
	}

	for _, page := range pages {
		pageID, err := r.ensureMenu(payMenuID, page.name, 2, page.sort, page.path, page.component, page.icon, "", now)
		if err != nil {
			return err
		}

		for index, button := range page.buttons {
			if _, err := r.ensureMenu(pageID, button[0], 3, index+1, "", "", "", button[1], now); err != nil {
				return err
			}
		}
	}

	if err := r.grantBuiltinRoles(now); err != nil {
		return err
	}

	// 新增权限点后前移权限缓存版本，避免管理员要等 30 分钟才能拿到新按钮权限
	_ = facades.Cache().Put("rbac:version", strconv.FormatInt(time.Now().UnixNano(), 10), 30*time.Minute)

	return nil
}

// ensureMenu 按「父级 + 类型 + 名称」查找菜单，不存在时插入并返回其 ID。
// 参数 parentID 为父菜单 ID（0 表示顶级），name 为菜单名，menuType 为 1 目录 / 2 页面 / 3 按钮，
// sort 为排序，path / component / icon 为前端路由与展示配置，permission 为权限点，now 为写入时间。
func (r *M20261027000026SeedPayMenus) ensureMenu(parentID types.ID, name string, menuType, sort int, path, component, icon, permission string, now time.Time) (types.ID, error) {
	var existing model.SystemMenuModel
	if err := facades.DB().Table("system_menu").
		Select("id").
		Where("parent_id", parentID).
		Where("type", menuType).
		Where("name", name).
		WhereNull("deleted_at").
		First(&existing); err != nil {
		return 0, err
	}
	if existing.ID > 0 {
		return existing.ID, nil
	}

	id := support.NewID()
	_, err := facades.DB().Table("system_menu").Insert(map[string]any{
		"id":             id,
		"name":           name,
		"permission":     permission,
		"type":           menuType,
		"sort":           sort,
		"parent_id":      parentID,
		"path":           path,
		"icon":           icon,
		"component":      component,
		"component_name": "",
		"status":         1,
		"visible":        1,
		"keep_alive":     1,
		"always_show":    0,
		"creator":        "system",
		"created_at":     now,
		"updater":        "system",
		"updated_at":     now,
	})

	return id, err
}

// grantBuiltinRoles 给全部内置角色补齐缺失的菜单授权，返回写入错误。
// 参数 now 为授权记录写入时间。
func (r *M20261027000026SeedPayMenus) grantBuiltinRoles(now time.Time) error {
	if !facades.Schema().HasTable("system_role") || !facades.Schema().HasTable("system_role_menu") {
		return nil
	}

	roles := make([]model.SystemRoleModel, 0)
	if err := facades.DB().Table("system_role").
		Select("id").
		Where("type", 1).
		WhereNull("deleted_at").
		Get(&roles); err != nil {
		return err
	}

	menus := make([]model.SystemMenuModel, 0)
	if err := facades.DB().Table("system_menu").
		Select("id").
		WhereNull("deleted_at").
		Get(&menus); err != nil {
		return err
	}

	for _, role := range roles {
		granted := make(map[types.ID]bool, len(menus))

		rows := make([]model.SystemRoleMenuModel, 0)
		if err := facades.DB().Table("system_role_menu").
			Select("menu_id").
			Where("role_id", role.ID).
			WhereNull("deleted_at").
			Get(&rows); err != nil {
			return err
		}
		for _, row := range rows {
			granted[row.MenuID] = true
		}

		for _, menu := range menus {
			if granted[menu.ID] {
				continue
			}

			if _, err := facades.DB().Table("system_role_menu").Insert(map[string]any{
				"id":         support.NewID(),
				"role_id":    role.ID,
				"menu_id":    menu.ID,
				"creator":    "system",
				"created_at": now,
				"updater":    "system",
				"updated_at": now,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// Down 删除支付管理目录及其下所有页面、按钮与角色授权，返回删除结果。
func (r *M20261027000026SeedPayMenus) Down() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	menus := make([]model.SystemMenuModel, 0)
	if err := facades.DB().Table("system_menu").
		Select("id").
		Where("name", "支付管理").
		Where("type", 1).
		WhereNull("deleted_at").
		Get(&menus); err != nil {
		return err
	}

	for _, menu := range menus {
		pages := make([]model.SystemMenuModel, 0)
		if err := facades.DB().Table("system_menu").Select("id").Where("parent_id", menu.ID).Get(&pages); err != nil {
			return err
		}

		for _, page := range pages {
			buttons := make([]model.SystemMenuModel, 0)
			if err := facades.DB().Table("system_menu").Select("id").Where("parent_id", page.ID).Get(&buttons); err != nil {
				return err
			}

			for _, button := range buttons {
				if facades.Schema().HasTable("system_role_menu") {
					if _, err := facades.DB().Table("system_role_menu").Where("menu_id", button.ID).Delete(); err != nil {
						return err
					}
				}
				if _, err := facades.DB().Table("system_menu").Where("id", button.ID).Delete(); err != nil {
					return err
				}
			}

			if facades.Schema().HasTable("system_role_menu") {
				if _, err := facades.DB().Table("system_role_menu").Where("menu_id", page.ID).Delete(); err != nil {
					return err
				}
			}
			if _, err := facades.DB().Table("system_menu").Where("id", page.ID).Delete(); err != nil {
				return err
			}
		}

		if facades.Schema().HasTable("system_role_menu") {
			if _, err := facades.DB().Table("system_role_menu").Where("menu_id", menu.ID).Delete(); err != nil {
				return err
			}
		}
		if _, err := facades.DB().Table("system_menu").Where("id", menu.ID).Delete(); err != nil {
			return err
		}
	}

	return nil
}
