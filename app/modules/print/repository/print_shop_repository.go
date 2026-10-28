package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/model"
)

const shopTable = model.PrintShopTable

// ShopFilter 是门店列表的筛选条件，零值表示不限制。
type ShopFilter struct {
	Name        string
	Status      *int8
	AuditStatus *int8
}

// PrintShopRepository 是 print_shop 表的仓储。
type PrintShopRepository struct{}

// NewPrintShopRepository 构造门店仓储。
func NewPrintShopRepository() *PrintShopRepository {
	return &PrintShopRepository{}
}

// Create 新增门店。
func (r *PrintShopRepository) Create(ctx context.Context, m *model.PrintShopModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(shopTable).Insert(r.values(m))

	return err
}

// Update 更新门店基础信息（不含审核状态）。
func (r *PrintShopRepository) Update(ctx context.Context, m *model.PrintShopModel) error {
	m.FillUpdate(ctx)

	values := r.values(m)
	delete(values, "id")
	delete(values, "creator")
	delete(values, "created_at")

	_, err := facades.DB().WithContext(ctx).Table(shopTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(values)

	return err
}

// UpdateStatus 只改营业状态，供列表开关使用。
func (r *PrintShopRepository) UpdateStatus(ctx context.Context, id types.ID, status int8, updater string) error {
	_, err := facades.DB().WithContext(ctx).Table(shopTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"status":     status,
			"updater":    updater,
			"updated_at": types.Now(),
		})

	return err
}

// UpdateAudit 更新审核状态与意见。
func (r *PrintShopRepository) UpdateAudit(ctx context.Context, id types.ID, auditStatus int8, remark, updater string) error {
	_, err := facades.DB().WithContext(ctx).Table(shopTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"audit_status": auditStatus,
			"audit_remark": remark,
			"updater":      updater,
			"updated_at":   types.Now(),
		})

	return err
}

// Delete 软删除门店；历史订单仍按门店 ID 追溯。
func (r *PrintShopRepository) Delete(ctx context.Context, id types.ID, updater string) error {
	_, err := facades.DB().WithContext(ctx).Table(shopTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"deleted_at": types.Now(),
			"updater":    updater,
			"updated_at": types.Now(),
		})

	return err
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintShopRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintShopModel, error) {
	var m model.PrintShopModel

	if err := facades.DB().WithContext(ctx).Table(shopTable).
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

// Paginate 按条件分页查询门店。
func (r *PrintShopRepository) Paginate(ctx context.Context, filter ShopFilter, pageNo, pageSize int) ([]model.PrintShopModel, int64, error) {
	list := make([]model.PrintShopModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(shopTable).WhereNull("deleted_at")
	if filter.Name != "" {
		query = query.WhereLike("name", "%"+filter.Name+"%")
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}
	if filter.AuditStatus != nil {
		query = query.Where("audit_status", *filter.AuditStatus)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// ListOpen 查询会员端可下单的门店：审核通过且营业中。
func (r *PrintShopRepository) ListOpen(ctx context.Context) ([]model.PrintShopModel, error) {
	list := make([]model.PrintShopModel, 0)

	if err := facades.DB().WithContext(ctx).Table(shopTable).
		Where("status", printEnums.ShopStatusOpen).
		Where("audit_status", printEnums.AuditStatusApproved).
		WhereNull("deleted_at").
		OrderBy("id", "desc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// values 把模型转成写入字段集合；business_hours 为空时写 NULL。
func (r *PrintShopRepository) values(m *model.PrintShopModel) map[string]any {
	var hours any
	if m.BusinessHours != nil && *m.BusinessHours != "" {
		hours = *m.BusinessHours
	}

	return map[string]any{
		"id":             m.ID,
		"name":           m.Name,
		"address":        m.Address,
		"longitude":      m.Longitude,
		"latitude":       m.Latitude,
		"avatar":         m.Avatar,
		"business_hours": hours,
		"status":         m.Status,
		"audit_status":   m.AuditStatus,
		"audit_remark":   m.AuditRemark,
		"remark":         m.Remark,
		"creator":        m.Creator,
		"created_at":     m.CreatedAt,
		"updater":        m.Updater,
		"updated_at":     m.UpdatedAt,
	}
}
