// Package errs 定义 pay 模块的业务错误码，编码规则见 docs 8.3(2)。
package errs

import baseErrs "eqray-eprint/app/modules/base/errs"

// pay 模块错误码：1-005-{子模块}-{序号}。
// channel 占用 001，wallet 占用 002，payment 占用 003，refund 占用 004。
var (
	PayChannelConfigNotExists  = baseErrs.ErrorCode{Code: 1_005_001_001, Msg: "支付渠道配置【%v】不存在"}
	PayChannelConfigDisabled   = baseErrs.ErrorCode{Code: 1_005_001_002, Msg: "支付渠道【%v】已停用，请联系管理员"}
	PayChannelConfigInvalid    = baseErrs.ErrorCode{Code: 1_005_001_003, Msg: "支付渠道配置不完整：%v"}
	PayChannelConfigCodeExists = baseErrs.ErrorCode{Code: 1_005_001_004, Msg: "支付渠道编码【%v】已存在"}
	PayChannelSceneExists      = baseErrs.ErrorCode{Code: 1_005_001_005, Msg: "该渠道下【%v】场景已配置"}

	PayWalletNotExists        = baseErrs.ErrorCode{Code: 1_005_002_001, Msg: "会员【%v】还没有钱包账户"}
	PayWalletDisabled         = baseErrs.ErrorCode{Code: 1_005_002_002, Msg: "钱包账户已被冻结，暂不可支付"}
	PayWalletBalanceNotEnough = baseErrs.ErrorCode{Code: 1_005_002_003, Msg: "钱包余额不足"}
	PayWalletAmountInvalid    = baseErrs.ErrorCode{Code: 1_005_002_004, Msg: "金额必须大于 0"}

	PayPaymentNotExists      = baseErrs.ErrorCode{Code: 1_005_003_001, Msg: "支付单【%v】不存在"}
	PayPaymentStatusInvalid  = baseErrs.ErrorCode{Code: 1_005_003_002, Msg: "支付单【%v】当前状态不允许该操作"}
	PayPaymentAmountInvalid  = baseErrs.ErrorCode{Code: 1_005_003_003, Msg: "支付金额不合法：%v"}
	PayPaymentAmountMismatch = baseErrs.ErrorCode{Code: 1_005_003_004, Msg: "回调金额与支付单金额不一致"}
	PayGatewayNotSupported   = baseErrs.ErrorCode{Code: 1_005_003_005, Msg: "不支持的支付方式【%v】"}
	PayGatewayRequestFailed  = baseErrs.ErrorCode{Code: 1_005_003_006, Msg: "支付渠道请求失败：%v"}

	PayRefundNotExists     = baseErrs.ErrorCode{Code: 1_005_004_001, Msg: "退款单【%v】不存在"}
	PayRefundAmountExceed  = baseErrs.ErrorCode{Code: 1_005_004_002, Msg: "退款金额超过可退金额"}
	PayRefundStatusInvalid = baseErrs.ErrorCode{Code: 1_005_004_003, Msg: "退款单【%v】当前状态不允许该操作"}
)
