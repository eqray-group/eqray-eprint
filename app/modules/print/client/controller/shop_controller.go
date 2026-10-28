// Package controller 提供 print 模块的会员端接口。
package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/print/service"
)

// ShopController 提供会员端门店、可售组合与轮播图查询。
type ShopController struct {
	shops      *service.ShopService
	capability *service.CapabilityService
	banners    *service.BannerService
}

// NewShopController 构造会员端门店控制器。
func NewShopController() *ShopController {
	return &ShopController{
		shops:      service.NewShopService(),
		capability: service.NewCapabilityService(),
		banners:    service.NewBannerService(),
	}
}

// Index 返回可下单门店（审核通过且营业中）。
func (c *ShopController) Index(ctx http.Context) http.Response {
	list, err := c.shops.ListOpen(ctx.Context())
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, list)
}

// Sellable 返回门店当前可下单的参数组合与单价。
// 组合 = 启用价格 ∩ 启用打印机能力，前端据此渲染纸张/色彩/单双面选项。
func (c *ShopController) Sellable(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	combos, err := c.capability.Sellable(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, combos)
}

// Banners 返回门店可见轮播图：门店自有 + 全平台，仅启用，按 sort 升序。
func (c *ShopController) Banners(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	list, err := c.banners.ListForShop(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, list)
}
