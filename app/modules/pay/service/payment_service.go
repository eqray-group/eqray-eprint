package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	payEnums "eqray-eprint/app/modules/pay/enums"
	payErrs "eqray-eprint/app/modules/pay/errs"
	"eqray-eprint/app/modules/pay/gateway"
	"eqray-eprint/app/modules/pay/model"
	"eqray-eprint/app/modules/pay/repository"
)

// paymentExpire 是支付单有效期，超时后需关单并重新下单。
const paymentExpire = 30 * time.Minute

// CreatePaymentInput 是创建支付单的入参。
// 金额由业务方（打印订单）算好后传入，支付模块不重新计价。
type CreatePaymentInput struct {
	MemberID    types.ID
	BizType     string
	BizID       types.ID
	Amount      int64
	ChannelCode string
	Description string
	OpenID      string
	ClientIP    string
}

// CreatePaymentOutput 是下单结果：余额支付已扣款，渠道支付需前端调起。
type CreatePaymentOutput struct {
	PaymentNo   string         `json:"paymentNo"`
	Paid        bool           `json:"paid"`
	ChannelCode string         `json:"channelCode"`
	Scene       string         `json:"scene"`
	Params      map[string]any `json:"params,omitempty"`
	MwebURL     string         `json:"mwebUrl,omitempty"`
	Balance     int64          `json:"balance,omitempty"`
}

// PaymentService 负责支付单的创建、查单与回调处理。
// 状态机：待支付 → 支付中 → 已支付；超时/取消 → 已关闭。
type PaymentService struct {
	payments *repository.PayPaymentRepository
	logs     *repository.PayNotifyLogRepository
	wallets  *WalletService
	configs  *ChannelConfigService
}

// NewPaymentService 构造支付单服务。
func NewPaymentService() *PaymentService {
	return &PaymentService{
		payments: repository.NewPayPaymentRepository(),
		logs:     repository.NewPayNotifyLogRepository(),
		wallets:  NewWalletService(),
		configs:  NewChannelConfigService(),
	}
}

// Create 创建支付单。
// 余额支付在事务内扣款并直接置已支付；渠道支付返回前端调起所需参数。
// 同一业务单存在未过期的同渠道支付单时复用，避免用户重复点击生成多张单。
func (s *PaymentService) Create(ctx context.Context, in CreatePaymentInput) (*CreatePaymentOutput, error) {
	if in.Amount <= 0 {
		return nil, baseErrs.New(payErrs.PayWalletAmountInvalid)
	}
	if !validBizType(in.BizType) {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "不支持的业务类型")
	}

	if in.ChannelCode == payEnums.ChannelBalance {
		return s.payByBalance(ctx, in)
	}

	cfg, err := s.configs.Config(ctx, in.ChannelCode)
	if err != nil {
		return nil, err
	}
	gw, ok := gateway.Get(in.ChannelCode)
	if !ok {
		return nil, baseErrs.New(payErrs.PayGatewayNotSupported, in.ChannelCode)
	}

	payment, err := s.pendingOrNew(ctx, in)
	if err != nil {
		return nil, err
	}

	result, err := gw.Create(ctx, gateway.Order{
		PaymentNo:   payment.PaymentNo,
		Description: in.Description,
		Amount:      payment.Amount,
		OpenID:      in.OpenID,
		ClientIP:    in.ClientIP,
		ExpireAt:    expireAt(payment),
	}, *cfg)
	if err != nil {
		return nil, baseErrs.New(payErrs.PayGatewayRequestFailed, err.Error())
	}

	if err := s.payments.MarkPaying(ctx, payment.ID, in.OpenID); err != nil {
		return nil, err
	}

	return &CreatePaymentOutput{
		PaymentNo:   payment.PaymentNo,
		ChannelCode: cfg.Code,
		Scene:       cfg.Scene,
		Params:      result.Params,
		MwebURL:     result.MwebURL,
	}, nil
}

// PaidPaymentOf 按业务定位已支付单，供业务模块发起退款与对账。
// 参数 bizType/bizID 与创建支付单时一致；未支付或不存在返回 (nil, nil)。
func (s *PaymentService) PaidPaymentOf(ctx context.Context, bizType string, bizID types.ID) (*model.PayPaymentModel, error) {
	return s.payments.FindPaidByBiz(ctx, bizType, bizID)
}

// AdminPaymentQuery 是后台支付流水列表的查询条件。
type AdminPaymentQuery struct {
	PaymentNo string
	MemberID  types.ID
	Channel   string
	Status    *int8
	PageNo    int
	PageSize  int
}

