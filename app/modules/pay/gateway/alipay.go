package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"

	payEnums "eqray-eprint/app/modules/pay/enums"
)

// alipaySuccess 是支付宝网关返回的成功码。
const alipaySuccess = "10000"

// AlipayGateway 是支付宝手机网站支付实现（alipay_wap）。
// 使用公钥模式（应用私钥 + 支付宝公钥），不依赖上传证书；
// 支付宝退款为同步返回，因此不存在退款回调。
type AlipayGateway struct{}

// NewAlipayWapGateway 构造支付宝手机网站支付网关。
func NewAlipayWapGateway() *AlipayGateway {
	return &AlipayGateway{}
}

// Code 返回渠道配置编码。
func (g *AlipayGateway) Code() string {
	return payEnums.CodeAlipayWap
}

// Create 调起支付宝手机网站支付，返回收银台跳转地址。
// 金额从分转成元字符串，支付宝要求两位小数。
func (g *AlipayGateway) Create(ctx context.Context, order Order, cfg Config) (*CreateResult, error) {
	client, err := g.client(cfg)
	if err != nil {
		return nil, err
	}

	bm := gopay.BodyMap{
		"out_trade_no":    order.PaymentNo,
		"total_amount":    centsToYuan(order.Amount),
		"subject":         order.Description,
		"product_code":    "QUICK_WAP_WAY",
		"timeout_express": timeoutExpress(order.ExpireAt),
	}
	if cfg.NotifyURL != "" {
		bm.Set("notify_url", cfg.NotifyURL)
	}
	if cfg.ReturnURL != "" {
		// 用户支付完成后回跳的页面，仅用于体验，结果以回调/查单为准
		bm.Set("quit_url", cfg.ReturnURL)
	}

	payURL, err := client.TradeWapPay(ctx, bm)
	if err != nil {
		return nil, err
	}

	return &CreateResult{Scene: payEnums.SceneWap, MwebURL: payURL}, nil
}

// Query 主动查单，用于回调丢失时确认支付结果。
func (g *AlipayGateway) Query(ctx context.Context, paymentNo string, cfg Config) (*QueryResult, error) {
	client, err := g.client(cfg)
	if err != nil {
		return nil, err
	}

	rsp, err := client.TradeQuery(ctx, gopay.BodyMap{"out_trade_no": paymentNo})
	if err != nil {
		return nil, err
	}
	if rsp.Response == nil || rsp.Response.Code != alipaySuccess {
		return &QueryResult{Paid: false}, nil
	}

	status := rsp.Response.TradeStatus

	return &QueryResult{
		Paid:       status == "TRADE_SUCCESS" || status == "TRADE_FINISHED",
		TradeNo:    rsp.Response.TradeNo,
		PaidAmount: yuanToCents(rsp.Response.TotalAmount),
	}, nil
}

// Refund 发起退款；支付宝同步返回结果，FundChange=Y 表示资金已变动。
func (g *AlipayGateway) Refund(ctx context.Context, in RefundInput, cfg Config) (*RefundResult, error) {
	client, err := g.client(cfg)
	if err != nil {
		return nil, err
	}

	bm := gopay.BodyMap{
		"out_trade_no":   in.PaymentNo,
		"refund_amount":  centsToYuan(in.Amount),
		"out_request_no": in.RefundNo,
	}
	if in.Reason != "" {
		bm.Set("refund_reason", in.Reason)
	}

	rsp, err := client.TradeRefund(ctx, bm)
	if err != nil {
		return nil, err
	}
	if rsp.Response == nil || rsp.Response.Code != alipaySuccess {
		msg := "支付宝退款失败"
		if rsp.Response != nil {
			msg = rsp.Response.SubMsg
		}

		return &RefundResult{Success: false}, errors.New(msg)
	}

	// FundChange=Y 表示本次确实发生资金变动；部分场景可能不返回该字段，按成功处理
	return &RefundResult{
		ChannelRefundNo: in.RefundNo,
		Success:         rsp.Response.FundChange == "Y" || rsp.Response.FundChange == "",
		Processing:      false,
	}, nil
}

// ParseNotify 解析并验签支付宝异步通知。
// 支付宝回调为 form 表单（application/x-www-form-urlencoded），
// 验签前先校验 app_id，避免其它应用的合法回调被当成本平台订单。
func (g *AlipayGateway) ParseNotify(ctx context.Context, header http.Header, body []byte, cfg Config) (*NotifyResult, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	bm, err := alipay.ParseNotifyByURLValues(values)
	if err != nil {
		return nil, err
	}
	if bm.GetString("app_id") != cfg.AppID {
		return nil, errors.New("支付宝回调 app_id 与配置不一致")
	}

	ok, err := alipay.VerifySign(cfg.PlatformPublicKey, bm)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("支付宝回调验签失败")
	}

	status := bm.GetString("trade_status")

	return &NotifyResult{
		PaymentNo:   bm.GetString("out_trade_no"),
		TradeNo:     bm.GetString("trade_no"),
		BuyerOpenID: bm.GetString("buyer_open_id"),
		PaidAmount:  yuanToCents(bm.GetString("total_amount")),
		Success:     status == "TRADE_SUCCESS" || status == "TRADE_FINISHED",
	}, nil
}

// ParseRefundNotify 支付宝退款为同步返回，不存在退款回调。
func (g *AlipayGateway) ParseRefundNotify(ctx context.Context, header http.Header, body []byte, cfg Config) (*RefundNotifyResult, error) {
	return nil, errors.New("支付宝退款为同步返回，无退款回调")
}

// Close 关闭支付宝未支付订单；已支付订单关闭会返回业务错误，由调用方忽略。
func (g *AlipayGateway) Close(ctx context.Context, paymentNo string, cfg Config) error {
	client, err := g.client(cfg)
	if err != nil {
		return err
	}

	rsp, err := client.TradeClose(ctx, gopay.BodyMap{"out_trade_no": paymentNo})
	if err != nil {
		return err
	}
	if rsp.Response == nil || rsp.Response.Code != alipaySuccess {
		return fmt.Errorf("支付宝关单失败")
	}

	return nil
}

// client 用商户配置构造支付宝客户端；isProd 固定 true，沙箱不走线上配置。
func (g *AlipayGateway) client(cfg Config) (*alipay.Client, error) {
	if cfg.PrivateKey == "" {
		return nil, errors.New("支付宝应用私钥未配置")
	}

	client, err := alipay.NewClient(cfg.AppID, cfg.PrivateKey, true)
	if err != nil {
		return nil, err
	}
	client.SetCharset(alipay.UTF8)
	client.SetSignType(alipay.RSA2)

	return client, nil
}

// centsToYuan 把分转成支付宝需要的元字符串（两位小数）。
func centsToYuan(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}

// yuanToCents 把支付宝返回的元字符串转成分，解析失败按 0 处理（会被金额校验拦下）。
func yuanToCents(yuan string) int64 {
	if yuan == "" {
		return 0
	}

	value, err := strconv.ParseFloat(yuan, 64)
	if err != nil {
		return 0
	}

	return int64(value*100 + 0.5)
}

// timeoutExpress 生成支付宝的相对超时时长，最长 15 天。
func timeoutExpress(expireAt time.Time) string {
	if expireAt.IsZero() {
		return "30m"
	}

	remain := time.Until(expireAt)
	if remain <= 0 {
		return "1m"
	}
	if remain > 15*24*time.Hour {
		remain = 15 * 24 * time.Hour
	}

	return fmt.Sprintf("%dm", int(remain.Minutes()))
}
