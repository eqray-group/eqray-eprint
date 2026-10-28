package migrations

import (
	"time"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/system/model"
)

// M20261027000028SeedMissingPerms 补齐「路由已有、菜单种子缺失」的权限点。
// 缺权限点时后台按钮会因 Permission 中间件 403 而不显示，
// 涉及：用户分配角色、角色分配菜单、短信渠道、短信模板四组。
type M20261027000028SeedMissingPerms struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000028SeedMissingPerms) Signature() string {
	return "20261027000028_seed_missing_perms"
}

// Up 按页面 component 定位父菜单，补按钮并授权内置角色。
// 先查后插，可重复执行；页面不存在时跳过，不阻断迁移。
func (r *M20261027000028SeedMissingPerms) Up() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	now := time.Now()

	pages := []struct {
		component string
		buttons   [][2]string
	}{
		{
			component: "/system/user",
			buttons:   [][2]string{{"分配角色", "system:user:assign-role"}},
		},
		{
			component: "/system/role",
			buttons:   [][2]string{{"分配菜单", "system:role:assign-menu"}},
		},
		{
			component: "/infra/sms-channel",
			buttons: [][2]string{
				{"渠道查询", "infra:sms-channel:list"},
				{"渠道新增", "infra:sms-channel:create"},
				{"渠道编辑", "infra:sms-channel:update"},
				{"渠道删除", "infra:sms-channel:delete"},
			},
		},
		{
			component: "/infra/sms-template",
			buttons: [][2]string{
				{"模板查询", "infra:sms-template:list"},
				{"模板新增", "infra:sms-template:create"},
				{"模板编辑", "infra:sms-template:update"},
				{"模板删除", "infra:sms-template:delete"},
			},
		},
	}

	for _, page := range pages {
		pageID, err := r.pageID(page.component)
		if err != nil {
			return err
		}
		if pageID == 0 {
			// 页面不存在说明该模块未启用，跳过更安全
			continue
		}

		for index, button := range page.buttons {
			if err := r.ensureButton(pageID, button[0], button[1], index+1, now); err != nil {
				return err
			}
		}
	}

	if err := r.grantBuiltinRoles(now); err != nil {
		return err
	}

	_ = facades.Cache().Put("rbac:version", time.Now().Format("20060102150405.000000000"), 30*time.Minute)

	return nil
}

// pageID 按 component 找页面菜单，找不到返回 0。
func (r *M20261027000028SeedMissingPerms) pageID(component string) (types.ID, error) {
	var existing model.SystemMenuModel

	if err := facades.DB().Table("system_menu").
		Select("id").
		Where("component", component).
		Where("type", 2).
		WhereNull("deleted_at").
		First(&existing); err != nil {
		return 0, err
	}

	return existing.ID, nil
}

// ensureButton 按「父级 + 类型 + 权限点」查找按钮，不存在时插入。
func (r *M20261027000028SeedMissingPerms) ensureButton(parentID types.ID, name, permission string, sort int, now time.Time) error {
	var existing model.SystemMenuModel

	if err := facades.DB().Table("system_menu").
		Select("id").
		Where("parent_id", parentID).
		Where("type", 3).
		Where("permission", permission).
		WhereNull("deleted_at").
		First(&existing); err != nil {
		return err
	}
	if existing.ID > 0 {
		return nil
	}

	_, err := facades.DB().Table("system_menu").Insert(map[string]any{
		"id":             support.NewID(),
		"name":           name,
		"permission":     permission,
		"type":           3,
		"sort":           sort,
		"parent_id":      parentID,
		"path":           "",
		"icon":           "",
		"component":      "",
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

	return err
}

// grantBuiltinRoles 给全部内置角色补齐缺失的菜单授权。
func (r *M20261027000028SeedMissingPerms) grantBuiltinRoles(now time.Time) error {
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
		rows := make([]model.SystemRoleMenuModel, 0)
		if err := facades.DB().Table("system_role_menu").
			Select("menu_id").
			Where("role_id", role.ID).
			WhereNull("deleted_at").
			Get(&rows); err != nil {
			return err
		}

		granted := make(map[types.ID]bool, len(rows))
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

// Down 删除本迁移补的按钮及其角色授权，页面菜单本身保留。
func (r *M20261027000028SeedMissingPerms) Down() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	permissions := []string{
		"system:user:assign-role",
		"system:role:assign-menu",
		"infra:sms-channel:list",
		"infra:sms-channel:create",
		"infra:sms-channel:update",
		"infra:sms-channel:delete",
		"infra:sms-template:list",
		"infra:sms-template:create",
		"infra:sms-template:update",
		"infra:sms-template:delete",
	}

	menus := make([]model.SystemMenuModel, 0)
	if err := facades.DB().Table("system_menu").
		Select("id").
		WhereIn("permission", toAny(permissions)).
		Where("type", 3).
		Get(&menus); err != nil {
		return err
	}

	for _, menu := range menus {
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

// toAny 把字符串切片转成查询需要的 []any。
func toAny(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}

	return result
}