// AdminPage 分页查询支付单，供后台支付流水页使用。
func (s *PaymentService) AdminPage(ctx context.Context, q AdminPaymentQuery) ([]model.PayPaymentModel, int64, error) {
	return s.payments.Paginate(ctx, repository.PaymentFilter{
		PaymentNo: q.PaymentNo,
		MemberID:  q.MemberID,
		Channel:   q.Channel,
		Status:    q.Status,
	}, q.PageNo, q.PageSize)
}

// CloseExpired 关闭超时未支付的支付单，返回实际关闭数量。
// 先尽力关闭渠道侧订单（微信未必受理，失败只影响渠道侧），再本地置关闭；
// 本地条件更新保证回调先到时不会被误关。
func (s *PaymentService) CloseExpired(ctx context.Context, limit int) (int, error) {
	list, err := s.payments.FindExpired(ctx, limit)
	if err != nil {
		return 0, err
	}

	closed := 0
	for i := range list {
		payment := &list[i]

		if payment.ChannelCode != "" && payment.ChannelCode != payEnums.ChannelBalance {
			if cfg, cfgErr := s.configs.Config(ctx, payment.ChannelCode); cfgErr == nil {
				if gw, ok := gateway.Get(payment.ChannelCode); ok {
					_ = gw.Close(ctx, payment.PaymentNo, *cfg)
				}
			}
		}

		rows, err := s.payments.MarkClosed(ctx, payment.ID)
		if err != nil {
			return closed, err
		}
		if rows > 0 {
			closed++
		}
	}

	return closed, nil
}

