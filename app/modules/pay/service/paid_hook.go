package service

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
)

// paidHook 是支付成功后的业务回调，由 bootstrap 注入（如打印订单派发任务）。
// pay 模块不 import 任何业务模块，避免反向依赖。
var paidHook func(ctx context.Context, bizType string, bizID types.ID) error

// SetPaidHook 注册支付成功回调；传 nil 表示取消注册。
// 回调在支付单已置成功的独立事务之后执行，失败不影响支付结果，只记录日志。
func SetPaidHook(hook func(ctx context.Context, bizType string, bizID types.ID) error) {
	paidHook = hook
}

// notifyPaid 触发支付成功回调；没有注册或回调报错时只记录日志。
func notifyPaid(ctx context.Context, bizType string, bizID types.ID) {
	if paidHook == nil {
		return
	}

	if err := paidHook(ctx, bizType, bizID); err != nil {
		facades.Log().Errorf("支付成功回调处理失败 bizType=%s bizID=%d err=%v", bizType, bizID, err)
	}
}

// refundedHook 是退款成功后的业务回调，由 bootstrap 注入（如打印订单置已退款）。
var refundedHook func(ctx context.Context, bizType string, bizID types.ID) error

// SetRefundedHook 注册退款成功回调；传 nil 表示取消注册。
// 余额退款同步完成、微信退款走回调确认，两条路径都会触发。
func SetRefundedHook(hook func(ctx context.Context, bizType string, bizID types.ID) error) {
	refundedHook = hook
}

// notifyRefunded 触发退款成功回调；没有注册或回调报错时只记录日志。
func notifyRefunded(ctx context.Context, bizType string, bizID types.ID) {
	if refundedHook == nil {
		return
	}

	if err := refundedHook(ctx, bizType, bizID); err != nil {
		facades.Log().Errorf("退款成功回调处理失败 bizType=%s bizID=%d err=%v", bizType, bizID, err)
	}
}
