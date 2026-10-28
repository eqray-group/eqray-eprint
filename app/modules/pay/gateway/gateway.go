// Package gateway 定义支付网关抽象，屏蔽微信 / 支付宝 / 余额的实现差异。
//
// 每个渠道配置（pay_channel_config.code）对应一个 Gateway 实例，
// 上层只按 channelCode 取实例，不出现渠道分支判断；
// 新增支付宝或下线微信时，只增删实现与配置，不动支付流程。
package gateway

import (
	"context"
	"net/http"
	"time"
)

// Config 是解密后的渠道商户配置，由 ChannelConfigService 组装后传入。
// 密钥字段为明文，只在内存中使用，不落日志、不返回给前端。
type Config struct {
	Code              string
	Channel           string
	Scene             string
	AppID             string
	MchID             string
	SerialNo          string
	PrivateKey        string
	APIKey            string
	PlatformPublicKey string
	GatewayURL        string
	NotifyURL         string
	ReturnURL         string
}

// Order 是下单必要信息，金额单位为分。
type Order struct {
	PaymentNo   string
	Description string
	Amount      int64
	// OpenID 仅 JSAPI（小程序 / 公众号）需要，且必须属于当前配置的 AppID
	OpenID string
	// ClientIP 仅 H5 需要，必须是用户真实公网 IP
	ClientIP string
	ExpireAt time.Time
}

// CreateResult 是下单结果：JSAPI 返回调起参数，H5 返回跳转地址。
type CreateResult struct {
	Scene   string
	Params  map[string]any
	MwebURL string
}

// QueryResult 是查单结果，Success 表示渠道侧已支付成功。
type QueryResult struct {
	Paid        bool
	TradeNo     string
	BuyerOpenID string
	PaidAmount  int64
}

// RefundInput 是退款入参，Total 为原支付单金额，Amount 为本次退款金额。
type RefundInput struct {
	RefundNo  string
	PaymentNo string
	Reason    string
	Amount    int64
	Total     int64
}

// RefundResult 是退款受理结果：Processing 为 true 表示渠道异步处理，需等回调。
type RefundResult struct {
	ChannelRefundNo string
	Success         bool
	Processing      bool
}

// NotifyResult 是回调解析结果，PaidAmount 用于与支付单金额比对。
type NotifyResult struct {
	PaymentNo   string
	TradeNo     string
	BuyerOpenID string
	PaidAmount  int64
	Success     bool
}

// RefundNotifyResult 是退款回调解析结果。
type RefundNotifyResult struct {
	PaymentNo       string
	RefundNo        string
	ChannelRefundNo string
	Success         bool
}

// Gateway 是一个支付场景的实现能力集合。
type Gateway interface {
	// Code 返回渠道配置编码，与 pay_channel_config.code 一致。
	Code() string
	// Create 向渠道下单，返回前端调起支付所需的数据。
	Create(ctx context.Context, order Order, cfg Config) (*CreateResult, error)
	// Query 主动查单，用于回调丢失时的兜底确认。
	Query(ctx context.Context, paymentNo string, cfg Config) (*QueryResult, error)
	// Refund 发起退款，余额渠道直接成功，微信渠道返回处理中。
	Refund(ctx context.Context, in RefundInput, cfg Config) (*RefundResult, error)
	// ParseNotify 验签并解密渠道支付回调，验签失败必须返回错误。
	ParseNotify(ctx context.Context, header http.Header, body []byte, cfg Config) (*NotifyResult, error)
	// ParseRefundNotify 验签并解密渠道退款回调。
	ParseRefundNotify(ctx context.Context, header http.Header, body []byte, cfg Config) (*RefundNotifyResult, error)
	// Close 关闭渠道侧订单，关单定时任务使用；已支付或已关闭时渠道会报错，可忽略。
	Close(ctx context.Context, paymentNo string, cfg Config) error
}
