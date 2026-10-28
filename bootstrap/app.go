package bootstrap

import (
	"context"

	contractsFoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"

	baseMiddleware "eqray-eprint/app/modules/base/middleware"
	"eqray-eprint/app/modules/base/types"
	infraMiddleware "eqray-eprint/app/modules/infra/middleware"
	payService "eqray-eprint/app/modules/pay/service"
	printService "eqray-eprint/app/modules/print/service"
	systemService "eqray-eprint/app/modules/system/service"
	"eqray-eprint/config"
	"eqray-eprint/routes"
)

func Boot() contractsFoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithRouting(func() {
			// 权限数据归 system，中间件在 base；启动时注入实现，保持模块单向依赖
			baseMiddleware.RegisterPermissionChecker(
				func(ctx context.Context, userID types.ID, code string) (bool, error) {
					return systemService.NewPermissionService().Allows(ctx, userID, code)
				},
			)

			// 操作日志中间件在 base，写库实现归 system，启动时注入
			baseMiddleware.RegisterOperateLogRecorder(
				func(ctx context.Context, in baseMiddleware.OperateLogInput) error {
					return systemService.NewOperateLogService().Record(ctx, systemService.OperateLogInput{
						UserID:        in.UserID,
						UserType:      in.UserType,
						Type:          in.Type,
						SubType:       in.SubType,
						Action:        in.Action,
						Success:       in.Success,
						Extra:         in.Extra,
						RequestMethod: in.RequestMethod,
						RequestURL:    in.RequestURL,
						RequestIP:     in.RequestIP,
						UserAgent:     in.UserAgent,
						Duration:      in.Duration,
					})
				},
			)

			// 支付成功回调：pay 不感知业务，由 bootstrap 注入打印订单的派发实现
			payService.SetPaidHook(func(ctx context.Context, bizType string, bizID types.ID) error {
				return printService.NewOrderService().OnPaid(ctx, bizType, bizID)
			})
			// 退款成功回调：同一份注入思路，打印订单据此落「已退款」
			payService.SetRefundedHook(func(ctx context.Context, bizType string, bizID types.ID) error {
				return printService.NewOrderService().OnRefunded(ctx, bizType, bizID)
			})

			routes.Web()
			routes.Grpc()
			routes.Modules()
		}).
		WithProviders(Providers).
		WithSchedule(Schedule).
		WithConfig(config.Boot).
		// API 访问 / 错误日志：全局中间件，覆盖所有模块的请求
		WithMiddleware(func(handler configuration.Middleware) {
			handler.Append(&infraMiddleware.RequestLog{})
		}).
		Create()
}
