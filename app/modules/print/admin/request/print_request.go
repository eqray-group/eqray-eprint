// Package request 定义 print 模块后台接口的入参结构。
package request

import (
	"github.com/goravel/framework/contracts/http"
)

// ShopSaveRequest 是门店新增 / 编辑入参。
type ShopSaveRequest struct {
	Name          string  `form:"name" json:"name"`
	Address       string  `form:"address" json:"address"`
	Longitude     float64 `form:"longitude" json:"longitude"`
	Latitude      float64 `form:"latitude" json:"latitude"`
	Avatar        string  `form:"avatar" json:"avatar"`
	BusinessHours string  `form:"businessHours" json:"businessHours"`
	Status        *int8   `form:"status" json:"status"`
	Remark        string  `form:"remark" json:"remark"`
}

// Rules 返回校验规则。
func (r *ShopSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"name":          "required|max:64",
		"address":       "max:255",
		"businessHours": "max:1000",
		"status":        "required|integer|in:0,1",
		"remark":        "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *ShopSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":   "门店名称不能为空",
		"status.required": "请选择营业状态",
		"status.in":       "营业状态取值不合法",
	}
}

// ShopStatusRequest 是营业状态开关入参。
type ShopStatusRequest struct {
	Status int8 `form:"status" json:"status"`
}

// Rules 返回校验规则。
func (r *ShopStatusRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{"status": "required|integer|in:0,1"}
}

// Messages 覆盖默认提示。
func (r *ShopStatusRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{"status.in": "营业状态取值不合法"}
}

// ShopAuditRequest 是门店审核入参：1 通过 / 2 拒绝。
type ShopAuditRequest struct {
	AuditStatus int8   `form:"auditStatus" json:"auditStatus"`
	Remark      string `form:"remark" json:"remark"`
}

// Rules 返回校验规则。
func (r *ShopAuditRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"auditStatus": "required|integer|in:1,2",
		"remark":      "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *ShopAuditRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{"auditStatus.in": "审核结果只能是 1 通过或 2 拒绝"}
}

// StaffSaveRequest 是门店成员绑定入参。
type StaffSaveRequest struct {
	UserID int64  `form:"userId" json:"userId"`
	Role   string `form:"role" json:"role"`
}

// Rules 返回校验规则。
func (r *StaffSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"userId": "required|integer|min:1",
		"role":   "required|in:owner,staff",
	}
}

// Messages 覆盖默认提示。
func (r *StaffSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"userId.required": "请选择要绑定的后台账号",
		"role.in":         "门店角色只能是 owner 或 staff",
	}
}

// DeviceSaveRequest 是设备新增 / 编辑入参。
type DeviceSaveRequest struct {
	Code        string `form:"code" json:"code"`
	PrintShopID int64  `form:"printShopId" json:"printShopId"`
	Name        string `form:"name" json:"name"`
	Status      *int8  `form:"status" json:"status"`
	Remark      string `form:"remark" json:"remark"`
}

// Rules 返回校验规则。
func (r *DeviceSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"printShopId": "required|integer|min:1",
		"name":        "max:64",
		"status":      "required|integer|in:0,1",
		"remark":      "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *DeviceSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"printShopId.required": "请选择所属门店",
		"status.in":            "设备状态取值不合法",
	}
}

// PrinterSaveRequest 是打印机新增 / 编辑入参，能力字段由后台勾选。
type PrinterSaveRequest struct {
	DeviceID      int64  `form:"deviceId" json:"deviceId"`
	Name          string `form:"name" json:"name"`
	SystemName    string `form:"systemName" json:"systemName"`
	Model         string `form:"model" json:"model"`
	PaperTypes    string `form:"paperTypes" json:"paperTypes"`
	ColorSupport  int8   `form:"colorSupport" json:"colorSupport"`
	DuplexSupport int8   `form:"duplexSupport" json:"duplexSupport"`
	Status        *int8  `form:"status" json:"status"`
	Remark        string `form:"remark" json:"remark"`
}

