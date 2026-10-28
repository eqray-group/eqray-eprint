package service

import (
	"context"
	"strings"

	baseEnums "eqray-eprint/app/modules/base/enums"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/types"
	memberErrs "eqray-eprint/app/modules/member/errs"
	"eqray-eprint/app/modules/member/model"
	"eqray-eprint/app/modules/member/repository"
)

// MemberSaveInput 是新增 / 编辑会员的入参，由 Controller 从请求体转换而来。
// 会员不再持有积分、经验、标签与分组，这些字段已随体系下线移除。
type MemberSaveInput struct {
	Nickname string
	Mobile   string
	Avatar   string
	Name     string
	Sex      int8
	Birthday *types.DateTime
	AreaID   int
	Mark     string
	// Status 为空表示调用方没有传状态：新增按启用处理，编辑保持原状态
	Status *int8
}

// MemberListQuery 是会员列表的查询条件。
type MemberListQuery struct {
	Nickname string
	Mobile   string
	Status   *baseEnums.CommonStatusEnum
	PageNo   int
	PageSize int
}

// MemberService 负责会员档案的业务规则。
type MemberService struct {
	members *repository.MemberRepository
}

// NewMemberService 构造会员档案服务，装配会员仓储。
func NewMemberService() *MemberService {
	return &MemberService{members: repository.NewMemberRepository()}
}

// Page 分页查询会员。
func (s *MemberService) Page(ctx context.Context, q MemberListQuery) ([]model.MemberUserModel, int64, error) {
	return s.members.Paginate(ctx, repository.MemberFilter{
		Nickname: strings.TrimSpace(q.Nickname),
		Mobile:   strings.TrimSpace(q.Mobile),
		Status:   q.Status,
	}, q.PageNo, q.PageSize)
}

// Detail 查询会员详情，不存在时返回业务错误码。
func (s *MemberService) Detail(ctx context.Context, id types.ID) (*model.MemberUserModel, error) {
	m, err := s.members.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, baseErrs.New(memberErrs.MemberNotExists, id)
	}

	return m, nil
}

// EnsureActive 校验会员存在且为启用状态，供 C 端只读接口复用。
func (s *MemberService) EnsureActive(ctx context.Context, id types.ID) (*model.MemberUserModel, error) {
	m, err := s.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	if m.Status != baseEnums.CommonStatusEnabled {
		return nil, baseErrs.New(memberErrs.MemberDisabled)
	}

	return m, nil
}

// Create 新增会员：手机号唯一；未传状态时默认启用。
func (s *MemberService) Create(ctx context.Context, in MemberSaveInput) (types.ID, error) {
	if err := normalizeMemberInput(&in); err != nil {
		return 0, err
	}

	exists, err := s.members.ExistsByMobile(ctx, in.Mobile, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, baseErrs.New(memberErrs.MemberMobileExists, in.Mobile)
	}

	m := &model.MemberUserModel{
		Mobile:   in.Mobile,
		Nickname: in.Nickname,
		Avatar:   in.Avatar,
		Name:     in.Name,
		Sex:      in.Sex,
		Birthday: in.Birthday,
		AreaID:   in.AreaID,
		Mark:     in.Mark,
		Status:   memberStatusOrDefault(in.Status),
	}
	if err := s.members.Create(ctx, m); err != nil {
		return 0, err
	}

	return m.ID, nil
}

// Update 更新会员档案；手机号变更时排除自身做唯一校验。
func (s *MemberService) Update(ctx context.Context, id types.ID, in MemberSaveInput) error {
	if err := normalizeMemberInput(&in); err != nil {
		return err
	}

	m, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	exists, err := s.members.ExistsByMobile(ctx, in.Mobile, id)
	if err != nil {
		return err
	}
	if exists {
		return baseErrs.New(memberErrs.MemberMobileExists, in.Mobile)
	}

	m.Nickname = in.Nickname
	m.Mobile = in.Mobile
	m.Avatar = in.Avatar
	m.Name = in.Name
	m.Sex = in.Sex
	m.Birthday = in.Birthday
	m.AreaID = in.AreaID
	m.Mark = in.Mark
	// 未传状态时保留原状态，避免编辑漏传把会员静默改成禁用
	if in.Status != nil {
		m.Status = baseEnums.CommonStatusEnum(*in.Status)
	}

	return s.members.Update(ctx, m)
}

// UpdateStatus 启用 / 禁用会员，供列表开关直接调用。
func (s *MemberService) UpdateStatus(ctx context.Context, id types.ID, status int8) error {
	if !baseEnums.CommonStatusEnum(status).IsValid() {
		return baseErrs.WithMsg(baseErrs.BadRequest, "会员状态取值不合法")
	}

	m, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	m.Status = baseEnums.CommonStatusEnum(status)

	return s.members.UpdateStatus(ctx, m)
}

// Delete 软删除会员。
func (s *MemberService) Delete(ctx context.Context, id types.ID) error {
	m, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	return s.members.Delete(ctx, m)
}

// Addresses 查询会员收货地址（后台只读）。
func (s *MemberService) Addresses(ctx context.Context, id types.ID) ([]model.MemberAddressModel, error) {
	if _, err := s.Detail(ctx, id); err != nil {
		return nil, err
	}

	return s.members.Addresses(ctx, id)
}

// memberStatusOrDefault 把可选的会员状态转成落库值：调用方没传时默认启用。
// 参数 status：请求里的状态指针，nil 表示未传；返回值：可直接写入 model 的状态枚举。
func memberStatusOrDefault(status *int8) baseEnums.CommonStatusEnum {
	if status == nil {
		return baseEnums.CommonStatusEnabled
	}

	return baseEnums.CommonStatusEnum(*status)
}

// normalizeMemberInput 统一清洗入参并补齐默认值，供新增与编辑复用。
func normalizeMemberInput(in *MemberSaveInput) error {
	in.Nickname = strings.TrimSpace(in.Nickname)
	in.Mobile = strings.TrimSpace(in.Mobile)
	in.Avatar = strings.TrimSpace(in.Avatar)
	in.Name = strings.TrimSpace(in.Name)
	in.Mark = strings.TrimSpace(in.Mark)

	if in.Nickname == "" {
		return baseErrs.WithMsg(baseErrs.BadRequest, "昵称不能为空")
	}
	if in.Mobile == "" {
		return baseErrs.WithMsg(baseErrs.BadRequest, "手机号不能为空")
	}
	// 状态可以不传，传了才校验取值范围
	if in.Status != nil && !baseEnums.CommonStatusEnum(*in.Status).IsValid() {
		return baseErrs.WithMsg(baseErrs.BadRequest, "会员状态取值不合法")
	}
	if in.Sex < 0 || in.Sex > 2 {
		return baseErrs.WithMsg(baseErrs.BadRequest, "性别取值不合法")
	}
	if in.AreaID < 0 {
		in.AreaID = 0
	}

	return nil
}
