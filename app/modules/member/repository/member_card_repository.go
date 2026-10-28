package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/member/model"
)

const memberCardTable = model.MemberCardTable

// MemberCardFilter 是会员卡列表的筛选条件，零值表示不限制。
type MemberCardFilter struct {
	CardNo string
	Status *int8
}

// MemberCardRepository 是 member_card 表的仓储。
// 会员卡不做软删除，停用与过期都改 status，因此查询不带 deleted_at 过滤。
type MemberCardRepository struct{}

// NewMemberCardRepository 构造会员卡仓储。
func NewMemberCardRepository() *MemberCardRepository {
	return &MemberCardRepository{}
}

// Create 插入会员卡；卡号唯一冲突由上游转换为业务错误码。
func (r *MemberCardRepository) Create(ctx context.Context, m *model.MemberCardModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(memberCardTable).Insert(map[string]any{
		"id":               m.ID,
		"member_user_id":   m.MemberUserID,
		"card_no":          m.CardNo,
		"card_level":       m.CardLevel,
		"card_name":        m.CardName,
		"discount_percent": m.DiscountPercent,
		"card_type":        m.CardType,
		"source":           m.Source,
		"status":           m.Status,
		"start_time":       m.StartTime,
		"expire_time":      m.ExpireTime,
		"remark":           m.Remark,
		"creator":          m.Creator,
		"created_at":       m.CreatedAt,
		"updater":          m.Updater,
		"updated_at":       m.UpdatedAt,
	})

	return err
}

// Update 更新会员卡的等级、折扣、状态与有效期；卡号不随编辑变化。
func (r *MemberCardRepository) Update(ctx context.Context, m *model.MemberCardModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(memberCardTable).
		Where("id", m.ID).
		Update(map[string]any{
			"card_level":       m.CardLevel,
			"card_name":        m.CardName,
			"discount_percent": m.DiscountPercent,
			"status":           m.Status,
			"start_time":       m.StartTime,
			"expire_time":      m.ExpireTime,
			"remark":           m.Remark,
			"updater":          m.Updater,
			"updated_at":       m.UpdatedAt,
		})

	return err
}

// UpdateStatus 只更新卡状态（停用 / 恢复正常）。
func (r *MemberCardRepository) UpdateStatus(ctx context.Context, id types.ID, status int8, updater string) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(memberCardTable).
		Where("id", id).
		Update(map[string]any{
			"status":     status,
			"updater":    updater,
			"updated_at": now,
		})

	return err
}

// FindByMemberID 按会员查询唯一的会员卡；未命中返回 (nil, nil)。
func (r *MemberCardRepository) FindByMemberID(ctx context.Context, memberUserID types.ID) (*model.MemberCardModel, error) {
	var m model.MemberCardModel

	if err := facades.DB().WithContext(ctx).Table(memberCardTable).
		Where("member_user_id", memberUserID).
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// Paginate 按条件分页查询会员卡。
func (r *MemberCardRepository) Paginate(ctx context.Context, filter MemberCardFilter, pageNo, pageSize int) ([]model.MemberCardModel, int64, error) {
	list := make([]model.MemberCardModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(memberCardTable)
	if filter.CardNo != "" {
		query = query.WhereLike("card_no", "%"+filter.CardNo+"%")
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
