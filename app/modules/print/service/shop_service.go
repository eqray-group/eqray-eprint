package service

import (
	"context"
	"strings"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

// ShopInput 是门店新增 / 编辑入参。
type ShopInput struct {
	Name          string
	Address       string
	Longitude     float64
	Latitude      float64
	Avatar        string
	BusinessHours string
	Status        int8
	Remark        string
}

// ShopListQuery 是门店列表查询条件。
type ShopListQuery struct {
	Name        string
	Status      *int8
	AuditStatus *int8
	PageNo      int
	PageSize    int
}

// MyShopOption 是客户端登录后可选的归属门店，只暴露选择门店需要的字段。
type MyShopOption struct {
	ID     types.ID `json:"id"`
	Name   string   `json:"name"`
	Role   string   `json:"role"`
	Status int8     `json:"status"`
}

// ShopService 负责门店、门店成员与审核。
type ShopService struct {
	shops  *repository.PrintShopRepository
	staffs *repository.PrintShopStaffRepository
}

// NewShopService 构造门店服务。
func NewShopService() *ShopService {
	return &ShopService{
		shops:  repository.NewPrintShopRepository(),
		staffs: repository.NewPrintShopStaffRepository(),
	}
}

// Page 分页查询门店。
func (s *ShopService) Page(ctx context.Context, q ShopListQuery) ([]model.PrintShopModel, int64, error) {
	return s.shops.Paginate(ctx, repository.ShopFilter{
		Name:        strings.TrimSpace(q.Name),
		Status:      q.Status,
		AuditStatus: q.AuditStatus,
	}, q.PageNo, q.PageSize)
}

// Detail 查询门店详情。
func (s *ShopService) Detail(ctx context.Context, id types.ID) (*model.PrintShopModel, error) {
	shop, err := s.shops.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if shop == nil {
		return nil, baseErrs.New(printErrs.ShopNotExists, id)
	}

	return shop, nil
}

// Create 新增门店，审核状态默认为待审核，需平台审核通过后才可接单。
func (s *ShopService) Create(ctx context.Context, in ShopInput) (types.ID, error) {
	if err := validateShopInput(in); err != nil {
		return 0, err
	}

	shop := &model.PrintShopModel{
		Name:          strings.TrimSpace(in.Name),
		Address:       strings.TrimSpace(in.Address),
		Longitude:     in.Longitude,
		Latitude:      in.Latitude,
		Avatar:        strings.TrimSpace(in.Avatar),
		BusinessHours: normalizeHours(in.BusinessHours),
		Status:        in.Status,
		AuditStatus:   printEnums.AuditStatusPending,
		Remark:        strings.TrimSpace(in.Remark),
	}
	shop.ID = support.NewID()

	if err := s.shops.Create(ctx, shop); err != nil {
		return 0, err
	}

	return shop.ID, nil
}

// Update 编辑门店基础信息；审核状态不走这里。
func (s *ShopService) Update(ctx context.Context, id types.ID, in ShopInput) error {
	if err := validateShopInput(in); err != nil {
		return err
	}

	shop, err := s.shopOrErr(ctx, id)
	if err != nil {
		return err
	}

	shop.Name = strings.TrimSpace(in.Name)
	shop.Address = strings.TrimSpace(in.Address)
	shop.Longitude = in.Longitude
	shop.Latitude = in.Latitude
	shop.Avatar = strings.TrimSpace(in.Avatar)
	shop.BusinessHours = normalizeHours(in.BusinessHours)
	shop.Status = in.Status
	shop.Remark = strings.TrimSpace(in.Remark)

	return s.shops.Update(ctx, shop)
}

// ChangeStatus 切换营业状态。
func (s *ShopService) ChangeStatus(ctx context.Context, id types.ID, status int8, updater string) error {
	if status != printEnums.ShopStatusOpen && status != printEnums.ShopStatusClosed {
		return baseErrs.WithMsg(baseErrs.BadRequest, "营业状态取值不合法")
	}
	if _, err := s.shopOrErr(ctx, id); err != nil {
		return err
	}

	return s.shops.UpdateStatus(ctx, id, status, updater)
}

// Audit 审核门店：只有待审核或已拒绝的门店可以再次审核。
func (s *ShopService) Audit(ctx context.Context, id types.ID, auditStatus int8, remark, updater string) error {
	if auditStatus != printEnums.AuditStatusApproved && auditStatus != printEnums.AuditStatusRejected {
		return baseErrs.WithMsg(baseErrs.BadRequest, "审核结果取值不合法")
	}

	shop, err := s.shopOrErr(ctx, id)
	if err != nil {
		return err
	}
	if shop.AuditStatus == printEnums.AuditStatusApproved {
		return baseErrs.New(printErrs.ShopAuditInvalid)
	}

	return s.shops.UpdateAudit(ctx, id, auditStatus, strings.TrimSpace(remark), updater)
}

// Delete 删除门店；有未完成订单时不允许删除。
func (s *ShopService) Delete(ctx context.Context, id types.ID, updater string) error {
	if _, err := s.shopOrErr(ctx, id); err != nil {
		return err
	}

	return s.shops.Delete(ctx, id, updater)
}

// ListOpen 返回会员端可下单的门店（审核通过且营业中）。
func (s *ShopService) ListOpen(ctx context.Context) ([]model.PrintShopModel, error) {
	return s.shops.ListOpen(ctx)
}

// BindStaff 把后台账号绑定到门店；owner 只能有一个由后台控制，重复绑定直接报错。
func (s *ShopService) BindStaff(ctx context.Context, userID, shopID types.ID, role string, operator types.ID) error {
	if role != printEnums.StaffRoleOwner && role != printEnums.StaffRoleStaff {
		return baseErrs.WithMsg(baseErrs.BadRequest, "门店角色只能是 owner 或 staff")
	}
	if _, err := s.shopOrErr(ctx, shopID); err != nil {
		return err
	}

	exists, err := s.staffs.Exists(ctx, userID, shopID)
	if err != nil {
		return err
	}
	if exists {
		return baseErrs.New(printErrs.StaffExists)
	}

	staff := &model.PrintShopStaffModel{
		UserID: userID,
		ShopID: shopID,
		Role:   role,
		Status: printEnums.StatusEnabled,
	}
	staff.ID = support.NewID()

	return s.staffs.Create(ctx, staff)
}

// UnbindStaff 解除成员绑定（物理删除）。
func (s *ShopService) UnbindStaff(ctx context.Context, staffID types.ID) error {
	staff, err := s.staffs.FindByID(ctx, staffID)
	if err != nil {
		return err
	}
	if staff == nil {
		return baseErrs.New(printErrs.StaffNotExists)
	}

	return s.staffs.Delete(ctx, staffID)
}

// ListStaff 查询门店成员。
func (s *ShopService) ListStaff(ctx context.Context, shopID types.ID) ([]model.PrintShopStaffModel, error) {
	return s.staffs.ListByShop(ctx, shopID)
}

// MyShops 查询当前后台账号绑定的门店，客户端登录后据此选门店。
func (s *ShopService) MyShops(ctx context.Context, userID types.ID) ([]model.PrintShopStaffModel, error) {
	return s.staffs.ListByUser(ctx, userID)
}

// MyShopOptions 返回账号绑定且未删除的门店，供客户端登录后选择设备归属。
// 参数 userID 为后台账号；返回空切片表示未绑定任何门店。
func (s *ShopService) MyShopOptions(ctx context.Context, userID types.ID) ([]MyShopOption, error) {
	staffs, err := s.staffs.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	options := make([]MyShopOption, 0, len(staffs))
	for _, staff := range staffs {
		shop, err := s.shops.FindByID(ctx, staff.ShopID)
		if err != nil {
			return nil, err
		}
		if shop == nil {
			continue
		}

		options = append(options, MyShopOption{
			ID:     shop.ID,
			Name:   shop.Name,
			Role:   staff.Role,
			Status: shop.Status,
		})
	}

	return options, nil
}

// EnsureStaff 校验账号是否属于门店，设备登记与核销前调用。
func (s *ShopService) EnsureStaff(ctx context.Context, userID, shopID types.ID) error {
	staff, err := s.staffs.FindByUserShop(ctx, userID, shopID)
	if err != nil {
		return err
	}
	if staff == nil {
		return baseErrs.New(printErrs.StaffNotBelong)
	}

	return nil
}

// shopOrErr 取门店，不存在转业务错误。
func (s *ShopService) shopOrErr(ctx context.Context, id types.ID) (*model.PrintShopModel, error) {
	shop, err := s.shops.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if shop == nil {
		return nil, baseErrs.New(printErrs.ShopNotExists, id)
	}

	return shop, nil
}

// validateShopInput 校验必填与坐标范围。
func validateShopInput(in ShopInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return baseErrs.WithMsg(baseErrs.BadRequest, "门店名称不能为空")
	}
	if in.Status != printEnums.ShopStatusOpen && in.Status != printEnums.ShopStatusClosed {
		return baseErrs.WithMsg(baseErrs.BadRequest, "营业状态取值不合法")
	}

	return nil
}

// normalizeHours 把空营业时间转成 nil（表示全天营业），并校验非空时必须是 JSON。
func normalizeHours(hours string) *string {
	value := strings.TrimSpace(hours)
	if value == "" {
		return nil
	}

	return &value
}
