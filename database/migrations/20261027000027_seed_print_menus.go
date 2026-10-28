package migrations

import (
	"strconv"
	"time"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/system/model"
)

// M20261027000027SeedPrintMenus 写入打印管理菜单与权限点。
// 五个页面：打印门店、打印设备、打印机、价格配置、打印订单；
// 前端页面未就绪时点菜单会白屏，需与 views/print/* 一起上线。
type M20261027000027SeedPrintMenus struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000027SeedPrintMenus) Signature() string {
	return "20261027000027_seed_print_menus"
}

// Up 依次写入目录、页面、按钮权限点，并给内置角色补齐授权。
// 无参数；先查后插，可重复执行；菜单表不存在时直接返回 nil。
func (r *M20261027000027SeedPrintMenus) Up() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	now := time.Now()

	printMenuID, err := r.ensureMenu(0, "打印管理", 1, 50, "/print", "", "ri:printer-line", "", now)
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
			name: "打印门店", path: "shop", component: "/print/shop", icon: "ri:store-2-line", sort: 1,
			buttons: [][2]string{
				{"门店查询", "print:shop:list"},
				{"门店新增", "print:shop:create"},
				{"门店编辑", "print:shop:update"},
				{"门店删除", "print:shop:delete"},
				{"门店审核", "print:shop:audit"},
				{"成员查询", "print:staff:list"},
				{"成员绑定", "print:staff:create"},
				{"成员解绑", "print:staff:delete"},
			},
		},
		{
			name: "打印设备", path: "device", component: "/print/device", icon: "ri:computer-line", sort: 2,
			buttons: [][2]string{
				{"设备查询", "print:device:list"},
				{"设备新增", "print:device:create"},
				{"设备编辑", "print:device:update"},
				{"设备删除", "print:device:delete"},
			},
		},
		{
			name: "打印机", path: "printer", component: "/print/printer", icon: "ri:printer-line", sort: 3,
			buttons: [][2]string{
				{"打印机查询", "print:printer:list"},
				{"打印机新增", "print:printer:create"},
				{"打印机编辑", "print:printer:update"},
				{"打印机停用", "print:printer:delete"},
			},
		},
		{
			name: "价格配置", path: "price", component: "/print/price", icon: "ri:price-tag-3-line", sort: 4,
			buttons: [][2]string{
				{"价格查询", "print:price:list"},
				{"价格新增", "print:price:create"},
				{"价格编辑", "print:price:update"},
				{"价格删除", "print:price:delete"},
			},
		},
		{
			name: "打印文件", path: "file", component: "/print/file", icon: "ri:file-text-line", sort: 5,
			buttons: [][2]string{
				{"文件查询", "print:file:list"},
			},
		},
		{
			name: "轮播图", path: "banner", component: "/print/banner", icon: "ri:image-line", sort: 6,
			buttons: [][2]string{
				{"轮播查询", "print:banner:list"},
				{"轮播新增", "print:banner:create"},
				{"轮播编辑", "print:banner:update"},
				{"轮播删除", "print:banner:delete"},
			},
		},
		{
			name: "打印订单", path: "order", component: "/print/order", icon: "ri:file-list-2-line", sort: 7,
			buttons: [][2]string{
				{"订单查询", "print:order:list"},
				{"取件核销", "print:order:pickup"},
				{"订单退款", "print:order:refund"},
			},
		},
	}

	for _, page := range pages {
		pageID, err := r.ensureMenu(printMenuID, page.name, 2, page.sort, page.path, page.component, page.icon, "", now)
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

	_ = facades.Cache().Put("rbac:version", strconv.FormatInt(time.Now().UnixNano(), 10), 30*time.Minute)

	return nil
}

// ensureMenu 按「父级 + 类型 + 名称」查找菜单，不存在时插入并返回其 ID。
func (r *M20261027000027SeedPrintMenus) ensureMenu(parentID types.ID, name string, menuType, sort int, path, component, icon, permission string, now time.Time) (types.ID, error) {
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
func (r *M20261027000027SeedPrintMenus) grantBuiltinRoles(now time.Time) error {
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

// Down 删除打印管理目录及其下所有页面、按钮与角色授权。
func (r *M20261027000027SeedPrintMenus) Down() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	menus := make([]model.SystemMenuModel, 0)
	if err := facades.DB().Table("system_menu").
		Select("id").
		Where("name", "打印管理").
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
