package service

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/model"
)

// StatisticsOverview 是控制台的打印业务概览，金额单位为分。
type StatisticsOverview struct {
	// TodayOrders 今日下单数（含未支付与已取消）
	TodayOrders int64 `json:"todayOrders"`
	// TodayPaidOrders 今日支付成功笔数
	TodayPaidOrders int64 `json:"todayPaidOrders"`
	// TodayPaidAmount 今日支付金额
	TodayPaidAmount int64 `json:"todayPaidAmount"`
	// PendingPickup 待取件订单数
	PendingPickup int64 `json:"pendingPickup"`
	// PrintingOrders 待打印 + 打印中的订单数
	PrintingOrders int64 `json:"printingOrders"`
	// FailedTasks 失败任务数（全量，便于排查存量问题）
	FailedTasks int64 `json:"failedTasks"`
	// OnlineDevices 在线且启用的设备数
	OnlineDevices int64 `json:"onlineDevices"`
}

// StatisticsService 汇总打印业务指标，供控制台卡片展示。
type StatisticsService struct{}

// NewStatisticsService 构造统计服务。
func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

// Overview 统计今日订单与支付、待取件、打印中、失败任务与在线设备。
// 无参数；单条 SQL 失败直接返回错误，避免给控制台展示半套数字。
func (s *StatisticsService) Overview(ctx context.Context) (*StatisticsOverview, error) {
	today := startOfToday()
	result := &StatisticsOverview{}

	var err error
	if result.TodayOrders, err = s.countOrders(ctx, func(q db.Query) db.Query {
		return q.Where("created_at", ">=", today)
	}); err != nil {
		return nil, err
	}

	paidFilter := func(q db.Query) db.Query {
		return q.Where("pay_status", printEnums.PayStatusPaid).Where("pay_time", ">=", today)
	}
	if result.TodayPaidOrders, err = s.countOrders(ctx, paidFilter); err != nil {
		return nil, err
	}
	if result.TodayPaidAmount, err = s.sumPaidAmount(ctx, today); err != nil {
		return nil, err
	}

	if result.PendingPickup, err = s.countOrders(ctx, func(q db.Query) db.Query {
		return q.Where("order_status", printEnums.OrderStatusPickup)
	}); err != nil {
		return nil, err
	}

	if result.PrintingOrders, err = s.countOrders(ctx, func(q db.Query) db.Query {
		return q.WhereIn("order_status", []any{printEnums.OrderStatusPending, printEnums.OrderStatusPrinting})
	}); err != nil {
		return nil, err
	}

	if result.FailedTasks, err = facades.DB().WithContext(ctx).Table(model.PrintOrderItemTable).
		Where("status", printEnums.TaskStatusFailed).
		WhereNull("deleted_at").
		Count(); err != nil {
		return nil, err
	}

	if result.OnlineDevices, err = facades.DB().WithContext(ctx).Table(model.PrintDeviceTable).
		Where("online_status", printEnums.DeviceOnline).
		Where("status", printEnums.StatusEnabled).
		WhereNull("deleted_at").
		Count(); err != nil {
		return nil, err
	}

	return result, nil
}

// countOrders 按附加条件统计订单数。
func (s *StatisticsService) countOrders(ctx context.Context, apply func(db.Query) db.Query) (int64, error) {
	query := apply(facades.DB().WithContext(ctx).Table(model.PrintOrderTable).WhereNull("deleted_at"))

	return query.Count()
}

// paidAmountRow 承接 SUM(pay_amount) 的聚合结果。
type paidAmountRow struct {
	Total int64 `db:"total" json:"total"`
}

// sumPaidAmount 统计今日支付金额；当天没有支付记录时 SUM 返回 NULL，
// 直接扫进 int64 会报 Scan error，这里用 COALESCE 兜底成 0。
func (s *StatisticsService) sumPaidAmount(ctx context.Context, today time.Time) (int64, error) {
	// 聚合结果只有一行，用 First 按结构体扫描
	var row paidAmountRow
	err := facades.DB().WithContext(ctx).Table(model.PrintOrderTable).
		Where("pay_status", printEnums.PayStatusPaid).
		Where("pay_time", ">=", today).
		WhereNull("deleted_at").
		Select("COALESCE(SUM(pay_amount), 0) AS total").
		First(&row)
	if err != nil {
		return 0, err
	}

	return row.Total, nil
}

// startOfToday 返回本地时区当天 00:00:00，统计口径与业务时区一致。
func startOfToday() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
