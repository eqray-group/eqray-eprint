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

// MemberProfile 是 C 端「我的资料」的返回结构。
// 只暴露会员自己能看懂的字段；会员卡信息走 /member/card 单独查询。
type MemberProfile struct {
	ID        types.ID        `json:"id"`
	Mobile    string          `json:"mobile"`
	Nickname  string          `json:"nickname"`
	Avatar    string          `json:"avatar"`
	Name      string          `json:"name"`
	Sex       int8            `json:"sex"`
	Birthday  *types.DateTime `json:"birthday"`
	AreaID    int             `json:"areaId"`
	CreatedAt string          `json:"createdAt"`
}

// MemberProfileUpdateInput 是 C 端可自助修改的资料字段。
// 手机号需要走换绑验证码流程，状态由后台维护，都不在这里。
type MemberProfileUpdateInput struct {
	Nickname string
	Avatar   string
	Name     string
	Sex      int8
	Birthday *types.DateTime
	AreaID   int
}

// MemberProfileService 负责 C 端会员资料的自助查询与修改。
type MemberProfileService struct {
	members *repository.MemberRepository
}

// NewMemberProfileService 构造会员资料服务，装配会员仓储。
func NewMemberProfileService() *MemberProfileService {
	return &MemberProfileService{members: repository.NewMemberRepository()}
}

// Profile 组装当前会员的资料。
// 参数 memberID 为令牌里的会员 ID；返回资料结构，会员不存在或禁用时返回业务错误码。
func (s *MemberProfileService) Profile(ctx context.Context, memberID types.ID) (*MemberProfile, error) {
	m, err := s.mustGet(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return &MemberProfile{
		ID:        m.ID,
		Mobile:    m.Mobile,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		Name:      m.Name,
		Sex:       m.Sex,
		Birthday:  m.Birthday,
		AreaID:    m.AreaID,
		CreatedAt: m.CreatedAt.String(),
	}, nil
}

// Update 修改会员自助可改的字段，只写 nickname / avatar / name / sex / birthday / area_id。
func (s *MemberProfileService) Update(ctx context.Context, memberID types.ID, in MemberProfileUpdateInput) error {
	m, err := s.mustGet(ctx, memberID)
	if err != nil {
		return err
	}

	in.Nickname = strings.TrimSpace(in.Nickname)
	in.Avatar = strings.TrimSpace(in.Avatar)
	in.Name = strings.TrimSpace(in.Name)

	if in.Nickname == "" {
		return baseErrs.WithMsg(baseErrs.BadRequest, "昵称不能为空")
	}
	if len([]rune(in.Nickname)) > 32 {
		return baseErrs.WithMsg(baseErrs.BadRequest, "昵称不能超过 32 个字符")
	}
	if len(in.Avatar) > 255 {
		return baseErrs.WithMsg(baseErrs.BadRequest, "头像地址过长")
	}
	if in.Sex < 0 || in.Sex > 2 {
		return baseErrs.WithMsg(baseErrs.BadRequest, "性别取值不合法")
	}
	if in.AreaID < 0 {
		in.AreaID = 0
	}
	if in.Birthday != nil && !in.Birthday.IsZero() && in.Birthday.Time().After(types.Now().Time()) {
		return baseErrs.WithMsg(baseErrs.BadRequest, "生日不能晚于今天")
	}

	m.Nickname = in.Nickname
	m.Avatar = in.Avatar
	m.Name = in.Name
	m.Sex = in.Sex
	m.Birthday = in.Birthday
	m.AreaID = in.AreaID

	return s.members.UpdateProfile(ctx, m)
}

// mustGet 取当前会员并校验状态，禁用会员不允许访问任何自助接口。
func (s *MemberProfileService) mustGet(ctx context.Context, memberID types.ID) (*model.MemberUserModel, error) {
	if memberID == 0 {
		return nil, baseErrs.New(memberErrs.MemberNotExists, memberID)
	}

	m, err := s.members.FindByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, baseErrs.New(memberErrs.MemberNotExists, memberID)
	}
	if m.Status != baseEnums.CommonStatusEnabled {
		return nil, baseErrs.New(memberErrs.MemberDisabled)
	}

	return m, nil
}