// Recharge 创建钱包充值支付单。
// 充值的 biz_id 取钱包 ID，回调入账时按它定位钱包，避免依赖会员 ID 反查。
func (s *PaymentService) Recharge(ctx context.Context, memberID types.ID, amount int64, channelCode, openID, clientIP string) (*CreatePaymentOutput, error) {
	wallet, err := s.wallets.Ensure(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return s.Create(ctx, CreatePaymentInput{
		MemberID:    memberID,
		BizType:     payEnums.BizTypeWalletRecharge,
		BizID:       wallet.ID,
		Amount:      amount,
		ChannelCode: channelCode,
		Description: "会员钱包充值",
		OpenID:      openID,
		ClientIP:    clientIP,
	})
}

// Detail 查询支付单，仅允许查询本人支付单。
func (s *PaymentService) Detail(ctx context.Context, memberID types.ID, paymentNo string) (*model.PayPaymentModel, error) {
	payment, err := s.mustGet(ctx, memberID, paymentNo)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// List 分页查询会员支付单。
func (s *PaymentService) List(ctx context.Context, memberID types.ID, pageNo, pageSize int) ([]model.PayPaymentModel, int64, error) {
	return s.payments.PaginateByMember(ctx, memberID, pageNo, pageSize)
}

// Sync 主动查单：渠道侧已支付而回调丢失时，用查单结果补记状态。
func (s *PaymentService) Sync(ctx context.Context, memberID types.ID, paymentNo string) (*model.PayPaymentModel, error) {
	payment, err := s.mustGet(ctx, memberID, paymentNo)
	if err != nil {
		return nil, err
	}
	if payment.Status == payEnums.PaymentStatusPaid || payment.Status == payEnums.PaymentStatusRefunded {
		return payment, nil
	}

	cfg, err := s.configs.Config(ctx, payment.ChannelCode)
	if err != nil {
		return nil, err
	}
	gw, ok := gateway.Get(payment.ChannelCode)
	if !ok {
		return nil, baseErrs.New(payErrs.PayGatewayNotSupported, payment.ChannelCode)
	}

	result, err := gw.Query(ctx, payment.PaymentNo, *cfg)
	if err != nil {
		return nil, baseErrs.New(payErrs.PayGatewayRequestFailed, err.Error())
	}
	if !result.Paid {
		return payment, nil
	}
	if result.PaidAmount != payment.Amount {
		return nil, baseErrs.New(payErrs.PayPaymentAmountMismatch)
	}

	if err := s.markPaid(ctx, payment, result.TradeNo, result.PaidAmount); err != nil {
		return nil, err
	}

	return s.payments.FindByNo(ctx, paymentNo)
}

// HandleWechatNotify 处理微信支付回调，等价于 HandleNotify(ctx, wechat, ...)。
// 返回 nil 表示可以给微信回 success；返回错误时微信会重试，因此必须保证幂等。
func (s *PaymentService) HandleWechatNotify(ctx context.Context, header http.Header, body []byte) error {
	return s.HandleNotify(ctx, payEnums.ChannelWechat, header, body)
}

// HandleNotify 处理渠道支付回调：验签 → 解密 → 金额校验 → 幂等落库。
// 参数 channel 决定用哪套渠道配置验签（wechat / alipay）；
// 返回 nil 表示可向渠道回成功，返回错误时渠道会重试。
func (s *PaymentService) HandleNotify(ctx context.Context, channel string, header http.Header, body []byte) error {
	cfg, err := s.configs.AnyEnabled(ctx, channel)
	if err != nil {
		return err
	}
	gw, ok := gateway.Get(cfg.Code)
	if !ok {
		return baseErrs.New(payErrs.PayGatewayNotSupported, cfg.Code)
	}

	logRow := &model.PayNotifyLogModel{
		Channel:    channel,
		NotifyType: "payment",
		Body:       string(body),
	}

	result, err := gw.ParseNotify(ctx, header, body, *cfg)
	if err != nil {
		logRow.Verified = 0
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = err.Error()
		_ = s.logs.Create(ctx, logRow)

		return err
	}

	logRow.Verified = 1
	logRow.OutTradeNo = result.PaymentNo
	logRow.ChannelNo = result.TradeNo

	payment, err := s.payments.FindByNo(ctx, result.PaymentNo)
	if err != nil {
		return err
	}
	if payment == nil {
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = "支付单不存在"
		_ = s.logs.Create(ctx, logRow)

		return baseErrs.New(payErrs.PayPaymentNotExists, result.PaymentNo)
	}
	// 金额不一致绝不入账，落日志留证后返回错误
	if result.PaidAmount != payment.Amount {
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = "回调金额与支付单不一致"
		_ = s.logs.Create(ctx, logRow)

		return baseErrs.New(payErrs.PayPaymentAmountMismatch)
	}
	// 已处理过：直接回 success，微信不再重试
	if payment.Status == payEnums.PaymentStatusPaid || payment.Status == payEnums.PaymentStatusRefunded {
		logRow.HandleStatus = payEnums.NotifyHandleIgnored
		logRow.HandleMsg = "支付单已处理"
		_ = s.logs.Create(ctx, logRow)

		return nil
	}

	if err := s.markPaid(ctx, payment, result.TradeNo, result.PaidAmount); err != nil {
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = err.Error()
		_ = s.logs.Create(ctx, logRow)

		return err
	}

	logRow.HandleStatus = payEnums.NotifyHandleSuccess
	_ = s.logs.Create(ctx, logRow)

	return nil
}

// payByBalance 用钱包余额支付：先落支付单，再在事务内扣款 + 写流水 + 置已支付。
func (s *PaymentService) payByBalance(ctx context.Context, in CreatePaymentInput) (*CreatePaymentOutput, error) {
	if in.BizType != payEnums.BizTypePrintOrder {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "充值不支持余额支付")
	}

	payment, err := s.newPayment(ctx, in, payEnums.ChannelBalance)
	if err != nil {
		return nil, err
	}

	var balance int64
	err = facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		var err error
		balance, err = s.wallets.Debit(ctx, tx, in.MemberID, payEnums.BizTypePrintConsume, in.BizID, in.Amount,
			balanceTitle(payEnums.BizTypePrintConsume, in.Amount))
		if err != nil {
			return err
		}

		now := types.Now()

		return s.payments.UpdateStatus(ctx, tx, payment.ID, payEnums.PaymentStatusPaid, map[string]any{
			"paid_amount":  in.Amount,
			"success_time": now,
		})
	})
	if err != nil {
		return nil, err
	}

	// 余额支付当次即成功，直接触发业务回调（如打印订单派发任务）
	notifyPaid(ctx, in.BizType, in.BizID)

	return &CreatePaymentOutput{
		PaymentNo:   payment.PaymentNo,
		Paid:        true,
		ChannelCode: payEnums.ChannelBalance,
		Balance:     balance,
	}, nil
}

