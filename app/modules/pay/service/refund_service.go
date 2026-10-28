package service

import (
	"context"
	"fmt"
	"net/http"

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

// RefundInput 是发起退款的入参。
// 金额按分传入，累计退款不得超过实付金额，校验在服务层完成。
type RefundInput struct {
	PaymentNo  string
	Amount     int64
	Reason     string
	OperatorID types.ID
}

// RefundListQuery 是后台退款单列表的查询条件。
type RefundListQuery struct {
	PaymentNo string
	MemberID  types.ID
	Status    *int8
	PageNo    int
	PageSize  int
}

// RefundService 负责退款单的创建、状态推进与回调处理。
// 余额支付原路退回钱包（立即成功），微信退款走渠道异步回调。
type RefundService struct {
	refunds  *repository.PayRefundRepository
	payments *repository.PayPaymentRepository
	wallets  *WalletService
	configs  *ChannelConfigService
	logs     *repository.PayNotifyLogRepository
}

// NewRefundService 构造退款服务。
func NewRefundService() *RefundService {
	return &RefundService{
		refunds:  repository.NewPayRefundRepository(),
		payments: repository.NewPayPaymentRepository(),
		wallets:  NewWalletService(),
		configs:  NewChannelConfigService(),
		logs:     repository.NewPayNotifyLogRepository(),
	}
}

// Create 发起退款。
// 余额支付：当次事务内退回钱包并置成功；
// 微信支付：调渠道退款，受理成功后置「退款中」，等回调或查单确认。
func (s *RefundService) Create(ctx context.Context, in RefundInput) (*model.PayRefundModel, error) {
	if in.Amount <= 0 {
		return nil, baseErrs.New(payErrs.PayWalletAmountInvalid)
	}

	payment, err := s.payments.FindByNo(ctx, in.PaymentNo)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, baseErrs.New(payErrs.PayPaymentNotExists, in.PaymentNo)
	}
	if payment.Status != payEnums.PaymentStatusPaid {
		return nil, baseErrs.New(payErrs.PayPaymentStatusInvalid, payment.PaymentNo)
	}
	if payment.PaidAmount-payment.RefundAmount < in.Amount {
		return nil, baseErrs.New(payErrs.PayRefundAmountExceed)
	}

	refund := &model.PayRefundModel{
		PaymentID:    payment.ID,
		PaymentNo:    payment.PaymentNo,
		BizType:      payment.BizType,
		BizID:        payment.BizID,
		MemberUserID: payment.MemberUserID,
		Channel:      payment.Channel,
		ChannelCode:  payment.ChannelCode,
		Amount:       in.Amount,
		Reason:       in.Reason,
		OperatorID:   in.OperatorID,
		Status:       payEnums.RefundStatusPending,
	}
	refund.ID = support.NewID()
	refund.RefundNo = fmt.Sprintf("R%d", int64(refund.ID))

	if err := s.refunds.Create(ctx, refund); err != nil {
		return nil, err
	}

	if payment.Channel == payEnums.ChannelBalance {
		if err := s.applyBalanceRefund(ctx, payment, refund); err != nil {
			return nil, err
		}

		return s.refunds.FindByNo(ctx, refund.RefundNo)
	}

	return s.refundByChannel(ctx, payment, refund)
}

// List 分页查询退款单。
func (s *RefundService) List(ctx context.Context, q RefundListQuery) ([]model.PayRefundModel, int64, error) {
	return s.refunds.Paginate(ctx, repository.RefundFilter{
		PaymentNo: q.PaymentNo,
		MemberID:  q.MemberID,
		Status:    q.Status,
	}, q.PageNo, q.PageSize)
}

// HandleWechatNotify 处理微信退款回调：验签 → 解密 → 幂等落库。
// 支付宝退款为同步返回，没有退款回调，因此这里只处理微信。
func (s *RefundService) HandleWechatNotify(ctx context.Context, header http.Header, body []byte) error {
	cfg, err := s.configs.AnyEnabled(ctx, payEnums.ChannelWechat)
	if err != nil {
		return err
	}
	gw, ok := gateway.Get(cfg.Code)
	if !ok {
		return baseErrs.New(payErrs.PayGatewayNotSupported, cfg.Code)
	}

	logRow := &model.PayNotifyLogModel{
		Channel:    payEnums.ChannelWechat,
		NotifyType: "refund",
		Body:       string(body),
	}

	result, err := gw.ParseRefundNotify(ctx, header, body, *cfg)
	if err != nil {
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = err.Error()
		_ = s.logs.Create(ctx, logRow)

		return err
	}

	logRow.Verified = 1
	logRow.OutTradeNo = result.RefundNo
	logRow.ChannelNo = result.ChannelRefundNo

	refund, err := s.refunds.FindByNo(ctx, result.RefundNo)
	if err != nil {
		return err
	}
	if refund == nil {
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = "退款单不存在"
		_ = s.logs.Create(ctx, logRow)

		return baseErrs.New(payErrs.PayRefundNotExists, result.RefundNo)
	}
	if refund.Status == payEnums.RefundStatusSuccess {
		logRow.HandleStatus = payEnums.NotifyHandleIgnored
		logRow.HandleMsg = "退款单已成功"
		_ = s.logs.Create(ctx, logRow)

		return nil
	}

	if result.Success {
		err = s.finish(ctx, refund, result.ChannelRefundNo)
	} else {
		err = facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
			return s.refunds.UpdateStatus(ctx, tx, refund.ID, payEnums.RefundStatusFailed, map[string]any{
				"fail_reason": "渠道退款失败",
			})
		})
	}
	if err != nil {
		logRow.HandleStatus = payEnums.NotifyHandleFailed
		logRow.HandleMsg = err.Error()
		_ = s.logs.Create(ctx, logRow)

		return err
	}

	logRow.HandleStatus = payEnums.NotifyHandleSuccess
	_ = s.logs.Create(ctx, logRow)

	return nil
}

