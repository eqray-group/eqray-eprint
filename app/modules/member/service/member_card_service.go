package service

import (
	"context"
	"fmt"
	"strings"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	memberErrs "eqray-eprint/app/modules/member/errs"
	"eqray-eprint/app/modules/member/model"
	"eqray-eprint/app/modules/member/repository"
)

// MemberCardDefaultDiscount 是无卡或卡不可用时使用的折扣百分比，100 表示原价。
const MemberCardDefaultDiscount = 100

// MemberCardInput 是后台发放 / 修改会员卡的入参。
type MemberCardInput struct {
	MemberUserID    types.ID
	CardLevel       string
	CardName        string
	DiscountPercent int
	Status          int8
	StartTime       *types.DateTime
	ExpireTime      *types.DateTime
	Remark          string
}

// MemberCardListQuery 是会员卡列表的查询条件。
type MemberCardListQuery struct {
	CardNo   string
	Status   *int8
	PageNo   int
	PageSize int
}

// MemberCardService 负责会员卡的发放、修改与折扣计算。
// 会员卡是与会员等级无关的独立权益载体：折扣只读卡上的 DiscountPercent。
type MemberCardService struct {
	cards   *repository.MemberCardRepository
	members *repository.MemberRepository
}

// NewMemberCardService 构造会员卡服务，装配会员卡与会员仓储。
func NewMemberCardService() *MemberCardService {
	return &MemberCardService{
		cards:   repository.NewMemberCardRepository(),
		members: repository.NewMemberRepository(),
	}
}

// IssueDefault 为会员签发默认卡（无折扣），注册流程与后台新增会员时调用。
// 已有卡时直接返回 nil，保证幂等；参数 memberID 为会员 ID，返回写入错误。
func (s *MemberCardService) IssueDefault(ctx context.Context, memberID types.ID) error {
	existing, err := s.cards.FindByMemberID(ctx, memberID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	now := types.Now()
	card := &model.MemberCardModel{
		MemberUserID:    memberID,
		CardLevel:       "normal",
		CardName:        "普通会员卡",
		DiscountPercent: MemberCardDefaultDiscount,
		CardType:        "level",
		Source:          model.MemberCardSourceRegister,
		Status:          model.MemberCardStatusNormal,
		StartTime:       &now,
	}
	card.ID = support.NewID()
	card.CardNo = newCardNo(card.ID)

	return s.cards.Create(ctx, card)
}

// Save 后台发放或修改会员卡：会员已有卡时更新，没有卡时新建。
// 入参为 MemberCardInput；返回会员卡 ID 与写入错误。
func (s *MemberCardService) Save(ctx context.Context, in MemberCardInput) (types.ID, error) {
	if in.DiscountPercent < 1 || in.DiscountPercent > 100 {
		return 0, baseErrs.New(memberErrs.MemberCardDiscountInvalid, in.DiscountPercent)
	}

	member, err := s.members.FindByID(ctx, in.MemberUserID)
	if err != nil {
		return 0, err
	}
	if member == nil {
		return 0, baseErrs.New(memberErrs.MemberNotExists, in.MemberUserID)
	}

	in.CardLevel = strings.TrimSpace(in.CardLevel)
	in.CardName = strings.TrimSpace(in.CardName)
	if in.CardLevel == "" || in.CardName == "" {
		return 0, baseErrs.WithMsg(baseErrs.BadRequest, "卡等级编码与卡名称不能为空")
	}

	existing, err := s.cards.FindByMemberID(ctx, in.MemberUserID)
	if err != nil {
		return 0, err
	}

	if existing != nil {
		existing.CardLevel = in.CardLevel
		existing.CardName = in.CardName
		existing.DiscountPercent = in.DiscountPercent
		existing.Status = in.Status
		existing.StartTime = in.StartTime
		existing.ExpireTime = in.ExpireTime
		existing.Remark = in.Remark

		return existing.ID, s.cards.Update(ctx, existing)
	}

	card := &model.MemberCardModel{
		MemberUserID:    in.MemberUserID,
		CardLevel:       in.CardLevel,
		CardName:        in.CardName,
		DiscountPercent: in.DiscountPercent,
		CardType:        "level",
		Source:          model.MemberCardSourceAdmin,
		Status:          in.Status,
		StartTime:       in.StartTime,
		ExpireTime:      in.ExpireTime,
		Remark:          in.Remark,
	}
	card.ID = support.NewID()
	card.CardNo = newCardNo(card.ID)

	return card.ID, s.cards.Create(ctx, card)
}

// Detail 查询会员的卡，不存在时返回业务错误码。
func (s *MemberCardService) Detail(ctx context.Context, memberID types.ID) (*model.MemberCardModel, error) {
	card, err := s.cards.FindByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, baseErrs.New(memberErrs.MemberCardNotExists, memberID)
	}

	card.Enabled = card.IsUsable()

	return card, nil
}

// Page 分页查询会员卡，返回列表与总数；每条记录带上「当前是否可用」标记。
func (s *MemberCardService) Page(ctx context.Context, q MemberCardListQuery) ([]model.MemberCardModel, int64, error) {
	list, total, err := s.cards.Paginate(ctx, repository.MemberCardFilter{
		CardNo: strings.TrimSpace(q.CardNo),
		Status: q.Status,
	}, q.PageNo, q.PageSize)
	if err != nil {
		return nil, 0, err
	}

	for i := range list {
		list[i].Enabled = list[i].IsUsable()
	}

	return list, total, nil
}

// Disable 停用会员卡：停用后折扣立即失效。
// 参数 memberID 为会员 ID；返回写入错误。
func (s *MemberCardService) Disable(ctx context.Context, memberID types.ID) error {
	card, err := s.cards.FindByMemberID(ctx, memberID)
	if err != nil {
		return err
	}
	if card == nil {
		return baseErrs.New(memberErrs.MemberCardNotExists, memberID)
	}

	return s.cards.UpdateStatus(ctx, card.ID, model.MemberCardStatusDisabled, "system")
}

// DiscountPercent 返回会员当前可享的折扣百分比，供计价使用。
// 无卡、卡停用或已过期一律返回 100（原价）；参数 memberID 为会员 ID。
func (s *MemberCardService) DiscountPercent(ctx context.Context, memberID types.ID) (int, error) {
	card, err := s.cards.FindByMemberID(ctx, memberID)
	if err != nil {
		return MemberCardDefaultDiscount, err
	}
	if !card.IsUsable() {
		return MemberCardDefaultDiscount, nil
	}
	if card.DiscountPercent < 1 || card.DiscountPercent > 100 {
		return MemberCardDefaultDiscount, nil
	}

	return card.DiscountPercent, nil
}

// newCardNo 按主键生成卡号，形如 C1234567890123456；主键唯一因此卡号唯一。
func newCardNo(id types.ID) string {
	return fmt.Sprintf("C%d", int64(id))
}
