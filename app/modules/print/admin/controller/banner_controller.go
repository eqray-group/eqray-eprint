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

// BannerController 提供门店轮播图的后台管理接口。
type BannerController struct {
	service *service.BannerService
}

// NewBannerController 构造轮播图控制器。
func NewBannerController() *BannerController {
	return &BannerController{service: service.NewBannerService()}
}

// Index 分页查询轮播图，printShopId 为 0 时查全平台。
func (c *BannerController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var shopID types.ID
	if raw := ctx.Request().QueryInt("printShopId", 0); raw > 0 {
		shopID = types.ID(raw)
	}

	list, total, err := c.service.Page(ctx.Context(), shopID, page.PageNo, page.PageSize)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询轮播图详情。
func (c *BannerController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	banner, err := c.service.Detail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, banner)
}

// Store 新增轮播图。
func (c *BannerController) Store(ctx http.Context) http.Response {
	var req request.BannerSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	id, err := c.service.Create(ctx.Context(), toBannerInput(req))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 编辑轮播图。
func (c *BannerController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.BannerSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, toBannerInput(req)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Destroy 删除轮播图。
func (c *BannerController) Destroy(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Delete(ctx.Context(), id, operatorName(ctx)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// toBannerInput 把请求体转成服务入参，状态缺省按启用处理。
func toBannerInput(req request.BannerSaveRequest) service.BannerInput {
	status := printEnums.StatusEnabled
	if req.Status != nil {
		status = *req.Status
	}

	return service.BannerInput{
		PrintShopID: types.ID(req.PrintShopID),
		Title:       req.Title,
		ImageURL:    req.ImageURL,
		LinkType:    req.LinkType,
		LinkURL:     req.LinkURL,
		Sort:        req.Sort,
		Status:      status,
	}
}
