package gateway

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"

	"eqray-eprint/app/modules/pay/enums"
)

// defaultExpire 是支付超时缺省值，与 pay_payment.expire_time 的默认 30 分钟一致。
const defaultExpire = 30 * time.Minute

// 微信退款单状态取值，gopay 未提供常量，这里本地定义避免魔法字符串散落。
const (
	wxRefundSuccess    = "SUCCESS"
	wxRefundProcessing = "PROCESSING"
)

// WechatGateway 是微信支付实现，一个小程序 / 公众号 / H5 场景一个实例。
// 三个场景共用同一个商户号与私钥，仅下单接口、返回结构和调起方式不同。
type WechatGateway struct {
	scene string
}

// NewWechatMiniGateway 构造小程序支付网关（JSAPI，需小程序 openid）。
func NewWechatMiniGateway() *WechatGateway {
	return &WechatGateway{scene: enums.SceneMini}
}

// NewWechatMpGateway 构造公众号支付网关（JSAPI，微信内置浏览器，需公众号 openid）。
func NewWechatMpGateway() *WechatGateway {
	return &WechatGateway{scene: enums.SceneJsapi}
}

// NewWechatH5Gateway 构造 H5 支付网关（MWEB，微信外手机浏览器，需 client_ip）。
func NewWechatH5Gateway() *WechatGateway {
	return &WechatGateway{scene: enums.SceneH5}
}

// Code 返回渠道配置编码，与 pay_channel_config.code 一致。
func (g *WechatGateway) Code() string {
	switch g.scene {
	case enums.SceneJsapi:
		return enums.CodeWechatMp
	case enums.SceneH5:
		return enums.CodeWechatH5
	default:
		return enums.CodeWechatMini
	}
}

// Create 向微信下单：H5 走 MWEB 返回跳转地址，小程序与公众号走 JSAPI 返回调起参数。
func (g *WechatGateway) Create(ctx context.Context, order Order, cfg Config) (*CreateResult, error) {
	client, err := g.client(cfg)
	if err != nil {
		return nil, err
	}

	bm := gopay.BodyMap{
		"appid":        cfg.AppID,
		"mchid":        cfg.MchID,
		"description":  order.Description,
		"out_trade_no": order.PaymentNo,
		"notify_url":   cfg.NotifyURL,
		"amount": gopay.BodyMap{
			"total":    order.Amount,
			"currency": "CNY",
		},
	}
	expireAt := order.ExpireAt
	if expireAt.IsZero() {
		expireAt = time.Now().Add(defaultExpire)
	}
	bm.Set("time_expire", expireAt.Format(time.RFC3339))

	if g.scene == enums.SceneH5 {
		// MWEB 必须带真实客户端 IP 与场景信息，否则微信直接拒单
		bm.Set("scene_info", gopay.BodyMap{
			"payer_client_ip": order.ClientIP,
			"h5_info": gopay.BodyMap{
				"type": "Wap",
			},
		})

		rsp, err := client.V3TransactionH5(ctx, bm)
		if err != nil {
			return nil, err
		}
		if rsp.Code != wechat.Success || rsp.Response == nil {
			return nil, fmt.Errorf("微信H5下单失败(%d)：%s", rsp.Code, rsp.Error)
		}

		return &CreateResult{Scene: g.scene, MwebURL: rsp.Response.H5Url}, nil
	}

	if order.OpenID == "" {
		return nil, errors.New("JSAPI 支付缺少 openid")
	}
	bm.Set("payer", gopay.BodyMap{"openid": order.OpenID})

	rsp, err := client.V3TransactionJsapi(ctx, bm)
	if err != nil {
		return nil, err
	}
	if rsp.Code != wechat.Success || rsp.Response == nil {
		return nil, fmt.Errorf("微信JSAPI下单失败(%d)：%s", rsp.Code, rsp.Error)
	}

	params, err := g.payParams(client, cfg.AppID, rsp.Response.PrepayId)
	if err != nil {
		return nil, err
	}

	return &CreateResult{Scene: g.scene, Params: params}, nil
}

// Query 主动查单，用于前端轮询与回调丢失兜底。
// 渠道返回非成功状态时按「未支付」处理，不当作错误。
func (g *WechatGateway) Query(ctx context.Context, paymentNo string, cfg Config) (*QueryResult, error) {
	client, err := g.client(cfg)
	if err != nil {
		return nil, err
	}

	rsp, err := client.V3TransactionQueryOrder(ctx, wechat.OutTradeNo, paymentNo)
	if err != nil {
		return nil, err
	}
	if rsp.Code != wechat.Success || rsp.Response == nil {
		return &QueryResult{Paid: false}, nil
	}

	res := &QueryResult{
		Paid:    rsp.Response.TradeState == wechat.TradeStateSuccess,
		TradeNo: rsp.Response.TransactionId,
	}
	if rsp.Response.Amount != nil {
		res.PaidAmount = int64(rsp.Response.Amount.Total)
	}
	if rsp.Response.Payer != nil {
		res.BuyerOpenID = rsp.Response.Payer.Openid
	}

	return res, nil
}

