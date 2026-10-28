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

// PrinterInput 是打印机新增 / 编辑入参，能力字段由门店后台勾选。
type PrinterInput struct {
	PrintShopID   types.ID
	DeviceID      types.ID
	Name          string
	SystemName    string
	Model         string
	PaperTypes    string
	ColorSupport  int8
	DuplexSupport int8
	Status        int8
	Remark        string
}

// PrinterListQuery 是打印机列表查询条件。
type PrinterListQuery struct {
	ShopID   types.ID
	DeviceID types.ID
	Status   *int8
	PageNo   int
	PageSize int
}

// PrinterService 负责打印机登记与能力维护。
type PrinterService struct {
	printers *repository.PrintPrinterRepository
	devices  *repository.PrintDeviceRepository
}

// NewPrinterService 构造打印机服务。
func NewPrinterService() *PrinterService {
	return &PrinterService{
		printers: repository.NewPrintPrinterRepository(),
		devices:  repository.NewPrintDeviceRepository(),
	}
}

// Page 分页查询打印机。
func (s *PrinterService) Page(ctx context.Context, q PrinterListQuery) ([]model.PrintPrinterModel, int64, error) {
	return s.printers.Paginate(ctx, repository.PrinterFilter{
		ShopID:   q.ShopID,
		DeviceID: q.DeviceID,
		Status:   q.Status,
	}, q.PageNo, q.PageSize)
}

// Detail 查询打印机详情。
func (s *PrinterService) Detail(ctx context.Context, id types.ID) (*model.PrintPrinterModel, error) {
	printer, err := s.printers.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if printer == nil {
		return nil, baseErrs.New(printErrs.PrinterNotExists, id)
	}

	return printer, nil
}

// Create 登记打印机；同一设备下 system_name 唯一（队列名是静默打印的唯一依据）。
func (s *PrinterService) Create(ctx context.Context, in PrinterInput) (types.ID, error) {
	normalized, err := s.validate(ctx, in)
	if err != nil {
		return 0, err
	}

	printer := &model.PrintPrinterModel{
		PrintShopID:   normalized.PrintShopID,
		DeviceID:      normalized.DeviceID,
		Name:          normalized.Name,
		SystemName:    normalized.SystemName,
		Model:         normalized.Model,
		PaperTypes:    normalized.PaperTypes,
		ColorSupport:  normalized.ColorSupport,
		DuplexSupport: normalized.DuplexSupport,
		Status:        normalized.Status,
		Remark:        normalized.Remark,
	}
	printer.ID = support.NewID()

	if err := s.printers.Create(ctx, printer); err != nil {
		if support.IsDuplicate(err) {
			return 0, baseErrs.New(printErrs.PrinterNameExists, normalized.SystemName)
		}

		return 0, err
	}

	return printer.ID, nil
}

// Update 编辑打印机配置与能力；改能力会影响已下单但未打印的任务，需谨慎。
func (s *PrinterService) Update(ctx context.Context, id types.ID, in PrinterInput) error {
	printer, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	normalized, err := s.validate(ctx, in)
	if err != nil {
		return err
	}

	printer.PrintShopID = normalized.PrintShopID
	printer.DeviceID = normalized.DeviceID
	printer.Name = normalized.Name
	printer.SystemName = normalized.SystemName
	printer.Model = normalized.Model
	printer.PaperTypes = normalized.PaperTypes
	printer.ColorSupport = normalized.ColorSupport
	printer.DuplexSupport = normalized.DuplexSupport
	printer.Status = normalized.Status
	printer.Remark = normalized.Remark

	if err := s.printers.Update(ctx, printer); err != nil {
		if support.IsDuplicate(err) {
			return baseErrs.New(printErrs.PrinterNameExists, normalized.SystemName)
		}

		return err
	}

	return nil
}

// Disable 停用打印机（不物理删除，避免队列名被历史记录占用）。
func (s *PrinterService) Disable(ctx context.Context, id types.ID, updater string) error {
	if _, err := s.Detail(ctx, id); err != nil {
		return err
	}

	return s.printers.Disable(ctx, id, updater)
}

// ListEnabled 查询门店启用的打印机，供能力服务比对。
func (s *PrinterService) ListEnabled(ctx context.Context, shopID types.ID) ([]model.PrintPrinterModel, error) {
	return s.printers.ListEnabledByShop(ctx, shopID)
}

// validate 归一化并校验入参：设备必须存在且属于同一门店，纸张能力必须是 A4/A3。
func (s *PrinterService) validate(ctx context.Context, in PrinterInput) (PrinterInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.SystemName = strings.TrimSpace(in.SystemName)
	in.Model = strings.TrimSpace(in.Model)
	in.Remark = strings.TrimSpace(in.Remark)
	in.PaperTypes = normalizePaperTypes(in.PaperTypes)

	if in.Name == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印机名称不能为空")
	}
	if in.SystemName == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "Windows 打印队列名称不能为空")
	}
	if in.PaperTypes == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "请至少选择一种纸张")
	}
	if in.ColorSupport != 0 && in.ColorSupport != 1 {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "彩色能力取值不合法")
	}
	if in.DuplexSupport != 0 && in.DuplexSupport != 1 {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "双面能力取值不合法")
	}
	if in.Status != printEnums.StatusEnabled && in.Status != printEnums.StatusDisabled {
		in.Status = printEnums.StatusEnabled
	}

	device, err := s.devices.FindByID(ctx, in.DeviceID)
	if err != nil {
		return in, err
	}
	if device == nil {
		return in, baseErrs.New(printErrs.DeviceNotExists, in.DeviceID)
	}
	// 打印机必须挂在同门店的设备上，避免跨店路由
	in.PrintShopID = device.PrintShopID

	return in, nil
}

// normalizePaperTypes 校验并归一化纸张能力，只允许 A4 / A3，去重后用逗号拼接。
func normalizePaperTypes(value string) string {
	allowed := map[string]bool{printEnums.PaperA4: true, printEnums.PaperA3: true}
	seen := make(map[string]bool)
	result := make([]string, 0, 2)

	for _, item := range strings.Split(value, ",") {
		paper := strings.ToUpper(strings.TrimSpace(item))
		if !allowed[paper] || seen[paper] {
			continue
		}
		seen[paper] = true
		result = append(result, paper)
	}

	return strings.Join(result, ",")
}
