package bootstrap

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/schedule"

	"eqray-eprint/app/facades"
	infraService "eqray-eprint/app/modules/infra/service"
	payService "eqray-eprint/app/modules/pay/service"
	printService "eqray-eprint/app/modules/print/service"
)

// Schedule 注册应用级定时任务，由 foundation.Setup().WithSchedule 装载，
// 随 app.Start() 的 schedule runner 一起启动。
func Schedule() []schedule.Event {
	return []schedule.Event{
		// 访问日志保留策略：每天 03:10 清理 infra_api_access_log 中超过保留天数的记录，
		// 保留天数见 config/infra.go（默认 7 天，可用 INFRA_API_LOG_RETENTION_DAYS 覆盖）。
		// OnOneServer 保证多实例部署时同一时刻只有一台机器执行。
		facades.Schedule().Call(purgeExpiredApiAccessLogs).
			DailyAt("03:10").
			Name("infra:purge-api-access-log").
			OnOneServer(),

		// 支付关单：每 5 分钟关闭超时未支付的支付单（默认 30 分钟超时），
		// 先尽力关渠道侧订单再本地置关闭；批次上限 200 条，避免一次拉太多。
		facades.Schedule().Call(closeExpiredPayments).
			EveryFiveMinutes().
			Name("pay:close-expired-payments").
			OnOneServer(),

		// 打印订单关单：每 5 分钟关闭支付超时的订单并取消未开始的任务
		facades.Schedule().Call(closeExpiredPrintOrders).
			EveryFiveMinutes().
			Name("print:close-expired-orders").
			OnOneServer(),

		// 打印文件清理：每天 03:40 删除超过保留期（7 天）的文件与存储对象
		facades.Schedule().Call(cleanupExpiredPrintFiles).
			DailyAt("03:40").
			Name("print:cleanup-expired-files").
			OnOneServer(),

		// 设备离线检测：每分钟把心跳超过 90 秒的设备置为离线
		facades.Schedule().Call(markOfflinePrintDevices).
			EveryMinute().
			Name("print:mark-offline-devices").
			OnOneServer(),
	}
}

// closeExpiredPrintOrders 执行一次打印订单关单，结果写入应用日志。
// 无入参与返回值：异常只记录日志，避免打断后续调度任务。
func closeExpiredPrintOrders() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	closed, err := printService.NewOrderService().CloseExpired(ctx, 200)
	if err != nil {
		facades.Log().Errorf("打印订单关单任务失败：%v", err)
		return
	}

	if closed > 0 {
		facades.Log().Infof("打印订单关单任务完成，共关闭 %d 笔", closed)
	}
}

// cleanupExpiredPrintFiles 清理过期打印文件，结果写入应用日志。
// 无入参与返回值：异常只记录日志，避免打断后续调度任务。
func cleanupExpiredPrintFiles() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cleaned, err := printService.NewFileService().CleanupExpired(ctx, 200)
	if err != nil {
		facades.Log().Errorf("清理打印文件失败：%v", err)
		return
	}

	if cleaned > 0 {
		facades.Log().Infof("清理打印文件完成，共删除 %d 个", cleaned)
	}
}

// markOfflinePrintDevices 把心跳超时的设备置为离线。
// 阈值 90 秒 = 客户端 30 秒心跳 × 3 次容忍，避免网络抖动就判离线。
func markOfflinePrintDevices() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	affected, err := printService.NewDeviceService().MarkOffline(ctx, 90)
	if err != nil {
		facades.Log().Errorf("打印设备离线检测失败：%v", err)
		return
	}

	if affected > 0 {
		facades.Log().Infof("打印设备离线检测完成，共置离线 %d 台", affected)
	}
}

// closeExpiredPayments 执行一次支付关单，结果写入应用日志。
// 无入参与返回值：异常只记录日志，避免打断后续调度任务。
func closeExpiredPayments() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	closed, err := payService.NewPaymentService().CloseExpired(ctx, 200)
	if err != nil {
		facades.Log().Errorf("支付关单任务失败：%v", err)
		return
	}

	if closed > 0 {
		facades.Log().Infof("支付关单任务完成，共关闭 %d 笔", closed)
	}
}

// purgeExpiredApiAccessLogs 执行一次访问日志清理，并把结果写入应用日志。
// 无入参与返回值：异常只记录日志，避免把错误抛给调度器打断后续任务。
func purgeExpiredApiAccessLogs() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	deleted, err := infraService.NewApiLogService().PurgeExpiredAccess(ctx)
	if err != nil {
		facades.Log().Errorf("清理 API 访问日志失败：%v", err)
		return
	}

	if deleted > 0 {
		facades.Log().Infof("清理 API 访问日志完成，共删除 %d 条", deleted)
	}
}
