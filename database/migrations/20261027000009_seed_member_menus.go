package migrations

import (
	"strconv"
	"time"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/system/model"
)

// M20261027000009SeedMemberMenus 写入会员管理菜单与权限点。
// 会员模块重建后只保留三个页面：会员列表、会员卡、社交客户端；
// 等级、签到、积分、分组、标签相关菜单与权限点不再生成。
type M20261027000009SeedMemberMenus struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000009SeedMemberMenus) Signature() string {
	return "20261027000009_seed_member_menus"
}

// Up 依次写入目录、页面、按钮权限点，并给内置角色补齐授权。
// 无参数；先查后插，可重复执行；菜单表不存在时直接返回 nil。
func (r *M20261027000009SeedMemberMenus) Up() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	now := time.Now()

	// 目录：会员管理
	memberID, err := r.ensureMenu(0, "会员管理", 1, 30, "/member", "", "ri:user-heart-line", "", now)
	if err != nil {
		return err
	}

	// 页面与其按钮权限点
	pages := []struct {
		name      string
		path      string
		component string
		icon      string
		sort      int
		buttons   [][2]string
	}{
		{
			name: "会员列表", path: "user", component: "/member/user", icon: "ri:user-3-line", sort: 1,
			buttons: [][2]string{
				{"会员查询", "member:user:list"},
				{"会员新增", "member:user:create"},
				{"会员修改", "member:user:update"},
				{"会员删除", "member:user:delete"},
			},
		},
		{
			name: "会员卡", path: "card", component: "/member/card", icon: "ri:bank-card-line", sort: 2,
			buttons: [][2]string{
				{"会员卡查询", "member:card:list"},
				{"会员卡发放与修改", "member:card:update"},
				{"会员卡停用", "member:card:delete"},
			},
		},
		{
			name: "社交客户端", path: "social-client", component: "/member/social-client", icon: "ri:wechat-line", sort: 3,
			buttons: [][2]string{
				{"客户端查询", "member:social-client:list"},
				{"客户端新增", "member:social-client:create"},
				{"客户端修改", "member:social-client:update"},
				{"客户端删除", "member:social-client:delete"},
			},
		},
	}

	for _, page := range pages {
		pageID, err := r.ensureMenu(memberID, page.name, 2, page.sort, page.path, page.component, page.icon, "", now)
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
func (r *M20261027000009SeedMemberMenus) ensureMenu(parentID types.ID, name string, menuType, sort int, path, component, icon, permission string, now time.Time) (types.ID, error) {
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
func (r *M20261027000009SeedMemberMenus) grantBuiltinRoles(now time.Time) error {
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

// Down 删除会员管理目录及其下所有页面、按钮与角色授权，返回删除结果。
func (r *M20261027000009SeedMemberMenus) Down() error {
	if !facades.Schema().HasTable("system_menu") {
		return nil
	}

	menus := make([]model.SystemMenuModel, 0)
	if err := facades.DB().Table("system_menu").
		Select("id").
		Where("name", "会员管理").
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
