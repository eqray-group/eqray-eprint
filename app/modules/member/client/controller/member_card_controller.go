package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/member/service"
)

// MemberCardController 提供 C 端会员查看自己会员卡的接口。
type MemberCardController struct {
	service *service.MemberCardService
}

// NewMemberCardController 构造 C 端会员卡控制器。
func NewMemberCardController() *MemberCardController {
	return &MemberCardController{service: service.NewMemberCardService()}
}

// Show 查看「我的会员卡」：卡号、卡等级、折扣与有效期。
// 无入参，会员身份取自令牌；未发卡返回业务错误码。
func (c *MemberCardController) Show(ctx http.Context) http.Response {
	card, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, card)
}
