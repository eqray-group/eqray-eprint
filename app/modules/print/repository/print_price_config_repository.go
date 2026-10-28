package repository

import (
	"context"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/model"
)

const priceTable = model.PrintPriceConfigTable

// PrintPriceConfigRepository 是 print_price_config 表的仓储。
// 表无软删，删除规则走物理删除（规则属于可重建数据）。
type PrintPriceConfigRepository struct{}

// NewPrintPriceConfigRepository 构造价格配置仓储。
func NewPrintPriceConfigRepository() *PrintPriceConfigRepository {
	return &PrintPriceConfigRepository{}
}

// Create 新增价格规则。
func (r *PrintPriceConfigRepository) Create(ctx context.Context, m *model.PrintPriceConfigModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(priceTable).Insert(map[string]any{
		"id":            m.ID,
		"print_shop_id": m.PrintShopID,
		"paper_type":    m.PaperType,
		"color_type":    m.ColorType,
		"print_mode":    m.PrintMode,
		"price":         m.Price,
		"status":        m.Status,
		"creator":       m.Creator,
		"created_at":    m.CreatedAt,
		"updater":       m.Updater,
		"updated_at":    m.UpdatedAt,
	})

	return err
}

// Update 更新价格与启用状态。
func (r *PrintPriceConfigRepository) Update(ctx context.Context, m *model.PrintPriceConfigModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(priceTable).
		Where("id", m.ID).
		Update(map[string]any{
			"price":      m.Price,
			"status":     m.Status,
			"updater":    m.Updater,
			"updated_at": m.UpdatedAt,
		})

	return err
}

// Delete 物理删除价格规则，未命中返回 NotFound。
func (r *PrintPriceConfigRepository) Delete(ctx context.Context, id types.ID) error {
	res, err := facades.DB().WithContext(ctx).Table(priceTable).Where("id", id).Delete()
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return baseErrs.WithMsg(baseErrs.NotFound, "价格规则不存在")
	}

	return nil
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintPriceConfigRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintPriceConfigModel, error) {
	var m model.PrintPriceConfigModel

	if err := facades.DB().WithContext(ctx).Table(priceTable).Where("id", id).First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindByRule 按「门店 + 纸张 + 色彩 + 单双面」取唯一规则，计价时调用。
func (r *PrintPriceConfigRepository) FindByRule(ctx context.Context, shopID types.ID, paperType, colorType, printMode string) (*model.PrintPriceConfigModel, error) {
	var m model.PrintPriceConfigModel

	if err := facades.DB().WithContext(ctx).Table(priceTable).
		Where("print_shop_id", shopID).
		Where("paper_type", paperType).
		Where("color_type", colorType).
		Where("print_mode", printMode).
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// ListByShop 查询门店全部价格规则（含停用），后台配置页与能力体检用。
func (r *PrintPriceConfigRepository) ListByShop(ctx context.Context, shopID types.ID) ([]model.PrintPriceConfigModel, error) {
	list := make([]model.PrintPriceConfigModel, 0)

	if err := facades.DB().WithContext(ctx).Table(priceTable).
		Where("print_shop_id", shopID).
		OrderBy("paper_type", "asc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// Paginate 按门店分页查询价格规则。
func (r *PrintPriceConfigRepository) Paginate(ctx context.Context, shopID types.ID, pageNo, pageSize int) ([]model.PrintPriceConfigModel, int64, error) {
	list := make([]model.PrintPriceConfigModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(priceTable)
	if shopID > 0 {
		query = query.Where("print_shop_id", shopID)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
