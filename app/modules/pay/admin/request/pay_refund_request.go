package request

import (
	"github.com/goravel/framework/contracts/http"
)

// RefundCreateRequest 是后台发起退款的入参。
// 金额按分传入，服务端会校验不超过可退余额。
type RefundCreateRequest struct {
	PaymentNo string `form:"paymentNo" json:"paymentNo"`
	Amount    int64  `form:"amount" json:"amount"`
	Reason    string `form:"reason" json:"reason"`
}

// Rules 返回校验规则。
func (r *RefundCreateRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"paymentNo": "required|max:32",
		"amount":    "required|integer|min:1",
		"reason":    "required|max:255",
	}
}

// Messages 覆盖默认提示。
func (r *RefundCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"paymentNo.required": "支付单号不能为空",
		"amount.required":    "退款金额不能为空",
		"amount.min":         "退款金额必须大于 0",
		"amount.integer":     "退款金额必须是整数分",
		"reason.required":    "请填写退款原因",
		"reason.max":         "退款原因不能超过 255 个字符",
	}
}
