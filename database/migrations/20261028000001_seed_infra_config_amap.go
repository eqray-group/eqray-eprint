package migrations

import (
	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	infraModel "eqray-eprint/app/modules/infra/model"
)

type M20261028000001SeedInfraConfigAmap struct{}

// Signature The unique signature for the migration.
func (r *M20261028000001SeedInfraConfigAmap) Signature() string {
	return "20261028000001_seed_infra_config_amap"
}

// Up 写入高德地图相关参数（门店经纬度选点用），value 留空由管理员在
// 「基础设施 → 参数配置」里填写。visible=1 表示允许登录用户通过
// /api/v1/infra/configs/public 读取，前端地图组件才能拿到 Key。
func (r *M20261028000001SeedInfraConfigAmap) Up() error {
	if !facades.Schema().HasTable(infraModel.InfraConfigTable) {
		return nil
	}

	now := types.Now()
	items := []struct {
		name   string
		key    string
		remark string
	}{
		{
			name:   "高德地图 Web 端 Key",
			key:    "print.amap.web.key",
			remark: "门店地图选点使用；在高德开放平台创建「Web 端（JS API）」应用后获取",
		},
		{
			name:   "高德地图安全密钥",
			key:    "print.amap.web.security_code",
			remark: "高德 JS API 2.0 启用安全密钥后填写，未启用留空",
		},
	}

	for _, item := range items {
		var exists infraModel.InfraConfigModel
		if err := facades.DB().Table(infraModel.InfraConfigTable).
			Where("config_key", item.key).
			WhereNull("deleted_at").
			First(&exists); err != nil {
			return err
		}
		if exists.ID != 0 {
			continue
		}

		if _, err := facades.DB().Table(infraModel.InfraConfigTable).Insert(map[string]any{
			"id":         support.NewID(),
			"category":   "地图配置",
			"name":       item.name,
			"config_key": item.key,
			"value":      "",
			"type":       1,
			"visible":    1,
			"remark":     item.remark,
			"creator":    "system",
			"created_at": now,
			"updater":    "system",
			"updated_at": now,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down 删除本迁移写入的两条参数（管理员改过值也一并删除）。
func (r *M20261028000001SeedInfraConfigAmap) Down() error {
	if !facades.Schema().HasTable(infraModel.InfraConfigTable) {
		return nil
	}

	_, err := facades.DB().Table(infraModel.InfraConfigTable).
		WhereIn("config_key", []any{"print.amap.web.key", "print.amap.web.security_code"}).
		Delete()

	return err
}
