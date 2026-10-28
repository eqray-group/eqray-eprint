// Package enums 定义支付与钱包的枚举常量，数据库存字符串或 tinyint，代码里只用常量。
package enums

// 支付渠道大类，对应 pay_channel_config.channel 与 pay_payment.channel。
const (
	ChannelWechat  = "wechat"
	ChannelAlipay  = "alipay"
	ChannelBalance = "balance"
)

// 支付场景，对应 pay_channel_config.scene 与 pay_payment.scene。
// mini 小程序、jsapi 公众号（微信内置浏览器）、h5 手机浏览器、wap 支付宝手机网站。
const (
	SceneMini  = "mini"
	SceneJsapi = "jsapi"
	SceneH5    = "h5"
	SceneWap   = "wap"
	SceneQr    = "qr"
)

// 渠道配置编码，对应 pay_channel_config.code，也是 Gateway 的注册键。
const (
	CodeWechatMini = "wechat_mini"
	CodeWechatMp   = "wechat_mp"
	CodeWechatH5   = "wechat_h5"
	CodeAlipayWap  = "alipay_wap"
)

// 业务类型，对应 pay_payment.biz_type 与 pay_wallet_transaction.biz_type。
const (
	BizTypePrintOrder     = "print_order"
	BizTypeWalletRecharge = "wallet_recharge"
	BizTypePrintConsume   = "print_consume"
	BizTypePrintRefund    = "print_refund"
	BizTypeManualAdjust   = "manual_adjust"
)

// 支付单状态。
const (
	PaymentStatusPending  int8 = 0 // 待支付
	PaymentStatusPaying   int8 = 1 // 支付中（已调起渠道，等回调）
	PaymentStatusPaid     int8 = 2 // 已支付
	PaymentStatusFailed   int8 = 3 // 支付失败
	PaymentStatusClosed   int8 = 4 // 已关闭（超时或取消）
	PaymentStatusRefunded int8 = 5 // 已全额退款
)

// 退款单状态。
const (
	RefundStatusPending   int8 = 0 // 待退款
	RefundStatusRefunding int8 = 1 // 退款中
	RefundStatusSuccess   int8 = 2 // 退款成功
	RefundStatusFailed    int8 = 3 // 退款失败
)

// 钱包流水方向。
const (
	WalletTxIncome  int8 = 1 // 收入
	WalletTxExpense int8 = 2 // 支出
)

// 钱包状态。
const (
	WalletStatusFrozen int8 = 0 // 冻结，禁止支付
	WalletStatusNormal int8 = 1 // 正常
)

// 渠道配置状态。
const (
	ChannelConfigDisabled int8 = 0
	ChannelConfigEnabled  int8 = 1
)

// 回调日志处理结果。
const (
	NotifyHandlePending int8 = 0 // 待处理
	NotifyHandleSuccess int8 = 1 // 处理成功
	NotifyHandleIgnored int8 = 2 // 幂等忽略
	NotifyHandleFailed  int8 = 3 // 处理失败
)
