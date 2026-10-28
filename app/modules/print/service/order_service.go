package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	memberService "eqray-eprint/app/modules/member/service"
	payEnums "eqray-eprint/app/modules/pay/enums"
	payService "eqray-eprint/app/modules/pay/service"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

// 订单相关超时：未支付 30 分钟关单，取件保留 7 天。
const (
	payTimeout    = 30 * time.Minute
	pickupKeep    = 7 * 24 * time.Hour
	maxTaskRetry  = 3
	taskPullLimit = 20
)

// CreateItemInput 是下单时的一个文件及其打印参数。
type CreateItemInput struct {
	FileID      types.ID
	PaperType   string
	ColorType   string
	PrintMode   string
	Orientation string
	PageRange   string
	Copies      int
}

// CreateOrderInput 是会员下单入参。
type CreateOrderInput struct {
	MemberID    types.ID
	ShopID      types.ID
	ChannelCode string
	OpenID      string
	ClientIP    string
	Remark      string
	Items       []CreateItemInput
}

// CreateOrderOutput 返回订单号与支付调起信息。
type CreateOrderOutput struct {
	OrderNo   string         `json:"orderNo"`
	PaymentNo string         `json:"paymentNo"`
	PayAmount int64          `json:"payAmount"`
	Paid      bool           `json:"paid"`
	Params    map[string]any `json:"params,omitempty"`
	MwebURL   string         `json:"mwebUrl,omitempty"`
	Balance   int64          `json:"balance,omitempty"`
}

// TaskView 是下发给客户端的打印任务，参数以服务端为准，客户端只读。
type TaskView struct {
	TaskNo        string         `json:"taskNo"`
	OrderNo       string         `json:"orderNo"`
	FileID        types.ID       `json:"fileId"`
	FileName      string         `json:"fileName"`
	FilePageCount int            `json:"filePageCount"`
	PaperType     string         `json:"paperType"`
	ColorType     string         `json:"colorType"`
	PrintMode     string         `json:"printMode"`
	Orientation   string         `json:"orientation"`
	PageRange     string         `json:"pageRange"`
	PageCount     int            `json:"pageCount"`
	Copies        int            `json:"copies"`
	SystemName    string         `json:"systemName"`
	PrintConfig   map[string]any `json:"printConfig,omitempty"`
}

// ReportTaskInput 是客户端回执入参。
type ReportTaskInput struct {
	DeviceID     types.ID
	TaskNo       string
	Status       int8
	PrintedPages int
	ErrorCode    string
	Reason       string
}

// OrderListQuery 是订单列表查询条件。
type OrderListQuery struct {
	ShopID      types.ID
	MemberID    types.ID
	OrderNo     string
	PayStatus   *int8
	OrderStatus *int8
	PageNo      int
	PageSize    int
}

// OrderService 负责下单、支付后派发、任务回执与取件核销。
type OrderService struct {
	orders     *repository.PrintOrderRepository
	items      *repository.PrintOrderItemRepository
	files      *repository.PrintFileRepository
	capability *CapabilityService
	shops      *ShopService
	payments   *payService.PaymentService
}

// NewOrderService 构造订单服务。
func NewOrderService() *OrderService {
	return &OrderService{
		orders:     repository.NewPrintOrderRepository(),
		items:      repository.NewPrintOrderItemRepository(),
		files:      repository.NewPrintFileRepository(),
		capability: NewCapabilityService(),
		shops:      NewShopService(),
		payments:   payService.NewPaymentService(),
	}
}

