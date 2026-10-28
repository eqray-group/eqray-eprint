package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/admin/request"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/service"
)

// PrinterController 提供打印机与能力配置的后台接口。
type PrinterController struct {
	service *service.PrinterService
}

// NewPrinterController 构造打印机控制器。
func NewPrinterController() *PrinterController {
	return &PrinterController{service: service.NewPrinterService()}
}

// Index 分页查询打印机。
func (c *PrinterController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.PrinterListQuery{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if shopID := ctx.Request().QueryInt("printShopId", 0); shopID > 0 {
		q.ShopID = types.ID(shopID)
	}
	if deviceID := ctx.Request().QueryInt("deviceId", 0); deviceID > 0 {
		q.DeviceID = types.ID(deviceID)
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := int8(status)
		q.Status = &v
	}

	list, total, err := c.service.Page(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询打印机详情。
func (c *PrinterController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	printer, err := c.service.Detail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, printer)
}

// Store 登记打印机并配置能力。
func (c *PrinterController) Store(ctx http.Context) http.Response {
	var req request.PrinterSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	id, err := c.service.Create(ctx.Context(), toPrinterInput(req))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 编辑打印机与能力；能力变更会影响后续任务路由。
func (c *PrinterController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.PrinterSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, toPrinterInput(req)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Disable 停用打印机（不物理删除）。
func (c *PrinterController) Disable(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Disable(ctx.Context(), id, operatorName(ctx)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// toPrinterInput 把请求体转成服务入参，状态缺省按启用处理。
func toPrinterInput(req request.PrinterSaveRequest) service.PrinterInput {
	status := printEnums.StatusEnabled
	if req.Status != nil {
		status = *req.Status
	}

	return service.PrinterInput{
		DeviceID:      types.ID(req.DeviceID),
		Name:          req.Name,
		SystemName:    req.SystemName,
		Model:         req.Model,
		PaperTypes:    req.PaperTypes,
		ColorSupport:  req.ColorSupport,
		DuplexSupport: req.DuplexSupport,
		Status:        status,
		Remark:        req.Remark,
	}
}
