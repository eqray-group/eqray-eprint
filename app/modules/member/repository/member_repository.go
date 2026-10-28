package repository

import (
	"context"

	"eqray-eprint/app/facades"
	baseEnums "eqray-eprint/app/modules/base/enums"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	memberErrs "eqray-eprint/app/modules/member/errs"
	"eqray-eprint/app/modules/member/model"
)

const memberTable = model.MemberUserTable

// MemberFilter 是会员列表的筛选条件，零值表示不限制。
type MemberFilter struct {
	Nickname string
	Mobile   string
	Status   *baseEnums.CommonStatusEnum
}

// MemberRepository 是 member_user 表的仓储。
type MemberRepository struct{}

// NewMemberRepository 构造会员仓储。
func NewMemberRepository() *MemberRepository {
	return &MemberRepository{}
}

// Create 插入会员；手机号冲突转换为统一的业务错误码。
func (r *MemberRepository) Create(ctx context.Context, m *model.MemberUserModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(memberTable).Insert(map[string]any{
		"id":          m.ID,
		"mobile":      m.Mobile,
		"password":    m.Password,
		"nickname":    m.Nickname,
		"avatar":      m.Avatar,
		"name":        m.Name,
		"sex":         m.Sex,
		"birthday":    m.Birthday,
		"area_id":     m.AreaID,
		"mark":        m.Mark,
		"register_ip": m.RegisterIP,
		"login_ip":    m.LoginIP,
		"login_date":  m.LoginDate,
		"status":      m.Status,
		"creator":     m.Creator,
		"created_at":  m.CreatedAt,
		"updater":     m.Updater,
		"updated_at":  m.UpdatedAt,
	})
	if err != nil && support.IsDuplicate(err) {
		return baseErrs.New(memberErrs.MemberMobileExists, m.Mobile)
	}

	return err
}

// Update 更新会员可变字段；手机号冲突同样转换为业务错误码。
func (r *MemberRepository) Update(ctx context.Context, m *model.MemberUserModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(memberTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"mobile":     m.Mobile,
			"nickname":   m.Nickname,
			"avatar":     m.Avatar,
			"name":       m.Name,
			"sex":        m.Sex,
			"birthday":   m.Birthday,
			"area_id":    m.AreaID,
			"mark":       m.Mark,
			"status":     m.Status,
			"updater":    m.Updater,
			"updated_at": m.UpdatedAt,
		})
	if err != nil && support.IsDuplicate(err) {
		return baseErrs.New(memberErrs.MemberMobileExists, m.Mobile)
	}

	return err
}

// UpdateProfile 只更新会员自助可改的字段。
// 手机号、标记、状态属于后台或独立流程，不走这里，避免越权覆盖。
func (r *MemberRepository) UpdateProfile(ctx context.Context, m *model.MemberUserModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(memberTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"nickname":   m.Nickname,
			"avatar":     m.Avatar,
			"name":       m.Name,
			"sex":        m.Sex,
			"birthday":   m.Birthday,
			"area_id":    m.AreaID,
			"updater":    m.Updater,
			"updated_at": m.UpdatedAt,
		})

	return err
}

// UpdateStatus 只更新启用状态，列表切换开关时避免整行覆盖。
func (r *MemberRepository) UpdateStatus(ctx context.Context, m *model.MemberUserModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(memberTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"status":     m.Status,
			"updater":    m.Updater,
			"updated_at": m.UpdatedAt,
		})

	return err
}

// UpdateLoginInfo 记录会员最近一次登录，只更新登录相关字段。
func (r *MemberRepository) UpdateLoginInfo(ctx context.Context, id types.ID, ip string, at types.DateTime) error {
	_, err := facades.DB().WithContext(ctx).Table(memberTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"login_ip":   ip,
			"login_date": at,
			"updated_at": types.Now(),
		})

	return err
}

// Delete 软删除会员。
func (r *MemberRepository) Delete(ctx context.Context, m *model.MemberUserModel) error {
	m.FillUpdate(ctx)

	res, err := facades.DB().WithContext(ctx).Table(memberTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"deleted_at": types.Now(),
			"updater":    m.Updater,
			"updated_at": m.UpdatedAt,
		})
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return baseErrs.New(memberErrs.MemberNotExists, m.ID)
	}

	return nil
}

// FindByID 按主键查询未软删的会员；未命中返回 (nil, nil)。
func (r *MemberRepository) FindByID(ctx context.Context, id types.ID) (*model.MemberUserModel, error) {
	var m model.MemberUserModel

	if err := facades.DB().WithContext(ctx).Table(memberTable).
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

// FindByMobile 按手机号查询未软删的会员；未命中返回 (nil, nil)。
func (r *MemberRepository) FindByMobile(ctx context.Context, mobile string) (*model.MemberUserModel, error) {
	var m model.MemberUserModel

	if err := facades.DB().WithContext(ctx).Table(memberTable).
		Where("mobile", mobile).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// ExistsByMobile 判断手机号是否已被占用，excludeID 用于编辑时排除自身。
func (r *MemberRepository) ExistsByMobile(ctx context.Context, mobile string, excludeID types.ID) (bool, error) {
	query := facades.DB().WithContext(ctx).Table(memberTable).
		Where("mobile", mobile).
		WhereNull("deleted_at")
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	total, err := query.Count()
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// Paginate 按条件分页查询会员。
func (r *MemberRepository) Paginate(ctx context.Context, filter MemberFilter, pageNo, pageSize int) ([]model.MemberUserModel, int64, error) {
	list := make([]model.MemberUserModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(memberTable).WhereNull("deleted_at")
	if filter.Nickname != "" {
		query = query.WhereLike("nickname", "%"+filter.Nickname+"%")
	}
	if filter.Mobile != "" {
		query = query.WhereLike("mobile", "%"+filter.Mobile+"%")
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// Addresses 查询会员收货地址，数量有限不分页。
func (r *MemberRepository) Addresses(ctx context.Context, userID types.ID) ([]model.MemberAddressModel, error) {
	list := make([]model.MemberAddressModel, 0)

	if err := facades.DB().WithContext(ctx).Table(model.MemberAddressTable).
		Where("user_id", userID).
		WhereNull("deleted_at").
		OrderByDesc("default_status").
		OrderByDesc("id").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}