// Create 会员下单：校验门店与文件 → 逐项匹配打印机并计价 → 落订单与明细 → 创建支付单。
// 价格与折扣全部由服务端计算，前端传入的金额一律忽略。
func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (*CreateOrderOutput, error) {
	if len(in.Items) == 0 {
		return nil, baseErrs.New(printErrs.OrderItemEmpty)
	}

	shop, err := s.shops.Detail(ctx, in.ShopID)
	if err != nil {
		return nil, err
	}
	if shop.AuditStatus != printEnums.AuditStatusApproved {
		return nil, baseErrs.New(printErrs.ShopNotAudited)
	}
	if shop.Status != printEnums.ShopStatusOpen {
		return nil, baseErrs.New(printErrs.ShopClosed)
	}

	items, totals, err := s.buildItems(ctx, in)
	if err != nil {
		return nil, err
	}

	percent, err := memberService.NewMemberCardService().DiscountPercent(ctx, in.MemberID)
	if err != nil {
		return nil, err
	}
	discountAmount := totals.TotalAmount * int64(100-percent) / 100
	payAmount := totals.TotalAmount - discountAmount

	now := types.Now()
	payExpire := types.DateTime(now.Time().Add(payTimeout))

	order := &model.PrintOrderModel{
		OrderNo:         fmt.Sprintf("PO%d", int64(support.NewID())),
		PrintShopID:     in.ShopID,
		MemberUserID:    in.MemberID,
		DeviceID:        totals.DeviceID,
		PickupCode:      s.generatePickupCode(ctx, in.ShopID),
		ItemCount:       len(items),
		TotalPages:      totals.TotalPages,
		TotalSheets:     totals.TotalSheets,
		TotalAmount:     totals.TotalAmount,
		DiscountPercent: percent,
		DiscountAmount:  discountAmount,
		PayAmount:       payAmount,
		PayStatus:       printEnums.PayStatusUnpaid,
		OrderStatus:     printEnums.OrderStatusUnpaid,
		PayExpireTime:   &payExpire,
		Remark:          strings.TrimSpace(in.Remark),
	}
	order.ID = support.NewID()
	order.PayType = payTypeOf(in.ChannelCode)

	for i := range items {
		items[i].OrderID = order.ID
	}

	if err := facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		if err := s.orders.CreateInTx(ctx, tx, order); err != nil {
			return err
		}

		return s.items.BatchCreate(ctx, tx, items)
	}); err != nil {
		return nil, err
	}

	return s.pay(ctx, in, order)
}

// PayOrder 为已创建但未支付的订单重新发起支付，用于用户关掉收银台后再次支付。
func (s *OrderService) PayOrder(ctx context.Context, memberID types.ID, orderNo, channelCode, openID, clientIP string) (*CreateOrderOutput, error) {
	order, err := s.orders.FindByNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	if order == nil || order.MemberUserID != memberID {
		return nil, baseErrs.New(printErrs.OrderNotExists, orderNo)
	}
	if order.OrderStatus != printEnums.OrderStatusUnpaid {
		return nil, baseErrs.New(printErrs.OrderStatusInvalid)
	}

	return s.pay(ctx, CreateOrderInput{
		MemberID:    memberID,
		ShopID:      order.PrintShopID,
		ChannelCode: channelCode,
		OpenID:      openID,
		ClientIP:    clientIP,
	}, order)
}

// OnPaid 是支付成功回调：把订单置为已支付并转「待打印」。
// bizType 非 print_order 时直接返回，支付模块只负责通知，不认识打印业务。
func (s *OrderService) OnPaid(ctx context.Context, bizType string, bizID types.ID) error {
	if bizType != payEnums.BizTypePrintOrder {
		return nil
	}

	order, err := s.orders.FindByID(ctx, bizID)
	if err != nil {
		return err
	}
	if order == nil {
		return baseErrs.New(printErrs.OrderNotExists, bizID)
	}
	if order.PayStatus == printEnums.PayStatusPaid {
		return nil
	}

	return facades.DB().WithContext(ctx).Transaction(func(tx db.Tx) error {
		rows, err := s.orders.MarkPaid(ctx, tx, order.ID, order.PayType, "")
		if err != nil {
			return err
		}
		if rows == 0 {
			// 已被其它回调处理，幂等返回
			return nil
		}

		return nil
	})
}

