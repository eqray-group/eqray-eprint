// Package controller 提供打印客户端使用的设备、任务与核销接口。
// 客户端用后台账号登录，设备归属门店取自该账号在 print_shop_staff 的绑定。
package controller

import (
	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/device/request"
	"eqray-eprint/app/modules/print/service"
)

// DeviceController 处理客户端登记、心跳、拉任务、回执与核销。
type DeviceController struct {
	devices *service.DeviceService
	orders  *service.OrderService
}

// NewDeviceController 构建设备控制器。
func NewDeviceController() *DeviceController {
	return &DeviceController{
		devices: service.NewDeviceService(),
		orders:  service.NewOrderService(),
	}
}

// Register 客户端登录后登记设备，已存在则更新归属并刷新心跳。
func (c *DeviceController) Register(ctx http.Context) http.Response {
	var req request.RegisterRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	device, err := c.devices.Register(
		ctx.Context(),
		types.ID(support.UserID(ctx.Context())),
		types.ID(req.PrintShopID),
		req.Code,
		req.Name,
		req.Version,
		clientIP(ctx),
	)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, device)
}

// Heartbeat 刷新心跳；客户端空转时每 30 秒调用一次。
func (c *DeviceController) Heartbeat(ctx http.Context) http.Response {
	deviceID, err := pathDeviceID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.HeartbeatRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.devices.Heartbeat(
		ctx.Context(),
		deviceID,
		types.ID(support.UserID(ctx.Context())),
		req.Version,
		clientIP(ctx),
	); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Tasks 拉取本机待打印任务，同时顺带刷新心跳，避免客户端多打一次请求。
func (c *DeviceController) Tasks(ctx http.Context) http.Response {
	deviceID, err := pathDeviceID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	// 心跳失败不阻断拉任务：网络抖动时客户端仍可继续打印已拉到的内容
	_ = c.devices.Heartbeat(ctx.Context(), deviceID, types.ID(support.UserID(ctx.Context())), ctx.Request().Query("version"), clientIP(ctx))

	tasks, err := c.orders.ClaimTasks(ctx.Context(), deviceID, ctx.Request().QueryInt("limit", 0))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, tasks)
}

// Report 任务回执：打印中 / 完成 / 失败。
func (c *DeviceController) Report(ctx http.Context) http.Response {
	deviceID, err := pathDeviceID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.ReportRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.orders.ReportTask(ctx.Context(), service.ReportTaskInput{
		DeviceID:     deviceID,
		TaskNo:       req.TaskNo,
		Status:       req.Status,
		PrintedPages: req.PrintedPages,
		ErrorCode:    req.ErrorCode,
		Reason:       req.Reason,
	}); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Pickup 门店核销取件码（客户端与店端共用）。
func (c *DeviceController) Pickup(ctx http.Context) http.Response {
	shopID, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.PickupRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	// 只有门店成员可以核销别的门店的订单
	if err := service.NewShopService().EnsureStaff(ctx.Context(), types.ID(support.UserID(ctx.Context())), shopID); err != nil {
		return result.Fail(ctx, err)
	}

	order, err := c.orders.PickupVerify(ctx.Context(), shopID, req.Code, types.ID(support.UserID(ctx.Context())))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, order)
}

// pathDeviceID 解析路径里的设备 ID。
func pathDeviceID(ctx http.Context) (types.ID, error) {
	value, err := param.PathID(ctx)
	if err != nil || value <= 0 {
		return 0, baseErrs.WithMsg(baseErrs.BadRequest, "设备 ID 不合法")
	}

	return value, nil
}

// clientIP 取客户端 IP，优先 X-Forwarded-For 的第一段，便于排查网络问题。
func clientIP(ctx http.Context) string {
	origin := ctx.Request().Origin()
	if forwarded := origin.Header.Get("X-Forwarded-For"); forwarded != "" {
		for i := 0; i < len(forwarded); i++ {
			if forwarded[i] == ',' {
				return forwarded[:i]
			}
		}

		return forwarded
	}

	return origin.RemoteAddr
}
