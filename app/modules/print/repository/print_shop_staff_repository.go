package repository

import (
	"context"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/model"
)

const shopStaffTable = model.PrintShopStaffTable

// PrintShopStaffRepository 是 print_shop_staff 表的仓储。
// 表无软删，解绑即物理删除，重新绑定不会撞唯一键。
type PrintShopStaffRepository struct{}

// NewPrintShopStaffRepository 构造门店成员仓储。
func NewPrintShopStaffRepository() *PrintShopStaffRepository {
	return &PrintShopStaffRepository{}
}

// Create 绑定账号到门店；重复绑定返回数据库错误，由服务层转业务错误。
func (r *PrintShopStaffRepository) Create(ctx context.Context, m *model.PrintShopStaffModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(shopStaffTable).Insert(map[string]any{
		"id":         m.ID,
		"user_id":    m.UserID,
		"shop_id":    m.ShopID,
		"role":       m.Role,
		"status":     m.Status,
		"creator":    m.Creator,
		"created_at": m.CreatedAt,
		"updater":    m.Updater,
		"updated_at": m.UpdatedAt,
	})

	return err
}

// Delete 解除绑定（物理删除），未命中返回 false。
func (r *PrintShopStaffRepository) Delete(ctx context.Context, id types.ID) error {
	res, err := facades.DB().WithContext(ctx).Table(shopStaffTable).Where("id", id).Delete()
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return baseErrs.WithMsg(baseErrs.NotFound, "门店成员不存在")
	}

	return nil
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintShopStaffRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintShopStaffModel, error) {
	var m model.PrintShopStaffModel

	if err := facades.DB().WithContext(ctx).Table(shopStaffTable).Where("id", id).First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindByUserShop 查询某账号在指定门店的绑定，用于鉴权与门店归属推导。
func (r *PrintShopStaffRepository) FindByUserShop(ctx context.Context, userID, shopID types.ID) (*model.PrintShopStaffModel, error) {
	var m model.PrintShopStaffModel

	if err := facades.DB().WithContext(ctx).Table(shopStaffTable).
		Where("user_id", userID).
		Where("shop_id", shopID).
		Where("status", 1).
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// ListByUser 查询账号可用的门店绑定，客户端登录后据此选门店。
func (r *PrintShopStaffRepository) ListByUser(ctx context.Context, userID types.ID) ([]model.PrintShopStaffModel, error) {
	list := make([]model.PrintShopStaffModel, 0)

	if err := facades.DB().WithContext(ctx).Table(shopStaffTable).
		Where("user_id", userID).
		Where("status", 1).
		OrderBy("id", "asc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// ListByShop 查询门店下的全部成员。
func (r *PrintShopStaffRepository) ListByShop(ctx context.Context, shopID types.ID) ([]model.PrintShopStaffModel, error) {
	list := make([]model.PrintShopStaffModel, 0)

	if err := facades.DB().WithContext(ctx).Table(shopStaffTable).
		Where("shop_id", shopID).
		OrderBy("id", "asc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// Exists 判断账号是否已绑定该门店。
func (r *PrintShopStaffRepository) Exists(ctx context.Context, userID, shopID types.ID) (bool, error) {
	total, err := facades.DB().WithContext(ctx).Table(shopStaffTable).
		Where("user_id", userID).
		Where("shop_id", shopID).
		Count()
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// NewID 暴露雪花主键生成，供服务层创建绑定记录时使用。
func NewID() types.ID {
	return support.NewID()
}
