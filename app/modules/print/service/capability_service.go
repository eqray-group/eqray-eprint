package service

import (
	"context"
	"strings"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

// SellableCombo 是可下单的参数组合：价格配置与打印机能力的交集。
type SellableCombo struct {
	PaperType   string     `json:"paperType"`
	ColorType   string     `json:"colorType"`
	PrintMode   string     `json:"printMode"`
	Price       int64      `json:"price"`
	PrinterIDs  []types.ID `json:"printerIds"`
	PrinterName string     `json:"printerName"`
}

// CapabilityWarning 是能力体检的提示项，不阻断保存，只在后台展示。
type CapabilityWarning struct {
	PaperType string `json:"paperType"`
	ColorType string `json:"colorType"`
	PrintMode string `json:"printMode"`
	Message   string `json:"message"`
}

// CapabilityService 负责「能不能卖」与「用哪台打印机」两个判定。
// 口径：可售 = 启用价格规则 ∩ 启用打印机能力，运行时求交集；
// 可售只看 status，不看在线状态，避免打印机断网导致商品全部下架。
type CapabilityService struct {
	prices   *repository.PrintPriceConfigRepository
	printers *repository.PrintPrinterRepository
	devices  *repository.PrintDeviceRepository
}

// NewCapabilityService 构造能力服务。
func NewCapabilityService() *CapabilityService {
	return &CapabilityService{
		prices:   repository.NewPrintPriceConfigRepository(),
		printers: repository.NewPrintPrinterRepository(),
		devices:  repository.NewPrintDeviceRepository(),
	}
}

// Sellable 返回门店当前可下单的参数组合（价格与打印机的交集）。
// 参数 shopID 为门店 ID；返回空集合表示该门店没有任何可下单组合。
func (s *CapabilityService) Sellable(ctx context.Context, shopID types.ID) ([]SellableCombo, error) {
	prices, err := s.prices.ListByShop(ctx, shopID)
	if err != nil {
		return nil, err
	}
	printers, err := s.printers.ListEnabledByShop(ctx, shopID)
	if err != nil {
		return nil, err
	}

	combos := make([]SellableCombo, 0, len(prices))
	for _, price := range prices {
		if price.Status != printEnums.StatusEnabled {
			continue
		}

		matched := matchPrinters(printers, price.PaperType, price.ColorType, price.PrintMode)
		if len(matched) == 0 {
			continue
		}

		printerIDs := make([]types.ID, 0, len(matched))
		for _, printer := range matched {
			printerIDs = append(printerIDs, printer.ID)
		}

		combos = append(combos, SellableCombo{
			PaperType:   price.PaperType,
			ColorType:   price.ColorType,
			PrintMode:   price.PrintMode,
			Price:       price.Price,
			PrinterIDs:  printerIDs,
			PrinterName: matched[0].Name,
		})
	}

	return combos, nil
}

// Warnings 做配置期体检：找出「有价格但没有打印机支持」的组合。
// 参数 shopID 为门店 ID；只返回提示，不阻断保存。
func (s *CapabilityService) Warnings(ctx context.Context, shopID types.ID) ([]CapabilityWarning, error) {
	prices, err := s.prices.ListByShop(ctx, shopID)
	if err != nil {
		return nil, err
	}
	printers, err := s.printers.ListEnabledByShop(ctx, shopID)
	if err != nil {
		return nil, err
	}

	warnings := make([]CapabilityWarning, 0)
	for _, price := range prices {
		if price.Status != printEnums.StatusEnabled {
			continue
		}
		if len(matchPrinters(printers, price.PaperType, price.ColorType, price.PrintMode)) > 0 {
			continue
		}

		warnings = append(warnings, CapabilityWarning{
			PaperType: price.PaperType,
			ColorType: price.ColorType,
			PrintMode: price.PrintMode,
			Message:   "当前没有支持该组合的启用打印机，该组合暂不可下单",
		})
	}

	return warnings, nil
}

// ResolvePrinter 为一次打印选定打印机，严格匹配、不兜底、不降级。
// 参数 shopID 为门店 ID，其余为打印参数；无匹配返回业务错误。
func (s *CapabilityService) ResolvePrinter(ctx context.Context, shopID types.ID, paperType, colorType, printMode string) (*model.PrintPrinterModel, error) {
	printers, err := s.printers.ListEnabledByShop(ctx, shopID)
	if err != nil {
		return nil, err
	}

	matched := matchPrinters(printers, paperType, colorType, printMode)
	if len(matched) == 0 {
		return nil, baseErrs.New(printErrs.PrinterNoMatch, describeCombo(paperType, colorType, printMode))
	}

	return &matched[0], nil
}

// Price 取指定组合的单价；不存在或停用返回业务错误。
func (s *CapabilityService) Price(ctx context.Context, shopID types.ID, paperType, colorType, printMode string) (int64, error) {
	rule, err := s.prices.FindByRule(ctx, shopID, paperType, colorType, printMode)
	if err != nil {
		return 0, err
	}
	if rule == nil || rule.Status != printEnums.StatusEnabled {
		return 0, baseErrs.New(printErrs.PrinterNoMatch, describeCombo(paperType, colorType, printMode))
	}

	return rule.Price, nil
}

// matchPrinters 按三个参数筛选打印机：纸张在支持列表中，彩色/双面要求对应能力为 1。
// 参数 printers 为候选集合；返回满足条件的打印机（保持原有顺序）。
func matchPrinters(printers []model.PrintPrinterModel, paperType, colorType, printMode string) []model.PrintPrinterModel {
	matched := make([]model.PrintPrinterModel, 0, len(printers))

	for _, printer := range printers {
		if printer.Status != printEnums.StatusEnabled {
			continue
		}
		if !paperSupported(printer.PaperTypes, paperType) {
			continue
		}
		if colorType == printEnums.ColorColor && printer.ColorSupport != 1 {
			continue
		}
		if printMode == printEnums.ModeDouble && printer.DuplexSupport != 1 {
			continue
		}

		matched = append(matched, printer)
	}

	return matched
}

// paperSupported 判断逗号分隔的纸张能力是否包含目标纸张。
func paperSupported(paperTypes, target string) bool {
	if target == "" {
		return true
	}

	for _, item := range strings.Split(paperTypes, ",") {
		if strings.TrimSpace(item) == target {
			return true
		}
	}

	return false
}

// describeCombo 生成参数组合的中文描述，供错误提示直接展示。
func describeCombo(paperType, colorType, printMode string) string {
	paper := map[string]string{printEnums.PaperA4: "A4", printEnums.PaperA3: "A3"}[paperType]
	if paper == "" {
		paper = paperType
	}

	color := "黑白"
	if colorType == printEnums.ColorColor {
		color = "彩色"
	}

	mode := "单面"
	if printMode == printEnums.ModeDouble {
		mode = "双面"
	}

	return paper + " " + color + " " + mode
}
