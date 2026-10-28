package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/model"
)

const printerTable = model.PrintPrinterTable

// PrinterFilter 是打印机列表的筛选条件，零值表示不限制。
type PrinterFilter struct {
	ShopID   types.ID
	DeviceID types.ID
	Status   *int8
}

// PrintPrinterRepository 是 print_printer 表的仓储。
type PrintPrinterRepository struct{}

// NewPrintPrinterRepository 构造打印机仓储。
func NewPrintPrinterRepository() *PrintPrinterRepository {
	return &PrintPrinterRepository{}
}

// Create 新增打印机；同一设备下 system_name 重复会返回数据库错误。
func (r *PrintPrinterRepository) Create(ctx context.Context, m *model.PrintPrinterModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(printerTable).Insert(map[string]any{
		"id":             m.ID,
		"print_shop_id":  m.PrintShopID,
		"device_id":      m.DeviceID,
		"name":           m.Name,
		"system_name":    m.SystemName,
		"model":          m.Model,
		"paper_types":    m.PaperTypes,
		"color_support":  m.ColorSupport,
		"duplex_support": m.DuplexSupport,
		"status":         m.Status,
		"remark":         m.Remark,
		"creator":        m.Creator,
		"created_at":     m.CreatedAt,
		"updater":        m.Updater,
		"updated_at":     m.UpdatedAt,
	})

	return err
}

// Update 更新打印机配置与能力字段。
func (r *PrintPrinterRepository) Update(ctx context.Context, m *model.PrintPrinterModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(printerTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"name":           m.Name,
			"system_name":    m.SystemName,
			"model":          m.Model,
			"paper_types":    m.PaperTypes,
			"color_support":  m.ColorSupport,
			"duplex_support": m.DuplexSupport,
			"status":         m.Status,
			"remark":         m.Remark,
			"updater":        m.Updater,
			"updated_at":     m.UpdatedAt,
		})

	return err
}

// Disable 停用打印机；不物理删除，避免队列名被删除记录长期占用。
func (r *PrintPrinterRepository) Disable(ctx context.Context, id types.ID, updater string) error {
	_, err := facades.DB().WithContext(ctx).Table(printerTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"status":     0,
			"updater":    updater,
			"updated_at": types.Now(),
		})

	return err
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintPrinterRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintPrinterModel, error) {
	var m model.PrintPrinterModel

	if err := facades.DB().WithContext(ctx).Table(printerTable).
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

// ListEnabledByShop 查询门店下启用的打印机，路由能力匹配用。
func (r *PrintPrinterRepository) ListEnabledByShop(ctx context.Context, shopID types.ID) ([]model.PrintPrinterModel, error) {
	list := make([]model.PrintPrinterModel, 0)

	if err := facades.DB().WithContext(ctx).Table(printerTable).
		Where("print_shop_id", shopID).
		Where("status", 1).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// Paginate 按条件分页查询打印机。
func (r *PrintPrinterRepository) Paginate(ctx context.Context, filter PrinterFilter, pageNo, pageSize int) ([]model.PrintPrinterModel, int64, error) {
	list := make([]model.PrintPrinterModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(printerTable).WhereNull("deleted_at")
	if filter.ShopID > 0 {
		query = query.Where("print_shop_id", filter.ShopID)
	}
	if filter.DeviceID > 0 {
		query = query.Where("device_id", filter.DeviceID)
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
