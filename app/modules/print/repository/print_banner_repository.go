package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/model"
)

const bannerTable = model.PrintBannerTable

// PrintBannerRepository 是 print_banner 表的仓储。
type PrintBannerRepository struct{}

// NewPrintBannerRepository 构造轮播图仓储。
func NewPrintBannerRepository() *PrintBannerRepository {
	return &PrintBannerRepository{}
}

// Create 新增轮播图。
func (r *PrintBannerRepository) Create(ctx context.Context, m *model.PrintBannerModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(bannerTable).Insert(map[string]any{
		"id":            m.ID,
		"print_shop_id": m.PrintShopID,
		"title":         m.Title,
		"image_url":     m.ImageURL,
		"link_type":     m.LinkType,
		"link_url":      m.LinkURL,
		"sort":          m.Sort,
		"status":        m.Status,
		"creator":       m.Creator,
		"created_at":    m.CreatedAt,
		"updater":       m.Updater,
		"updated_at":    m.UpdatedAt,
	})

	return err
}

// Update 编辑轮播图。
func (r *PrintBannerRepository) Update(ctx context.Context, m *model.PrintBannerModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(bannerTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"print_shop_id": m.PrintShopID,
			"title":         m.Title,
			"image_url":     m.ImageURL,
			"link_type":     m.LinkType,
			"link_url":      m.LinkURL,
			"sort":          m.Sort,
			"status":        m.Status,
			"updater":       m.Updater,
			"updated_at":    m.UpdatedAt,
		})

	return err
}

// Delete 软删除轮播图。
func (r *PrintBannerRepository) Delete(ctx context.Context, id types.ID, updater string) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(bannerTable).
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
func (r *PrintBannerRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintBannerModel, error) {
	var m model.PrintBannerModel

	if err := facades.DB().WithContext(ctx).Table(bannerTable).
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

// Paginate 按门店分页查询；shopID 为 0 表示不限门店。
func (r *PrintBannerRepository) Paginate(ctx context.Context, shopID types.ID, pageNo, pageSize int) ([]model.PrintBannerModel, int64, error) {
	list := make([]model.PrintBannerModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(bannerTable).WhereNull("deleted_at")
	if shopID > 0 {
		query = query.Where("print_shop_id", shopID)
	}

	if err := query.OrderBy("sort", "asc").OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// ListEnabled 查询门店可见的轮播图：门店自有 + 全平台（shop_id=0）。
func (r *PrintBannerRepository) ListEnabled(ctx context.Context, shopID types.ID) ([]model.PrintBannerModel, error) {
	list := make([]model.PrintBannerModel, 0)

	if err := facades.DB().WithContext(ctx).Table(bannerTable).
		WhereIn("print_shop_id", []any{0, shopID}).
		Where("status", printEnums.StatusEnabled).
		WhereNull("deleted_at").
		OrderBy("sort", "asc").
		OrderByDesc("id").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}
