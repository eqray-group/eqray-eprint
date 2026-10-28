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

// PrinterReportItem 是客户端上报的一台本机打印机及其能力配置。
// SystemName 为 Windows 打印队列名，是同一设备内的唯一标识。
type PrinterReportItem struct {
	SystemName    string
	Name          string
	Model         string
	PaperTypes    string
	ColorSupport  int8
	DuplexSupport int8
}

// maxReportPrinters 限制单次上报的打印机数量，避免客户端误传超长数组。
const maxReportPrinters = 50

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

// ListByDevice 返回设备下已上报的打印机（含停用），客户端用来回填勾选状态与能力。
func (s *PrinterService) ListByDevice(ctx context.Context, deviceID types.ID) ([]model.PrintPrinterModel, error) {
	return s.printers.ListByDevice(ctx, deviceID)
}

// SyncDevicePrinters 用客户端勾选结果覆盖设备的打印机集合。
// 参数 device 为设备记录，items 为本次勾选的完整集合：已存在的按队列名更新并启用，
// 新增的直接创建，本次未上报的停用（不物理删除，保留历史任务引用）。
// 返回同步后的完整列表。
func (s *PrinterService) SyncDevicePrinters(ctx context.Context, device *model.PrintDeviceModel, items []PrinterReportItem) ([]model.PrintPrinterModel, error) {
	if device == nil {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "设备不能为空")
	}
	if len(items) > maxReportPrinters {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "单次最多上报 50 台打印机")
	}

	existing, err := s.printers.ListByDevice(ctx, device.ID)
	if err != nil {
		return nil, err
	}

	// Windows 队列名不区分大小写，这里统一用小写键匹配，避免仅大小写变化时重复创建
	bySystem := make(map[string]*model.PrintPrinterModel, len(existing))
	for i := range existing {
		bySystem[strings.ToLower(existing[i].SystemName)] = &existing[i]
	}

	seen := make(map[string]bool, len(items))
	for _, raw := range items {
		item, err := normalizeReportItem(raw)
		if err != nil {
			return nil, err
		}
		key := strings.ToLower(item.SystemName)
		if seen[key] {
			return nil, baseErrs.WithMsg(baseErrs.BadRequest, "打印机【"+item.SystemName+"】重复上报")
		}
		seen[key] = true

		if printer, ok := bySystem[key]; ok {
			printer.PrintShopID = device.PrintShopID
			printer.Name = item.Name
			printer.Model = item.Model
			printer.PaperTypes = item.PaperTypes
			printer.ColorSupport = item.ColorSupport
			printer.DuplexSupport = item.DuplexSupport
			printer.Status = printEnums.StatusEnabled

			if err := s.printers.Update(ctx, printer); err != nil {
				return nil, err
			}

			continue
		}

		printer := &model.PrintPrinterModel{
			PrintShopID:   device.PrintShopID,
			DeviceID:      device.ID,
			Name:          item.Name,
			SystemName:    item.SystemName,
			Model:         item.Model,
			PaperTypes:    item.PaperTypes,
			ColorSupport:  item.ColorSupport,
			DuplexSupport: item.DuplexSupport,
			Status:        printEnums.StatusEnabled,
		}
		printer.ID = support.NewID()

		if err := s.printers.Create(ctx, printer); err != nil {
			if support.IsDuplicate(err) {
				return nil, baseErrs.New(printErrs.PrinterNameExists, item.SystemName)
			}

			return nil, err
		}
	}

	// 本次未勾选的打印机统一停用：队列已不存在或用户主动取消上报
	for i := range existing {
		if seen[strings.ToLower(existing[i].SystemName)] || existing[i].Status != printEnums.StatusEnabled {
			continue
		}
		if err := s.printers.Disable(ctx, existing[i].ID, "打印客户端"); err != nil {
			return nil, err
		}
	}

	return s.printers.ListByDevice(ctx, device.ID)
}

// normalizeReportItem 清洗客户端上报的打印机字段并校验能力取值。
// 返回归一化后的条目；队列名为空或纸张能力非法时返回 400 业务错误。
func normalizeReportItem(in PrinterReportItem) (PrinterReportItem, error) {
	in.SystemName = strings.TrimSpace(in.SystemName)
	in.Name = strings.TrimSpace(in.Name)
	in.Model = strings.TrimSpace(in.Model)

	if in.SystemName == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印队列名称不能为空")
	}
	if len(in.SystemName) > 128 {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印队列名称不能超过 128 个字符")
	}
	if in.Name == "" {
		in.Name = in.SystemName
	}
	if len(in.Name) > 64 || len(in.Model) > 64 {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印机名称与型号不能超过 64 个字符")
	}

	in.PaperTypes = normalizePaperTypes(in.PaperTypes)
	if in.PaperTypes == "" {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印机【"+in.SystemName+"】至少选择一种纸张")
	}
	if in.ColorSupport != 0 && in.ColorSupport != 1 {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印机【"+in.SystemName+"】彩色能力取值不合法")
	}
	if in.DuplexSupport != 0 && in.DuplexSupport != 1 {
		return in, baseErrs.WithMsg(baseErrs.BadRequest, "打印机【"+in.SystemName+"】双面能力取值不合法")
	}

	return in, nil
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
