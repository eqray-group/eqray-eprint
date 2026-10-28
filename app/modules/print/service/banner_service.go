package service

import (
	"context"
	"strings"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

// BannerInput 是轮播图新增 / 编辑入参；PrintShopID 为 0 表示全平台。
type BannerInput struct {
	PrintShopID types.ID
	Title       string
	ImageURL    string
	LinkType    string
	LinkURL     string
	Sort        int
	Status      int8
}

// BannerService 负责门店轮播图的维护与会员端展示。
type BannerService struct {
	banners *repository.PrintBannerRepository
	shops   *ShopService
}

// NewBannerService 构造轮播图服务。
func NewBannerService() *BannerService {
	return &BannerService{
		banners: repository.NewPrintBannerRepository(),
		shops:   NewShopService(),
	}
}

// Page 分页查询轮播图，shopID 为 0 时查全部。
func (s *BannerService) Page(ctx context.Context, shopID types.ID, pageNo, pageSize int) ([]model.PrintBannerModel, int64, error) {
	return s.banners.Paginate(ctx, shopID, pageNo, pageSize)
}

// Detail 查询轮播图详情。
func (s *BannerService) Detail(ctx context.Context, id types.ID) (*model.PrintBannerModel, error) {
	banner, err := s.banners.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if banner == nil {
		return nil, baseErrs.New(printErrs.BannerNotExists, id)
	}

	return banner, nil
}

// Create 新增轮播图。
func (s *BannerService) Create(ctx context.Context, in BannerInput) (types.ID, error) {
	normalized, err := s.validate(ctx, in)
	if err != nil {
		return 0, err
	}

	banner := &model.PrintBannerModel{
		PrintShopID: normalized.PrintShopID,
		Title:       normalized.Title,
		ImageURL:    normalized.ImageURL,
		LinkType:    normalized.LinkType,
		LinkURL:     normalized.LinkURL,
		Sort:        normalized.Sort,
		Status:      normalized.Status,
	}
	banner.ID = support.NewID()

	if err := s.banners.Create(ctx, banner); err != nil {
		return 0, err
	}

	return banner.ID, nil
}

// Update 编辑轮播图。
func (s *BannerService) Update(ctx context.Context, id types.ID, in BannerInput) error {
	banner, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	normalized, err := s.validate(ctx, in)
	if err != nil {
		return err
	}

	banner.PrintShopID = normalized.PrintShopID
	banner.Title = normalized.Title
	banner.ImageURL = normalized.ImageURL
	banner.LinkType = normalized.LinkType
	banner.LinkURL = normalized.LinkURL
	banner.Sort = normalized.Sort
	banner.Status = normalized.Status

	return s.banners.Update(ctx, banner)
}

// Delete 删除轮播图。
func (s *BannerService) Delete(ctx context.Context, id types.ID, updater string) error {
	if _, err := s.Detail(ctx, id); err != nil {
		return err
	}

	return s.banners.Delete(ctx, id, updater)
}

// ListForShop 返回门店可见轮播图：门店自有 + 全平台，只取启用。
func (s *BannerService) ListForShop(ctx context.Context, shopID types.ID) ([]model.PrintBannerModel, error) {
	if _, err := s.shops.Detail(ctx, shopID); err != nil {
		return nil, err
	}

	return s.banners.ListEnabled(ctx, shopID)
}

// validate 归一化并校验入参；跳转类型为空表示不跳转。
func (s *BannerService) validate(ctx context.Context, in BannerInput) (BannerInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.ImageURL = strings.TrimSpace(in.ImageURL)
	in.LinkURL = strings.TrimSpace(in.LinkURL)
	in.LinkType = strings.TrimSpace(in.LinkType)
	if in.LinkType == "" {
		in.LinkType = printEnums.BannerLinkNone
	}

	if in.ImageURL == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "请上传轮播图图片")
	}
	if in.LinkType != printEnums.BannerLinkNone &&
		in.LinkType != printEnums.BannerLinkPage &&
		in.LinkType != printEnums.BannerLinkURL {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "跳转类型不合法")
	}
	if in.LinkType != printEnums.BannerLinkNone && in.LinkURL == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "请填写跳转地址")
	}
	if in.LinkType == printEnums.BannerLinkNone {
		in.LinkURL = ""
	}
	if in.Status != printEnums.StatusEnabled && in.Status != printEnums.StatusDisabled {
		in.Status = printEnums.StatusEnabled
	}
	if in.Sort < 0 {
		in.Sort = 0
	}

	// print_shop_id = 0 表示全平台，不需要校验门店
	if in.PrintShopID > 0 {
		if _, err := s.shops.Detail(ctx, in.PrintShopID); err != nil {
			return in, err
		}
	}

	return in, nil
}