// Refund 发起退款；微信为异步受理，Success 表示受理成功而非到账成功。
func (g *WechatGateway) Refund(ctx context.Context, in RefundInput, cfg Config) (*RefundResult, error) {
	client, err := g.client(cfg)
	if err != nil {
		return nil, err
	}

	bm := gopay.BodyMap{
		"out_trade_no":  in.PaymentNo,
		"out_refund_no": in.RefundNo,
		"reason":        in.Reason,
		"amount": gopay.BodyMap{
			"refund":   in.Amount,
			"total":    in.Total,
			"currency": "CNY",
		},
	}
	if cfg.NotifyURL != "" {
		// 配置里存的是支付回调地址，退款回调走同前缀的 /refund 路由，避免两个事件打到同一入口
		bm.Set("notify_url", refundNotifyURL(cfg.NotifyURL))
	}

	rsp, err := client.V3Refund(ctx, bm)
	if err != nil {
		return nil, err
	}
	if rsp.Code != wechat.Success || rsp.Response == nil {
		return &RefundResult{Success: false}, fmt.Errorf("微信退款失败(%d)：%s", rsp.Code, rsp.Error)
	}

	return &RefundResult{
		ChannelRefundNo: rsp.Response.RefundId,
		Success:         rsp.Response.Status == wxRefundSuccess,
		Processing:      rsp.Response.Status == wxRefundProcessing,
	}, nil
}

// ParseNotify 验签并解密支付回调。
// 必须先验签再解密，验签失败不得落库为有效回调；
// 平台公钥未配置时直接报错，避免「静默放过未验签回调」。
func (g *WechatGateway) ParseNotify(ctx context.Context, header http.Header, body []byte, cfg Config) (*NotifyResult, error) {
	pubKey, err := parseRSAPublicKey(cfg.PlatformPublicKey)
	if err != nil {
		return nil, fmt.Errorf("微信平台公钥不可用：%w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/notify", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = header

	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		return nil, err
	}
	if err := notifyReq.VerifySignByPK(pubKey); err != nil {
		return nil, fmt.Errorf("微信回调验签失败：%w", err)
	}

	result, err := notifyReq.DecryptPayCipherText(cfg.APIKey)
	if err != nil {
		return nil, err
	}

	res := &NotifyResult{
		PaymentNo: result.OutTradeNo,
		TradeNo:   result.TransactionId,
		Success:   result.TradeState == wechat.TradeStateSuccess,
	}
	if result.Amount != nil {
		res.PaidAmount = int64(result.Amount.Total)
	}
	if result.Payer != nil {
		res.BuyerOpenID = result.Payer.Openid
	}

	return res, nil
}

// ParseRefundNotify 验签并解密退款回调。
// 退款回调与支付回调同一套平台公钥，解密后按 out_refund_no 定位退款单。
func (g *WechatGateway) ParseRefundNotify(ctx context.Context, header http.Header, body []byte, cfg Config) (*RefundNotifyResult, error) {
	pubKey, err := parseRSAPublicKey(cfg.PlatformPublicKey)
	if err != nil {
		return nil, fmt.Errorf("微信平台公钥不可用：%w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/notify", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = header

	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		return nil, err
	}
	if err := notifyReq.VerifySignByPK(pubKey); err != nil {
		return nil, fmt.Errorf("微信退款回调验签失败：%w", err)
	}

	result, err := notifyReq.DecryptRefundCipherText(cfg.APIKey)
	if err != nil {
		return nil, err
	}

	return &RefundNotifyResult{
		PaymentNo:       result.OutTradeNo,
		RefundNo:        result.OutRefundNo,
		ChannelRefundNo: result.RefundId,
		Success:         result.RefundStatus == wxRefundSuccess,
	}, nil
}

// Close 关闭微信侧未支付订单；用户已支付时微信会返回错误，由调用方忽略。
func (g *WechatGateway) Close(ctx context.Context, paymentNo string, cfg Config) error {
	client, err := g.client(cfg)
	if err != nil {
		return err
	}

	rsp, err := client.V3TransactionCloseOrder(ctx, paymentNo)
	if err != nil {
		return err
	}
	if rsp.Code != wechat.Success {
		return fmt.Errorf("微信关单失败(%d)：%s", rsp.Code, rsp.Error)
	}

	return nil
}

// client 用商户配置构造 gopay 客户端。
// 不开启请求自动验签：回调验签统一走 ParseNotify 里的公钥校验，避免依赖证书下载接口。
func (g *WechatGateway) client(cfg Config) (*wechat.ClientV3, error) {
	client, err := wechat.NewClientV3(cfg.MchID, cfg.SerialNo, cfg.APIKey, cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// payParams 生成前端调起参数：小程序用 Applet 结构，公众号用 JSAPI 结构。
func (g *WechatGateway) payParams(client *wechat.ClientV3, appID, prepayID string) (map[string]any, error) {
	var target any

	if g.scene == enums.SceneMini {
		params, err := client.PaySignOfApplet(appID, prepayID)
		if err != nil {
			return nil, err
		}
		target = params
	} else {
		params, err := client.PaySignOfJSAPI(appID, prepayID)
		if err != nil {
			return nil, err
		}
		target = params
	}

	// 结构体字段用 json tag 描述，转 map 后前端可直接透传给 wx.requestPayment
	bs, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	if err := json.Unmarshal(bs, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// refundNotifyURL 由支付回调地址推导退款回调地址。
// 约定：支付回调以 /payment 结尾时替换为 /refund，否则在末尾追加 /refund。
func refundNotifyURL(base string) string {
	if strings.HasSuffix(base, "/payment") {
		return strings.TrimSuffix(base, "/payment") + "/refund"
	}

	return strings.TrimSuffix(base, "/") + "/refund"
}

// parseRSAPublicKey 解析 PEM 形式的平台证书或公钥，支持 CERTIFICATE 与 PUBLIC KEY 两种块。
func parseRSAPublicKey(content string) (*rsa.PublicKey, error) {
	if content == "" {
		return nil, errors.New("未配置平台公钥")
	}

	block, _ := pem.Decode([]byte(content))
	if block == nil {
		return nil, errors.New("平台公钥不是合法的 PEM")
	}

	if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
		pub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("证书公钥不是 RSA")
		}

		return pub, nil
	}

	if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}

		return nil, errors.New("公钥不是 RSA")
	}

	return nil, errors.New("无法解析平台公钥")
}
