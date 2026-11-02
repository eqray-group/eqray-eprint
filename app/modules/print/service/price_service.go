package service

import (
	"context"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

// PriceInput 是价格规则新增 / 编辑入参，价格单位为分/页。
type PriceInput struct {
	PrintShopID types.ID
	PaperType   string
	ColorType   string
	PrintMode   string
	Price       int64
	Status      int8
}

// PriceBatchItem 是 SKU 矩阵里单个组合的保存项，价格单位为分/页。
// Price 为 0 表示该组合不售卖，已配置的规则会被删除。
type PriceBatchItem struct {
	PaperType string
	ColorType string
	PrintMode string
	Price     int64
	Status    int8
}

// PriceListQuery 是价格列表查询条件。
type PriceListQuery struct {
	ShopID   types.ID
	PageNo   int
	PageSize int
}

// PriceService 负责打印价格规则的维护。
// 保存时只做软校验提示，不阻断：门店常见的顺序是先配价格、后装打印机。
type PriceService struct {
	prices     *repository.PrintPriceConfigRepository
	shops      *ShopService
	capability *CapabilityService
}

// NewPriceService 构造价格服务。
func NewPriceService() *PriceService {
	return &PriceService{
		prices:     repository.NewPrintPriceConfigRepository(),
		shops:      NewShopService(),
		capability: NewCapabilityService(),
	}
}

// Page 分页查询价格规则。
func (s *PriceService) Page(ctx context.Context, q PriceListQuery) ([]model.PrintPriceConfigModel, int64, error) {
	return s.prices.Paginate(ctx, q.ShopID, q.PageNo, q.PageSize)
}

// ListByShop 查询门店全部价格规则（含停用）。
func (s *PriceService) ListByShop(ctx context.Context, shopID types.ID) ([]model.PrintPriceConfigModel, error) {
	return s.prices.ListByShop(ctx, shopID)
}

// Create 新增价格规则，同一组合不允许重复配置。
func (s *PriceService) Create(ctx context.Context, in PriceInput) (types.ID, error) {
	normalized, err := s.validate(ctx, in)
	if err != nil {
		return 0, err
	}

	exists, err := s.prices.FindByRule(ctx, normalized.PrintShopID, normalized.PaperType, normalized.ColorType, normalized.PrintMode)
	if err != nil {
		return 0, err
	}
	if exists != nil {
		return 0, baseErrs.New(printErrs.PriceRuleExists)
	}

	rule := &model.PrintPriceConfigModel{
		PrintShopID: normalized.PrintShopID,
		PaperType:   normalized.PaperType,
		ColorType:   normalized.ColorType,
		PrintMode:   normalized.PrintMode,
		Price:       normalized.Price,
		Status:      normalized.Status,
	}
	rule.ID = support.NewID()

	if err := s.prices.Create(ctx, rule); err != nil {
		return 0, err
	}

	return rule.ID, nil
}

// Update 修改价格与启用状态；组合本身不可改（改组合等于换一条规则）。
func (s *PriceService) Update(ctx context.Context, id types.ID, price int64, status int8) error {
	rule, err := s.prices.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if rule == nil {
		return baseErrs.New(printErrs.PriceNotExists)
	}
	if price <= 0 {
		return baseErrs.New(printErrs.PriceInvalid)
	}
	if status != printEnums.StatusEnabled && status != printEnums.StatusDisabled {
		return baseErrs.WithMsg(baseErrs.BadRequest, "状态取值不合法")
	}

	rule.Price = price
	rule.Status = status

	return s.prices.Update(ctx, rule)
}

// SaveBatch 批量保存门店价格（SKU 矩阵页一次提交）：price > 0 的组合有则改价、无则新增，
// price = 0 视为该组合不售卖并删除已有规则；整批在一个事务里完成，避免只保存一半。
// 参数 shopID 为门店，items 为矩阵行；同一组合重复提交时以最后一条为准。
func (s *PriceService) SaveBatch(ctx context.Context, shopID types.ID, items []PriceBatchItem) error {
	if _, err := s.shops.Detail(ctx, shopID); err != nil {
		return err
	}

	// 组合做归一化并按「纸张|色彩|单双面」去重，避免同一次提交里重复写同一条
	unique := make(map[string]PriceBatchItem, len(items))
	order := make([]string, 0, len(items))
	for _, item := range items {
		item.PaperType = normalizePaper(item.PaperType)
		item.ColorType = normalizeColor(item.ColorType)
		item.PrintMode = normalizeMode(item.PrintMode)
		if item.PaperType == "" || item.ColorType == "" || item.PrintMode == "" {
			return baseErrs.WithMsg(baseErrs.BadRequest, "纸张、色彩、单双面取值不合法")
		}
		if item.Price < 0 {
			return baseErrs.New(printErrs.PriceInvalid)
		}
		if item.Status != printEnums.StatusEnabled && item.Status != printEnums.StatusDisabled {
			item.Status = printEnums.StatusEnabled
		}

		key := item.PaperType + "|" + item.ColorType + "|" + item.PrintMode
		if _, ok := unique[key]; !ok {
			order = append(order, key)
		}
		unique[key] = item
	}

	if len(order) == 0 {
		return nil
	}

	return facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		for _, key := range order {
			item := unique[key]

			existing, err := s.prices.FindByRule(ctx, shopID, item.PaperType, item.ColorType, item.PrintMode)
			if err != nil {
				return err
			}

			// 价格清空表示这个组合不卖，已配置就删掉
			if item.Price == 0 {
				if existing == nil {
					continue
				}
				if err := s.prices.DeleteInTx(ctx, tx, existing.ID); err != nil {
					return err
				}
				continue
			}

			if existing == nil {
				rule := &model.PrintPriceConfigModel{
					PrintShopID: shopID,
					PaperType:   item.PaperType,
					ColorType:   item.ColorType,
					PrintMode:   item.PrintMode,
					Price:       item.Price,
					Status:      item.Status,
				}
				rule.ID = support.NewID()

				if err := s.prices.CreateInTx(ctx, tx, rule); err != nil {
					return err
				}
				continue
			}

			existing.Price = item.Price
			existing.Status = item.Status
			if err := s.prices.UpdateInTx(ctx, tx, existing); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 物理删除价格规则。
func (s *PriceService) Delete(ctx context.Context, id types.ID) error {
	return s.prices.Delete(ctx, id)
}

// Warnings 返回该门店「有价格但无打印机支持」的组合，保存后展示给门店核对。
func (s *PriceService) Warnings(ctx context.Context, shopID types.ID) ([]CapabilityWarning, error) {
	return s.capability.Warnings(ctx, shopID)
}

// validate 归一化并校验：门店存在、枚举合法、单价大于 0。
func (s *PriceService) validate(ctx context.Context, in PriceInput) (PriceInput, error) {
	in.PaperType = normalizePaper(in.PaperType)
	in.ColorType = normalizeColor(in.ColorType)
	in.PrintMode = normalizeMode(in.PrintMode)

	if in.Price <= 0 {
		return in, baseErrs.New(printErrs.PriceInvalid)
	}
	if in.PaperType == "" || in.ColorType == "" || in.PrintMode == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "纸张、色彩、单双面取值不合法")
	}
	if in.Status != printEnums.StatusEnabled && in.Status != printEnums.StatusDisabled {
		in.Status = printEnums.StatusEnabled
	}
	if _, err := s.shops.Detail(ctx, in.PrintShopID); err != nil {
		return in, err
	}

	return in, nil
}

// normalizePaper 归一化纸张枚举。
func normalizePaper(value string) string {
	switch value {
	case printEnums.PaperA4, printEnums.PaperA3:
		return value
	default:
		return ""
	}
}

// normalizeColor 归一化色彩枚举。
func normalizeColor(value string) string {
	switch value {
	case printEnums.ColorBW, printEnums.ColorColor:
		return value
	default:
		return ""
	}
}

// normalizeMode 归一化单双面枚举。
func normalizeMode(value string) string {
	switch value {
	case printEnums.ModeSingle, printEnums.ModeDouble:
		return value
	default:
		return ""
	}
}
