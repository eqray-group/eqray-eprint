package service

import (
	"context"
	"strings"
	"time"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

// DeviceInput 是设备新增 / 编辑入参。
type DeviceInput struct {
	PrintShopID types.ID
	Name        string
	Status      int8
	Remark      string
}

// DeviceListQuery 是设备列表查询条件。
type DeviceListQuery struct {
	ShopID       types.ID
	Keyword      string
	Status       *int8
	OnlineStatus *int8
	PageNo       int
	PageSize     int
}

// DeviceService 负责打印设备的登记、心跳与离线判定。
type DeviceService struct {
	devices *repository.PrintDeviceRepository
	shops   *ShopService
}

// NewDeviceService 构造设备服务。
func NewDeviceService() *DeviceService {
	return &DeviceService{
		devices: repository.NewPrintDeviceRepository(),
		shops:   NewShopService(),
	}
}

// Page 分页查询设备。
func (s *DeviceService) Page(ctx context.Context, q DeviceListQuery) ([]model.PrintDeviceModel, int64, error) {
	return s.devices.Paginate(ctx, repository.DeviceFilter{
		ShopID:       q.ShopID,
		Keyword:      strings.TrimSpace(q.Keyword),
		Status:       q.Status,
		OnlineStatus: q.OnlineStatus,
	}, q.PageNo, q.PageSize)
}

// Detail 查询设备详情。
func (s *DeviceService) Detail(ctx context.Context, id types.ID) (*model.PrintDeviceModel, error) {
	device, err := s.devices.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, baseErrs.New(printErrs.DeviceNotExists, id)
	}

	return device, nil
}

// Register 客户端登录后登记设备，返回设备记录（已存在则更新归属与心跳）。
// 参数 userID 为登录账号，shopID 为选中的门店，code 为机器指纹，ip / version 为本次上报。
func (s *DeviceService) Register(ctx context.Context, userID, shopID types.ID, code, name, version, ip string) (*model.PrintDeviceModel, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "设备编码不能为空")
	}
	if err := s.shops.EnsureStaff(ctx, userID, shopID); err != nil {
		return nil, err
	}

	device, err := s.devices.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if device == nil {
		device = &model.PrintDeviceModel{
			PrintShopID:     shopID,
			Code:            code,
			Name:            defaultDeviceName(name, code),
			ClientVersion:   version,
			Status:          printEnums.StatusEnabled,
			OnlineStatus:    printEnums.DeviceOnline,
			IP:              ip,
			LastLoginUserID: userID,
		}
		device.ID = support.NewID()

		if err := s.devices.Create(ctx, device); err != nil {
			if support.IsDuplicate(err) {
				return nil, baseErrs.New(printErrs.DeviceCodeExists, code)
			}

			return nil, err
		}

		now := types.Now()
		device.LastHeartbeatAt = &now

		return device, nil
	}

	if device.Status != printEnums.StatusEnabled {
		return nil, baseErrs.New(printErrs.DeviceDisabled)
	}

	// 已存在的设备换门店或换账号登录时，更新归属并刷新心跳
	if err := s.devices.UpdateHeartbeat(ctx, device.ID, ip, version, userID); err != nil {
		return nil, err
	}
	if device.PrintShopID != shopID {
		device.PrintShopID = shopID
		if err := s.devices.Update(ctx, device); err != nil {
			return nil, err
		}
	}

	return s.devices.FindByID(ctx, device.ID)
}

// Heartbeat 刷新设备心跳（客户端空转或拉任务时调用）。
// 参数 deviceID 为设备 ID，其余为本次上报信息；设备不存在或禁用返回业务错误。
func (s *DeviceService) Heartbeat(ctx context.Context, deviceID types.ID, userID types.ID, version, ip string) error {
	device, err := s.devices.FindByID(ctx, deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return baseErrs.New(printErrs.DeviceNotExists, deviceID)
	}
	if device.Status != printEnums.StatusEnabled {
		return baseErrs.New(printErrs.DeviceDisabled)
	}

	return s.devices.UpdateHeartbeat(ctx, deviceID, ip, version, userID)
}

// Create 后台新增设备（不经过客户端登记，用于预录设备）。
func (s *DeviceService) Create(ctx context.Context, code string, in DeviceInput) (types.ID, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, baseErrs.WithMsg(baseErrs.BadRequest, "设备编码不能为空")
	}
	if _, err := s.shops.Detail(ctx, in.PrintShopID); err != nil {
		return 0, err
	}

	device := &model.PrintDeviceModel{
		PrintShopID: in.PrintShopID,
		Code:        code,
		Name:        defaultDeviceName(in.Name, code),
		Status:      in.Status,
		Remark:      strings.TrimSpace(in.Remark),
	}
	device.ID = support.NewID()

	if err := s.devices.Create(ctx, device); err != nil {
		if support.IsDuplicate(err) {
			return 0, baseErrs.New(printErrs.DeviceCodeExists, code)
		}

		return 0, err
	}

	return device.ID, nil
}

// Update 编辑设备名称、状态与归属门店。
func (s *DeviceService) Update(ctx context.Context, id types.ID, in DeviceInput) error {
	device, err := s.Detail(ctx, id)
	if err != nil {
		return err
	}

	device.PrintShopID = in.PrintShopID
	device.Name = defaultDeviceName(in.Name, device.Code)
	device.Status = in.Status
	device.Remark = strings.TrimSpace(in.Remark)

	return s.devices.Update(ctx, device)
}

// Delete 删除设备；其下打印机需先停用。
func (s *DeviceService) Delete(ctx context.Context, id types.ID, updater string) error {
	if _, err := s.Detail(ctx, id); err != nil {
		return err
	}

	return s.devices.Delete(ctx, id, updater)
}

// MarkOffline 把心跳超时的设备置为离线，供定时任务调用，返回处理数量。
// 参数 timeoutSeconds 为离线阈值（默认 90 秒）。
func (s *DeviceService) MarkOffline(ctx context.Context, timeoutSeconds int) (int64, error) {
	before := types.DateTime(types.Now().Time().Add(-time.Duration(timeoutSeconds) * time.Second))

	return s.devices.MarkOffline(ctx, before)
}

// ListByShop 查询门店启用设备。
func (s *DeviceService) ListByShop(ctx context.Context, shopID types.ID) ([]model.PrintDeviceModel, error) {
	return s.devices.ListByShop(ctx, shopID)
}

// defaultDeviceName 生成设备名称：未填时用「打印客户端-编码后 4 位」。
func defaultDeviceName(name, code string) string {
	value := strings.TrimSpace(name)
	if value != "" {
		return value
	}
	if len(code) > 4 {
		return "打印客户端-" + code[len(code)-4:]
	}

	return "打印客户端-" + code
}