// Rules 返回校验规则。
func (r *PrinterSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"deviceId":      "required|integer|min:1",
		"name":          "required|max:64",
		"systemName":    "required|max:128",
		"model":         "max:64",
		"paperTypes":    "required|max:32",
		"colorSupport":  "required|integer|in:0,1",
		"duplexSupport": "required|integer|in:0,1",
		"status":        "required|integer|in:0,1",
		"remark":        "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *PrinterSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"deviceId.required":   "请选择所属设备",
		"systemName.required": "请填写 Windows 打印队列名称",
		"paperTypes.required": "请至少选择一种纸张",
	}
}

// PriceSaveRequest 是价格规则新增入参，价格单位为分/页。
type PriceSaveRequest struct {
	PrintShopID int64  `form:"printShopId" json:"printShopId"`
	PaperType   string `form:"paperType" json:"paperType"`
	ColorType   string `form:"colorType" json:"colorType"`
	PrintMode   string `form:"printMode" json:"printMode"`
	Price       int64  `form:"price" json:"price"`
	Status      *int8  `form:"status" json:"status"`
}

// Rules 返回校验规则。
func (r *PriceSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"printShopId": "required|integer|min:1",
		"paperType":   "required|in:A4,A3",
		"colorType":   "required|in:bw,color",
		"printMode":   "required|in:single,double",
		"price":       "required|integer|min:1",
		"status":      "required|integer|in:0,1",
	}
}

// Messages 覆盖默认提示。
func (r *PriceSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"price.min":    "单价必须大于 0",
		"paperType.in": "纸张只能是 A4 或 A3",
		"colorType.in": "色彩只能是 bw 或 color",
		"printMode.in": "单双面只能是 single 或 double",
	}
}

// PickupVerifyRequest 是取件码核销入参。
type PickupVerifyRequest struct {
	PrintShopID int64  `form:"printShopId" json:"printShopId"`
	Code        string `form:"code" json:"code"`
}

// Rules 返回校验规则。
func (r *PickupVerifyRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"printShopId": "required|integer|min:1",
		"code":        "required|max:16",
	}
}

// Messages 覆盖默认提示。
func (r *PickupVerifyRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"printShopId.required": "请选择门店",
		"code.required":        "请输入取件码",
	}
}

// BannerSaveRequest 是轮播图新增 / 编辑入参；printShopId 传 0 表示全平台。
type BannerSaveRequest struct {
	PrintShopID int64  `form:"printShopId" json:"printShopId"`
	Title       string `form:"title" json:"title"`
	ImageURL    string `form:"imageUrl" json:"imageUrl"`
	LinkType    string `form:"linkType" json:"linkType"`
	LinkURL     string `form:"linkUrl" json:"linkUrl"`
	Sort        int    `form:"sort" json:"sort"`
	Status      *int8  `form:"status" json:"status"`
}

// Rules 返回校验规则。
func (r *BannerSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"printShopId": "required|integer|min:0",
		"title":       "max:64",
		"imageUrl":    "required|max:255",
		"linkType":    "in:none,page,url",
		"linkUrl":     "max:255",
		"sort":        "integer|min:0",
		"status":      "required|integer|in:0,1",
	}
}

// Messages 覆盖默认提示。
func (r *BannerSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"imageUrl.required": "请上传轮播图图片",
		"linkType.in":       "跳转类型只能是 none / page / url",
	}
}

// OrderRefundRequest 是订单退款入参：只支持全额退款，金额由服务端取订单实付。
type OrderRefundRequest struct {
	Reason string `form:"reason" json:"reason"`
}

// Rules 返回校验规则。
func (r *OrderRefundRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{"reason": "required|max:255"}
}

// Messages 覆盖默认提示。
func (r *OrderRefundRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{"reason.required": "请填写退款原因"}
}

// PriceUpdateRequest 是价格修改入参：只允许改价格与状态。
type PriceUpdateRequest struct {
	Price  int64 `form:"price" json:"price"`
	Status int8  `form:"status" json:"status"`
}

// Rules 返回校验规则。
func (r *PriceUpdateRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"price":  "required|integer|min:1",
		"status": "required|integer|in:0,1",
	}
}

// Messages 覆盖默认提示。
func (r *PriceUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{"price.min": "单价必须大于 0"}
}
