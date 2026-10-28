package request

import (
	"github.com/goravel/framework/contracts/http"
)

// WalletAdjustRequest 是后台人工调账入参。
// amount 为正是加款、为负是扣款，单位为分。
type WalletAdjustRequest struct {
	Amount int64  `form:"amount" json:"amount"`
	Reason string `form:"reason" json:"reason"`
}

// Rules 返回校验规则；0 金额由 service 拒绝。
func (r *WalletAdjustRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"amount": "required|integer",
		"reason": "required|max:200",
	}
}

// Messages 覆盖默认提示。
func (r *WalletAdjustRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"amount.required": "调账金额不能为空",
		"amount.integer":  "调账金额必须是整数分",
		"reason.required": "请填写调账原因",
		"reason.max":      "调账原因不能超过 200 个字符",
	}
}