// markPaid 把支付单置为已支付；充值时同时给钱包入账，两步在同一事务内完成。
// 事务成功后触发业务回调（如打印订单派发任务），回调失败只记日志。
func (s *PaymentService) markPaid(ctx context.Context, payment *model.PayPaymentModel, tradeNo string, paidAmount int64) error {
	changed := false

	err := facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		now := types.Now()

		// 条件更新：只有非已支付状态才会被改动，天然幂等
		res, err := tx.Table(model.PayPaymentTable).
			Where("id", payment.ID).
			Where("status", "<", payEnums.PaymentStatusPaid).
			Update(map[string]any{
				"status":           payEnums.PaymentStatusPaid,
				"paid_amount":      paidAmount,
				"channel_trade_no": tradeNo,
				"success_time":     now,
				"updated_at":       now,
			})
		if err != nil {
			return err
		}
		if res.RowsAffected == 0 {
			return nil
		}
		changed = true

		if payment.BizType == payEnums.BizTypeWalletRecharge {
			_, err = s.wallets.Credit(ctx, tx, payment.MemberUserID, payEnums.BizTypeWalletRecharge, payment.ID, paidAmount,
				balanceTitle(payEnums.BizTypeWalletRecharge, paidAmount))
		}

		return err
	})
	if err != nil {
		return err
	}

	// 只有本次真正改了状态才通知业务方，重复回调不会重复派发
	if changed {
		notifyPaid(ctx, payment.BizType, payment.BizID)
	}

	return nil
}

// pendingOrNew 复用未过期的待支付单，否则新建一张。
func (s *PaymentService) pendingOrNew(ctx context.Context, in CreatePaymentInput) (*model.PayPaymentModel, error) {
	pending, err := s.payments.FindPending(ctx, in.BizType, in.BizID, in.ChannelCode)
	if err != nil {
		return nil, err
	}
	if pending != nil && pending.Amount == in.Amount && !expired(pending) {
		return pending, nil
	}

	return s.newPayment(ctx, in, in.ChannelCode)
}

// newPayment 生成支付单号并落库，状态为待支付。
func (s *PaymentService) newPayment(ctx context.Context, in CreatePaymentInput, channel string) (*model.PayPaymentModel, error) {
	now := types.Now()
	expire := types.DateTime(now.Time().Add(paymentExpire))

	channelCode := in.ChannelCode
	if channel == payEnums.ChannelBalance {
		channelCode = payEnums.ChannelBalance
	}

	payment := &model.PayPaymentModel{
		BizType:      in.BizType,
		BizID:        in.BizID,
		MemberUserID: in.MemberID,
		// Channel 只存渠道大类（wechat / alipay / balance），
		// 具体场景由 ChannelCode 决定，统计与筛选都按大类走
		Channel:     channelOf(channelCode),
		ChannelCode: channelCode,
		Scene:       sceneOf(channelCode),
		ClientIP:    in.ClientIP,
		Amount:      in.Amount,
		Status:      payEnums.PaymentStatusPending,
		OpenID:      in.OpenID,
		ExpireTime:  &expire,
	}
	payment.ID = support.NewID()
	payment.PaymentNo = fmt.Sprintf("P%d", int64(payment.ID))

	if err := s.payments.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// mustGet 取当前会员的支付单，越权查询按不存在处理。
func (s *PaymentService) mustGet(ctx context.Context, memberID types.ID, paymentNo string) (*model.PayPaymentModel, error) {
	payment, err := s.payments.FindByNo(ctx, paymentNo)
	if err != nil {
		return nil, err
	}
	if payment == nil || payment.MemberUserID != memberID {
		return nil, baseErrs.New(payErrs.PayPaymentNotExists, paymentNo)
	}

	return payment, nil
}

// expired 判断支付单是否已过有效期。
func expired(payment *model.PayPaymentModel) bool {
	if payment.ExpireTime == nil {
		return false
	}

	return payment.ExpireTime.Time().Before(types.Now().Time())
}

// expireAt 返回支付单过期时间，空值按当前加默认有效期。
func expireAt(payment *model.PayPaymentModel) time.Time {
	if payment.ExpireTime != nil {
		return payment.ExpireTime.Time()
	}

	return time.Now().Add(paymentExpire)
}

// validBizType 限定支付单可承载的业务类型，避免乱传导致后续无法处理。
func validBizType(bizType string) bool {
	return bizType == payEnums.BizTypePrintOrder || bizType == payEnums.BizTypeWalletRecharge
}

// channelOf 由渠道配置编码推断渠道大类，用于支付单统计口径。
func channelOf(channelCode string) string {
	switch channelCode {
	case payEnums.ChannelBalance:
		return payEnums.ChannelBalance
	case payEnums.CodeAlipayWap:
		return payEnums.ChannelAlipay
	default:
		return payEnums.ChannelWechat
	}
}

// sceneOf 按渠道配置编码推断场景，写入支付单便于统计。
func sceneOf(channelCode string) string {
	switch channelCode {
	case payEnums.CodeWechatMini:
		return payEnums.SceneMini
	case payEnums.CodeWechatMp:
		return payEnums.SceneJsapi
	case payEnums.CodeWechatH5:
		return payEnums.SceneH5
	case payEnums.CodeAlipayWap:
		return payEnums.SceneWap
	default:
		return ""
	}
}