// Refund 对已支付订单发起全额退款；打印失败自动退款与后台人工退款共用。
// 参数 orderID 为订单 ID，reason 为退款原因，operatorID 为操作人（0 表示系统）。
// 幂等：已退款直接返回，退款中重复调用会再次进入支付模块，由支付单金额校验兜住。
func (s *OrderService) Refund(ctx context.Context, orderID types.ID, reason string, operatorID types.ID) error {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return baseErrs.New(printErrs.OrderNotExists, orderID)
	}
	if order.PayStatus == printEnums.PayStatusRefunded {
		return nil
	}
	if order.PayStatus != printEnums.PayStatusPaid {
		return baseErrs.New(printErrs.OrderStatusInvalid)
	}

	payment, err := s.payments.PaidPaymentOf(ctx, payEnums.BizTypePrintOrder, order.ID)
	if err != nil {
		return err
	}
	if payment == nil {
		return baseErrs.WithMsg(baseErrs.BadRequest, "未找到该订单的支付单，无法退款")
	}

	// 先落「退款中」，避免并发重复发起；失败时按原状态回滚
	previousStatus := order.OrderStatus
	if err := s.orders.UpdateStatus(ctx, order.ID, printEnums.OrderStatusRefunding, map[string]any{
		"pay_status":    printEnums.PayStatusRefunding,
		"cancel_reason": truncateText(reason, 255),
	}); err != nil {
		return err
	}

	if _, err := payService.NewRefundService().Create(ctx, payService.RefundInput{
		PaymentNo:  payment.PaymentNo,
		Amount:     order.PayAmount,
		Reason:     reason,
		OperatorID: operatorID,
	}); err != nil {
		// 退款没发出去，订单退回原状态（仍为已支付），便于重试
		_ = s.orders.UpdateStatus(ctx, order.ID, previousStatus, map[string]any{
			"pay_status": printEnums.PayStatusPaid,
		})

		return err
	}

	return nil
}

// OnRefunded 是退款成功回调：订单置「已退款」并取消未开始的任务。
// bizType 非 print_order 时直接返回，支付模块不认识打印业务。
func (s *OrderService) OnRefunded(ctx context.Context, bizType string, bizID types.ID) error {
	if bizType != payEnums.BizTypePrintOrder {
		return nil
	}

	order, err := s.orders.FindByID(ctx, bizID)
	if err != nil {
		return err
	}
	if order == nil {
		return nil
	}
	if order.PayStatus == printEnums.PayStatusRefunded {
		return nil
	}

	now := types.Now()
	if err := s.orders.UpdateStatus(ctx, order.ID, printEnums.OrderStatusRefunded, map[string]any{
		"pay_status":    printEnums.PayStatusRefunded,
		"refund_amount": order.PayAmount,
		"cancel_time":   now,
	}); err != nil {
		return err
	}

	return s.items.CancelByOrder(ctx, order.ID)
}