// applyBalanceRefund 余额退款：退回钱包 + 更新退款单 + 累计支付单退款金额，同一事务完成。
func (s *RefundService) applyBalanceRefund(ctx context.Context, payment *model.PayPaymentModel, refund *model.PayRefundModel) error {
	err := facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		_, err := s.wallets.Credit(ctx, tx, payment.MemberUserID, payEnums.BizTypePrintRefund, payment.BizID, refund.Amount,
			balanceTitle(payEnums.BizTypePrintRefund, refund.Amount))
		if err != nil {
			return err
		}

		now := types.Now()
		if err := s.refunds.UpdateStatus(ctx, tx, refund.ID, payEnums.RefundStatusSuccess, map[string]any{
			"success_time": now,
		}); err != nil {
			return err
		}

		return s.addPaymentRefund(ctx, tx, payment, refund.Amount)
	})
	if err != nil {
		return err
	}

	// 余额退款当次即到账，事务成功后通知业务方
	notifyRefunded(ctx, refund.BizType, refund.BizID)

	return nil
}

// refundByChannel 调渠道退款；受理成功先落「退款中」，到账以回调为准。
func (s *RefundService) refundByChannel(ctx context.Context, payment *model.PayPaymentModel, refund *model.PayRefundModel) (*model.PayRefundModel, error) {
	cfg, err := s.configs.Config(ctx, payment.ChannelCode)
	if err != nil {
		return nil, err
	}
	gw, ok := gateway.Get(payment.ChannelCode)
	if !ok {
		return nil, baseErrs.New(payErrs.PayGatewayNotSupported, payment.ChannelCode)
	}

	result, err := gw.Refund(ctx, gateway.RefundInput{
		RefundNo:  refund.RefundNo,
		PaymentNo: payment.PaymentNo,
		Reason:    refund.Reason,
		Amount:    refund.Amount,
		Total:     payment.PaidAmount,
	}, *cfg)
	if err != nil {
		_ = facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
			return s.refunds.UpdateStatus(ctx, tx, refund.ID, payEnums.RefundStatusFailed, map[string]any{
				"fail_reason": truncateReason(err.Error()),
			})
		})

		return nil, baseErrs.New(payErrs.PayGatewayRequestFailed, err.Error())
	}

	if result.Success {
		if err := s.finish(ctx, refund, result.ChannelRefundNo); err != nil {
			return nil, err
		}

		return s.refunds.FindByNo(ctx, refund.RefundNo)
	}

	// 支付宝等同步渠道返回「未受理」时不存在后续回调，必须直接置失败
	if !result.Processing {
		_ = facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
			return s.refunds.UpdateStatus(ctx, tx, refund.ID, payEnums.RefundStatusFailed, map[string]any{
				"fail_reason": "渠道未受理退款",
			})
		})

		return nil, baseErrs.New(payErrs.PayGatewayRequestFailed, "渠道未受理退款")
	}

	err = facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		return s.refunds.UpdateStatus(ctx, tx, refund.ID, payEnums.RefundStatusRefunding, map[string]any{
			"channel_refund_no": result.ChannelRefundNo,
		})
	})
	if err != nil {
		return nil, err
	}

	return s.refunds.FindByNo(ctx, refund.RefundNo)
}

// finish 把退款单置成功并累计支付单退款金额。
// 全额退款时支付单转「已退款」，部分退款保持「已支付」。
func (s *RefundService) finish(ctx context.Context, refund *model.PayRefundModel, channelRefundNo string) error {
	err := facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		now := types.Now()

		if err := s.refunds.UpdateStatus(ctx, tx, refund.ID, payEnums.RefundStatusSuccess, map[string]any{
			"success_time":      now,
			"channel_refund_no": channelRefundNo,
		}); err != nil {
			return err
		}

		payment, err := s.payments.LockByID(ctx, tx, refund.PaymentID)
		if err != nil {
			return err
		}

		return s.addPaymentRefund(ctx, tx, payment, refund.Amount)
	})
	if err != nil {
		return err
	}

	// 微信退款是异步的，到这一步才确定资金已退，通知业务方落状态
	notifyRefunded(ctx, refund.BizType, refund.BizID)

	return nil
}

// addPaymentRefund 在事务内累加支付单退款金额并推导状态。
func (s *RefundService) addPaymentRefund(ctx context.Context, tx db.Tx, payment *model.PayPaymentModel, amount int64) error {
	payment.RefundAmount += amount
	if payment.RefundAmount >= payment.PaidAmount {
		payment.Status = payEnums.PaymentStatusRefunded
	}

	return s.payments.UpdateRefund(ctx, tx, payment)
}

// truncateReason 截断失败原因，避免超出 fail_reason 字段长度导致整笔回滚。
func truncateReason(reason string) string {
	runes := []rune(reason)
	if len(runes) <= 200 {
		return reason
	}

	return string(runes[:200])
}
