package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/pay/admin/request"
	"eqray-eprint/app/modules/pay/service"
)

// WalletController 提供后台会员钱包的人工调账接口。
type WalletController struct {
	service *service.WalletService
}

// NewWalletController 构造后台钱包控制器。
func NewWalletController() *WalletController {
	return &WalletController{service: service.NewWalletService()}
}

// Show 查询会员钱包余额与累计收支。
// 入参：路径 memberId；没有钱包时返回 0 余额，不触发懒创建。
func (c *WalletController) Show(ctx http.Context) http.Response {
	memberID, err := param.PathValue(ctx, "memberId")
	if err != nil {
		return result.Fail(ctx, err)
	}

	wallet, err := c.service.AdminDetail(ctx.Context(), memberID)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, wallet)
}

// Index 分页查询钱包流水，可按会员、收支方向、业务类型筛选。
func (c *WalletController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.AdminTransactionQuery{
		BizType:  ctx.Request().Query("bizType"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if raw := ctx.Request().QueryInt("memberId", 0); raw > 0 {
		q.MemberID = types.ID(raw)
	}
	if direction := ctx.Request().QueryInt("type", -1); direction >= 0 {
		v := int8(direction)
		q.Type = &v
	}

	list, total, err := c.service.AdminTransactions(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Adjust 给会员钱包人工调账：正数加款、负数扣款，返回调账后余额。
// 入参：路径 memberId + WalletAdjustRequest；返回余额。
func (c *WalletController) Adjust(ctx http.Context) http.Response {
	memberID, err := param.PathValue(ctx, "memberId")
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.WalletAdjustRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	balance, err := c.service.Adjust(ctx.Context(), memberID, req.Amount, req.Reason)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"balance": balance})
}
