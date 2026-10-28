package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/pay/service"
)

// WalletController 提供会员钱包的查询接口。
type WalletController struct {
	service *service.WalletService
}

// NewWalletController 构造钱包控制器。
func NewWalletController() *WalletController {
	return &WalletController{service: service.NewWalletService()}
}

// Show 查看我的钱包：余额与累计收支。
// 无入参，会员身份取自令牌；没有钱包时自动创建空账户。
func (c *WalletController) Show(ctx http.Context) http.Response {
	wallet, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, wallet)
}

// Transactions 分页查询我的钱包流水。
// 入参：pageNo、pageSize；返回流水列表与总数。
func (c *WalletController) Transactions(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	list, total, err := c.service.Transactions(ctx.Context(), support.MemberID(ctx.Context()), page.PageNo, page.PageSize)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}
