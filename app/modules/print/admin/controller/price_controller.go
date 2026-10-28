package controller

import (
	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/admin/request"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/service"
)

// PriceController 提供门店打印价格的后台接口。
type PriceController struct {
	service *service.PriceService
}

// NewPriceController 构造价格控制器。
func NewPriceController() *PriceController {
	return &PriceController{service: service.NewPriceService()}
}

// Index 分页查询价格规则。
func (c *PriceController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var shopID types.ID
	if raw := ctx.Request().QueryInt("printShopId", 0); raw > 0 {
		shopID = types.ID(raw)
	}

	list, total, err := c.service.Page(ctx.Context(), service.PriceListQuery{
		ShopID:   shopID,
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Store 新增价格规则，同一门店同一组合唯一。
func (c *PriceController) Store(ctx http.Context) http.Response {
	var req request.PriceSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	status := printEnums.StatusEnabled
	if req.Status != nil {
		status = *req.Status
	}

	id, err := c.service.Create(ctx.Context(), service.PriceInput{
		PrintShopID: types.ID(req.PrintShopID),
		PaperType:   req.PaperType,
		ColorType:   req.ColorType,
		PrintMode:   req.PrintMode,
		Price:       req.Price,
		Status:      status,
	})
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 修改价格与启用状态（组合本身不可改）。
func (c *PriceController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.PriceUpdateRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, req.Price, req.Status); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Destroy 删除价格规则。
func (c *PriceController) Destroy(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Delete(ctx.Context(), id); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Warnings 返回该门店「有价格但无打印机支持」的组合，供配置页提示。
func (c *PriceController) Warnings(ctx http.Context) http.Response {
	raw := ctx.Request().QueryInt("printShopId", 0)
	if raw <= 0 {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "请选择门店"))
	}

	warnings, err := c.service.Warnings(ctx.Context(), types.ID(raw))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, warnings)
}
