package request

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/types"
)

// MemberCardSaveRequest 是后台发放 / 修改会员卡的入参。
// 一个会员只有一张卡，因此 memberUserId 是定位键，卡号由服务端生成。
type MemberCardSaveRequest struct {
	MemberUserID    types.ID        `form:"memberUserId" json:"memberUserId"`
	CardLevel       string          `form:"cardLevel" json:"cardLevel"`
	CardName        string          `form:"cardName" json:"cardName"`
	DiscountPercent int             `form:"discountPercent" json:"discountPercent"`
	Status          int8            `form:"status" json:"status"`
	StartTime       *types.DateTime `form:"startTime" json:"startTime"`
	ExpireTime      *types.DateTime `form:"expireTime" json:"expireTime"`
	Remark          string          `form:"remark" json:"remark"`
}

// Rules 返回校验规则。
func (r *MemberCardSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"memberUserId":    "required|min:1",
		"cardLevel":       "required|max:32",
		"cardName":        "required|max:64",
		"discountPercent": "required|integer|between:1,100",
		"status":          "integer|in:0,1,2",
		"remark":          "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *MemberCardSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"memberUserId.required":    "请选择会员",
		"memberUserId.min":         "请选择会员",
		"cardLevel.required":       "卡等级编码不能为空",
		"cardLevel.max":            "卡等级编码不能超过 32 个字符",
		"cardName.required":        "卡名称不能为空",
		"cardName.max":             "卡名称不能超过 64 个字符",
		"discountPercent.required": "折扣不能为空",
		"discountPercent.integer":  "折扣必须是整数",
		"discountPercent.between":  "折扣取值范围为 1-100，100 表示无折扣",
		"status.integer":           "会员卡状态取值不合法",
		"status.in":                "会员卡状态取值不合法",
		"remark.max":               "备注不能超过 255 个字符",
	}
}