// Detail 查询会员自己的订单（含明细）。
func (s *OrderService) Detail(ctx context.Context, memberID types.ID, orderNo string) (*model.PrintOrderModel, []model.PrintOrderItemModel, error) {
	order, err := s.orders.FindByNo(ctx, orderNo)
	if err != nil {
		return nil, nil, err
	}
	if order == nil || order.MemberUserID != memberID {
		return nil, nil, baseErrs.New(printErrs.OrderNotExists, orderNo)
	}

	items, err := s.items.ListByOrder(ctx, order.ID)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

// AdminDetail 按主键查询订单详情（含明细），供后台订单详情弹窗使用。
// 后台不做会员归属校验，只要能进这个接口就有 print:order:list 权限。
func (s *OrderService) AdminDetail(ctx context.Context, id types.ID) (*model.PrintOrderModel, []model.PrintOrderItemModel, error) {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if order == nil {
		return nil, nil, baseErrs.New(printErrs.OrderNotExists, id)
	}

	items, err := s.items.ListByOrder(ctx, order.ID)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

// List 分页查询订单，会员端传 MemberID，后台传 ShopID 或都不传。
func (s *OrderService) List(ctx context.Context, q OrderListQuery) ([]model.PrintOrderModel, int64, error) {
	return s.orders.Paginate(ctx, repository.OrderFilter{
		ShopID:      q.ShopID,
		MemberID:    q.MemberID,
		OrderNo:     strings.TrimSpace(q.OrderNo),
		PayStatus:   q.PayStatus,
		OrderStatus: q.OrderStatus,
	}, q.PageNo, q.PageSize)
}

// ClaimTasks 客户端拉取本机待打印任务。
// 参数 deviceID 为设备 ID，limit 为单次上限；只返回已分配且未开始的任务。
func (s *OrderService) ClaimTasks(ctx context.Context, deviceID types.ID, limit int) ([]TaskView, error) {
	if limit <= 0 || limit > taskPullLimit {
		limit = taskPullLimit
	}

	items, err := s.items.ListPendingByDevice(ctx, deviceID, limit)
	if err != nil {
		return nil, err
	}

	views := make([]TaskView, 0, len(items))
	for _, item := range items {
		order, err := s.orders.FindByID(ctx, item.OrderID)
		if err != nil {
			return nil, err
		}
		// 未支付或已关闭的订单绝不下发给设备，防止关单后仍被打印
		if order == nil || order.PayStatus != printEnums.PayStatusPaid {
			continue
		}
		orderNo := order.OrderNo

		// 打印队列名从打印机表取，客户端按它静默打印
		systemName := ""
		printer, err := NewPrinterService().Detail(ctx, item.PrinterID)
		if err != nil || printer == nil {
			// 打印机被停用或删除时任务保持已分配，由后台处理，不下发脏任务
			continue
		}
		systemName = printer.SystemName

		views = append(views, TaskView{
			TaskNo:        item.TaskNo,
			OrderNo:       orderNo,
			FileID:        item.FileID,
			FileName:      item.FileName,
			FilePageCount: item.FilePageCount,
			PaperType:     item.PaperType,
			ColorType:     item.ColorType,
			PrintMode:     item.PrintMode,
			Orientation:   item.Orientation,
			PageRange:     item.PageRange,
			PageCount:     item.PageCount,
			Copies:        item.Copies,
			SystemName:    systemName,
			PrintConfig:   decodeConfig(item.PrintConfig),
		})
	}

	return views, nil
}

// ReportTask 处理客户端回执：打印中 / 完成 / 失败。
// 参数 in 为回执内容；全部任务完成时订单自动转「待取件」。
func (s *OrderService) ReportTask(ctx context.Context, in ReportTaskInput) error {
	item, err := s.items.FindByTaskNo(ctx, strings.TrimSpace(in.TaskNo))
	if err != nil {
		return err
	}
	if item == nil {
		return baseErrs.New(printErrs.TaskNotExists)
	}
	if item.DeviceID != in.DeviceID {
		return baseErrs.New(printErrs.TaskNotExists)
	}

	switch in.Status {
	case printEnums.TaskStatusPrinting:
		return s.items.MarkPrinting(ctx, item.ID)
	case printEnums.TaskStatusDone:
		if err := s.items.MarkDone(ctx, item.ID, in.PrintedPages); err != nil {
			return err
		}

		return s.completeIfAllDone(ctx, item.OrderID)
	case printEnums.TaskStatusFailed:
		return s.handleTaskFailure(ctx, item, in)
	default:
		return baseErrs.WithMsg(baseErrs.BadRequest, "回执状态不合法")
	}
}

// PickupVerify 门店核销取件码：仅待取件订单可核销，记录核销人与时间。
func (s *OrderService) PickupVerify(ctx context.Context, shopID types.ID, code string, operatorID types.ID) (*model.PrintOrderModel, error) {
	order, err := s.orders.FindByPickup(ctx, shopID, strings.TrimSpace(code))
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, baseErrs.New(printErrs.OrderPickupInvalid)
	}

	now := types.Now()
	if err := s.orders.UpdateStatus(ctx, order.ID, printEnums.OrderStatusCompleted, map[string]any{
		"finish_time":           now,
		"pickup_verify_at":      now,
		"pickup_verify_user_id": operatorID,
	}); err != nil {
		return nil, err
	}

	return s.orders.FindByID(ctx, order.ID)
}

// pay 创建支付单并把金额交给支付模块；金额完全由服务端计算。
func (s *OrderService) pay(ctx context.Context, in CreateOrderInput, order *model.PrintOrderModel) (*CreateOrderOutput, error) {
	if order.PayAmount <= 0 {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "订单金额异常，请重新下单")
	}

	result, err := s.payments.Create(ctx, payService.CreatePaymentInput{
		MemberID:    in.MemberID,
		BizType:     payEnums.BizTypePrintOrder,
		BizID:       order.ID,
		Amount:      order.PayAmount,
		ChannelCode: in.ChannelCode,
		Description: "打印订单 " + order.OrderNo,
		OpenID:      in.OpenID,
		ClientIP:    in.ClientIP,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderOutput{
		OrderNo:   order.OrderNo,
		PaymentNo: result.PaymentNo,
		PayAmount: order.PayAmount,
		Paid:      result.Paid,
		Params:    result.Params,
		MwebURL:   result.MwebURL,
		Balance:   result.Balance,
	}, nil
}

// buildItems 校验文件与参数、匹配打印机并计价，返回待插入的明细与汇总。
func (s *OrderService) buildItems(ctx context.Context, in CreateOrderInput) ([]model.PrintOrderItemModel, orderTotals, error) {
	fileIDs := make([]types.ID, 0, len(in.Items))
	for _, item := range in.Items {
		fileIDs = append(fileIDs, item.FileID)
	}

	files, err := s.files.FindByIDs(ctx, fileIDs)
	if err != nil {
		return nil, orderTotals{}, err
	}

	items := make([]model.PrintOrderItemModel, 0, len(in.Items))
	var totals orderTotals

	for _, input := range in.Items {
		file, ok := files[input.FileID]
		if !ok || file.MemberUserID != in.MemberID {
			return nil, totals, baseErrs.New(printErrs.FileNotExists)
		}
		if file.ConvertStatus != printEnums.ConvertStatusReady || file.PageCount <= 0 {
			return nil, totals, baseErrs.New(printErrs.FileNotReady)
		}

		paper := normalizePaper(input.PaperType)
		color := normalizeColor(input.ColorType)
		mode := normalizeMode(input.PrintMode)
		if paper == "" || color == "" || mode == "" {
			return nil, totals, baseErrs.WithMsg(baseErrs.BadRequest, "打印参数不合法")
		}
		orientation := input.Orientation
		if orientation != printEnums.OrientationLandscape {
			orientation = printEnums.OrientationPortrait
		}

		copies := input.Copies
		if copies <= 0 {
			copies = 1
		}

		// 严格匹配打印机：无匹配直接报错，不降级
		printer, err := s.capability.ResolvePrinter(ctx, in.ShopID, paper, color, mode)
		if err != nil {
			return nil, totals, err
		}
		price, err := s.capability.Price(ctx, in.ShopID, paper, color, mode)
		if err != nil {
			return nil, totals, err
		}

		pageCount := countPages(input.PageRange, file.PageCount)
		if pageCount <= 0 {
			return nil, totals, baseErrs.WithMsg(baseErrs.BadRequest, "页码范围超出文件页数")
		}

		sheets := pageCount * copies
		amount := int64(sheets) * price
		config := buildPrintConfig(paper, color, mode, orientation, input.PageRange, copies)

		item := model.PrintOrderItemModel{
			PrintShopID:   in.ShopID,
			TaskNo:        fmt.Sprintf("T%d", int64(support.NewID())),
			FileID:        file.ID,
			FileName:      file.FileName,
			FilePageCount: file.PageCount,
			PaperType:     paper,
			ColorType:     color,
			PrintMode:     mode,
			Orientation:   orientation,
			PageRange:     strings.TrimSpace(input.PageRange),
			PageCount:     pageCount,
			Copies:        copies,
			SheetCount:    sheets,
			PrintConfig:   config,
			Amount:        amount,
			Status:        printEnums.TaskStatusAssigned,
			DeviceID:      printer.DeviceID,
			PrinterID:     printer.ID,
		}
		item.ID = support.NewID()

		totals.TotalAmount += amount
		totals.TotalPages += pageCount
		totals.TotalSheets += sheets
		if totals.DeviceID == 0 {
			totals.DeviceID = printer.DeviceID
		}

		items = append(items, item)
	}

	return items, totals, nil
}

// handleTaskFailure 处理失败回执：不支持的能力直接失败，其它按次数重试，超限退款。
func (s *OrderService) handleTaskFailure(ctx context.Context, item *model.PrintOrderItemModel, in ReportTaskInput) error {
	code := in.ErrorCode
	if code == "" {
		code = printEnums.ErrorUnknown
	}

	permanent := code == printEnums.ErrorColorUnsupported ||
		code == printEnums.ErrorPaperUnsupported ||
		code == printEnums.ErrorDuplexUnsupported ||
		code == printEnums.ErrorFileBroken

	retry := item.RetryCount + 1
	if permanent || retry >= maxTaskRetry {
		if err := s.items.MarkFailed(ctx, item.ID, code, in.Reason, retry); err != nil {
			return err
		}

		return s.refundIfAllFailed(ctx, item.OrderID)
	}

	return s.items.Retry(ctx, item.ID, retry)
}

// completeIfAllDone 所有任务完成时把订单转「待取件」并设置取件保留期。
func (s *OrderService) completeIfAllDone(ctx context.Context, orderID types.ID) error {
	pending, err := s.items.CountByStatus(ctx, orderID, printEnums.TaskStatusAssigned)
	if err != nil {
		return err
	}
	printing, err := s.items.CountByStatus(ctx, orderID, printEnums.TaskStatusPrinting)
	if err != nil {
		return err
	}
	if pending+printing > 0 {
		return nil
	}

	pickupExpire := types.DateTime(types.Now().Time().Add(pickupKeep))

	return s.orders.UpdateStatus(ctx, orderID, printEnums.OrderStatusPickup, map[string]any{
		"pickup_expire_time": pickupExpire,
	})
}

// refundIfAllFailed 全部任务失败时把订单转「退款中」，退款由后台或自动任务发起。
func (s *OrderService) refundIfAllFailed(ctx context.Context, orderID types.ID) error {
	failed, err := s.items.CountByStatus(ctx, orderID, printEnums.TaskStatusFailed)
	if err != nil {
		return err
	}
	pending, err := s.items.CountByStatus(ctx, orderID, printEnums.TaskStatusAssigned)
	if err != nil {
		return err
	}
	printing, err := s.items.CountByStatus(ctx, orderID, printEnums.TaskStatusPrinting)
	if err != nil {
		return err
	}
	if failed == 0 || pending+printing > 0 {
		return nil
	}

	// 状态推进与退款都交给 Refund：它负责置「退款中」并在失败时回滚原状态
	// 自动退款失败不阻断回执流程，订单会停在原状态等后台人工退款
	if err := s.Refund(ctx, orderID, "打印失败自动退款", 0); err != nil {
		facades.Log().Errorf("打印失败自动退款失败 orderID=%d err=%v", orderID, err)
	}

	return nil
}

// CloseExpired 关闭超时未支付的订单：订单置「超期关闭」，未开始的明细置「已取消」。
// 支付单本身由 pay 模块的关单任务处理，两边都以 pay_expire_time 为准，互不依赖。
func (s *OrderService) CloseExpired(ctx context.Context, limit int) (int, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}

	list, err := s.orders.FindExpiredUnpaid(ctx, limit)
	if err != nil {
		return 0, err
	}

	closed := 0
	now := types.Now()
	for i := range list {
		order := &list[i]

		// 条件更新：只有仍是待支付的订单才会被关，回调先到不会被误关
		rows, err := facades.DB().WithContext(ctx).Table(model.PrintOrderTable).
			Where("id", order.ID).
			Where("order_status", printEnums.OrderStatusUnpaid).
			Update(map[string]any{
				"order_status":  printEnums.OrderStatusExpired,
				"close_time":    now,
				"cancel_reason": "支付超时自动关闭",
				"updated_at":    now,
			})
		if err != nil {
			return closed, err
		}
		if rows.RowsAffected == 0 {
			continue
		}

		if err := s.items.CancelByOrder(ctx, order.ID); err != nil {
			return closed, err
		}
		closed++
	}

	return closed, nil
}

// generatePickupCode 生成门店内未占用的 6 位取件码。
func (s *OrderService) generatePickupCode(ctx context.Context, shopID types.ID) string {
	for i := 0; i < 20; i++ {
		code := fmt.Sprintf("%06d", rand.Intn(1000000))
		exists, err := s.orders.ExistsPickupCode(ctx, shopID, code)
		if err == nil && !exists {
			return code
		}
	}

	// 极端情况下退化为时间戳后 6 位，仍保证门店内可查
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

// orderTotals 是下单汇总，用于写订单聚合字段。
type orderTotals struct {
	TotalPages  int
	TotalSheets int
	TotalAmount int64
	DeviceID    types.ID
}

// buildPrintConfig 生成打印参数全量快照，服务端下发、客户端只读。
func buildPrintConfig(paper, color, mode, orientation, pageRange string, copies int) *string {
	config := map[string]any{
		"paperType":   paper,
		"colorType":   color,
		"printMode":   mode,
		"orientation": orientation,
		"pageRange":   pageRange,
		"copies":      copies,
	}

	bs, err := json.Marshal(config)
	if err != nil {
		return nil
	}
	value := string(bs)

	return &value
}

// decodeConfig 解析打印参数快照，解析失败返回 nil（客户端按明细字段打印）。
func decodeConfig(raw *string) map[string]any {
	if raw == nil || *raw == "" {
		return nil
	}

	result := make(map[string]any)
	if err := json.Unmarshal([]byte(*raw), &result); err != nil {
		return nil
	}

	return result
}

// countPages 按页码区间计算计费页数；区间为空表示全部页。
func countPages(pageRange string, total int) int {
	pageRange = strings.TrimSpace(pageRange)
	if pageRange == "" || total <= 0 {
		return total
	}

	// 用集合去重：区间重叠时（如 1-5,3-8）同一页只算一次，否则会多收钱
	picked := make(map[int]struct{}, total)
	for _, part := range strings.Split(pageRange, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		start, end := parseRange(part)
		if start < 1 || end < start || start > total {
			continue
		}
		if end > total {
			end = total
		}
		for page := start; page <= end; page++ {
			picked[page] = struct{}{}
		}
	}

	return len(picked)
}

// parseRange 解析「5」或「1-5」形式的区间。
func parseRange(part string) (int, int) {
	if strings.Contains(part, "-") {
		var start, end int
		if _, err := fmt.Sscanf(part, "%d-%d", &start, &end); err != nil {
			return 0, 0
		}

		return start, end
	}

	var page int
	if _, err := fmt.Sscanf(part, "%d", &page); err != nil {
		return 0, 0
	}

	return page, page
}

// payTypeOf 把渠道配置编码映射成订单上的支付方式，便于统计。
func payTypeOf(channelCode string) string {
	switch channelCode {
	case payEnums.ChannelBalance:
		return printEnums.PayTypeWallet
	case payEnums.CodeAlipayWap:
		return printEnums.PayTypeAlipay
	default:
		return printEnums.PayTypeWechat
	}
}
