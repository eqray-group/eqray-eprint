package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/model"
)

const deviceTable = model.PrintDeviceTable

// DeviceFilter 是设备列表的筛选条件，零值表示不限制。
type DeviceFilter struct {
	ShopID       types.ID
	Keyword      string
	Status       *int8
	OnlineStatus *int8
}

// PrintDeviceRepository 是 print_device 表的仓储。
type PrintDeviceRepository struct{}

// NewPrintDeviceRepository 构造设备仓储。
func NewPrintDeviceRepository() *PrintDeviceRepository {
	return &PrintDeviceRepository{}
}

// Create 新增设备。
func (r *PrintDeviceRepository) Create(ctx context.Context, m *model.PrintDeviceModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(deviceTable).Insert(map[string]any{
		"id":                 m.ID,
		"print_shop_id":      m.PrintShopID,
		"code":               m.Code,
		"name":               m.Name,
		"client_version":     m.ClientVersion,
		"status":             m.Status,
		"online_status":      m.OnlineStatus,
		"ip":                 m.IP,
		"last_heartbeat_at":  m.LastHeartbeatAt,
		"last_login_user_id": m.LastLoginUserID,
		"last_error_info":    m.LastErrorInfo,
		"remark":             m.Remark,
		"creator":            m.Creator,
		"created_at":         m.CreatedAt,
		"updater":            m.Updater,
		"updated_at":         m.UpdatedAt,
	})

	return err
}

// Update 更新设备基础信息（名称、备注、状态、所属门店）。
func (r *PrintDeviceRepository) Update(ctx context.Context, m *model.PrintDeviceModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"print_shop_id": m.PrintShopID,
			"name":          m.Name,
			"status":        m.Status,
			"remark":        m.Remark,
			"updater":       m.Updater,
			"updated_at":    m.UpdatedAt,
		})

	return err
}

// UpdateHeartbeat 记录心跳：刷新在线状态、IP、版本与最近登录账号。
// 客户端每次拉任务或空转心跳都会调用，字段少、开销小。
func (r *PrintDeviceRepository) UpdateHeartbeat(ctx context.Context, id types.ID, ip, version string, loginUserID types.ID) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"online_status":      1,
			"ip":                 ip,
			"client_version":     version,
			"last_heartbeat_at":  now,
			"last_login_user_id": loginUserID,
			"updated_at":         now,
		})

	return err
}

// MarkOffline 把超过阈值的设备置为离线，供定时任务批量清理。
func (r *PrintDeviceRepository) MarkOffline(ctx context.Context, before types.DateTime) (int64, error) {
	res, err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("online_status", 1).
		Where("last_heartbeat_at", "<", before).
		WhereNull("deleted_at").
		Update(map[string]any{
			"online_status": 0,
			"updated_at":    types.Now(),
		})
	if err != nil {
		return 0, err
	}

	return res.RowsAffected, nil
}

// Delete 软删除设备；其下打印机需先停用，避免任务被派到已删设备。
func (r *PrintDeviceRepository) Delete(ctx context.Context, id types.ID, updater string) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"deleted_at": now,
			"updater":    updater,
			"updated_at": now,
		})

	return err
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintDeviceRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintDeviceModel, error) {
	var m model.PrintDeviceModel

	if err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("id", id).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindByCode 按机器指纹查询；客户端登录登记设备时用。
func (r *PrintDeviceRepository) FindByCode(ctx context.Context, code string) (*model.PrintDeviceModel, error) {
	var m model.PrintDeviceModel

	if err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("code", code).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// ListByShop 查询门店下启用的设备。
func (r *PrintDeviceRepository) ListByShop(ctx context.Context, shopID types.ID) ([]model.PrintDeviceModel, error) {
	list := make([]model.PrintDeviceModel, 0)

	if err := facades.DB().WithContext(ctx).Table(deviceTable).
		Where("print_shop_id", shopID).
		Where("status", 1).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// Paginate 按条件分页查询设备。
func (r *PrintDeviceRepository) Paginate(ctx context.Context, filter DeviceFilter, pageNo, pageSize int) ([]model.PrintDeviceModel, int64, error) {
	list := make([]model.PrintDeviceModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(deviceTable).WhereNull("deleted_at")
	if filter.ShopID > 0 {
		query = query.Where("print_shop_id", filter.ShopID)
	}
	if filter.Keyword != "" {
		query = query.WhereLike("name", "%"+filter.Keyword+"%")
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}
	if filter.OnlineStatus != nil {
		query = query.Where("online_status", *filter.OnlineStatus)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
