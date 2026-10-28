// Package request 定义打印客户端（Windows 端）接口入参。
package request

import (
	"github.com/goravel/framework/contracts/http"
)

// RegisterRequest 是客户端登录后登记设备入参，code 为机器指纹。
type RegisterRequest struct {
	PrintShopID int64  `form:"printShopId" json:"printShopId"`
	Code        string `form:"code" json:"code"`
	Name        string `form:"name" json:"name"`
	Version     string `form:"version" json:"version"`
}

// Rules 返回校验规则。
func (r *RegisterRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"printShopId": "required|integer|min:1",
		"code":        "required|max:128",
		"name":        "max:64",
		"version":     "max:32",
	}
}

// Messages 覆盖默认提示。
func (r *RegisterRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"printShopId.required": "请选择所属门店",
		"code.required":        "设备编码不能为空",
	}
}

// HeartbeatRequest 是心跳入参。
type HeartbeatRequest struct {
	Version string `form:"version" json:"version"`
}

// Rules 返回校验规则。
func (r *HeartbeatRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{"version": "max:32"}
}

// SyncPrintersRequest 是客户端上报本机打印机的入参。
// printers 为当前勾选的完整集合，服务端据此启用 / 停用设备下的打印机。
type SyncPrintersRequest struct {
	Printers []PrinterReportRequest `form:"printers" json:"printers"`
}

// PrinterReportRequest 是单台打印机的上报内容，能力字段由客户端界面勾选。
type PrinterReportRequest struct {
	SystemName    string `form:"systemName" json:"systemName"`
	Name          string `form:"name" json:"name"`
	Model         string `form:"model" json:"model"`
	PaperTypes    string `form:"paperTypes" json:"paperTypes"`
	ColorSupport  int8   `form:"colorSupport" json:"colorSupport"`
	DuplexSupport int8   `form:"duplexSupport" json:"duplexSupport"`
}

// Rules 返回校验规则；数组内部字段由服务层逐条校验，此处只校验外层类型。
func (r *SyncPrintersRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{}
}

// ReportRequest 是任务回执入参：2 打印中 / 3 完成 / 4 失败。
type ReportRequest struct {
	TaskNo       string `form:"taskNo" json:"taskNo"`
	Status       int8   `form:"status" json:"status"`
	PrintedPages int    `form:"printedPages" json:"printedPages"`
	ErrorCode    string `form:"errorCode" json:"errorCode"`
	Reason       string `form:"reason" json:"reason"`
}

// Rules 返回校验规则。
func (r *ReportRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"taskNo":    "required|max:32",
		"status":    "required|integer|in:2,3,4",
		"errorCode": "max:64",
		"reason":    "max:255",
	}
}

// Messages 覆盖默认提示。
func (r *ReportRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"taskNo.required": "任务号不能为空",
		"status.in":       "回执状态只能是 2 打印中、3 完成、4 失败",
	}
}

// PickupRequest 是取件码核销入参。
type PickupRequest struct {
	Code string `form:"code" json:"code"`
}

// Rules 返回校验规则。
func (r *PickupRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{"code": "required|max:16"}
}

// Messages 覆盖默认提示。
func (r *PickupRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{"code.required": "请输入取件码"}
}
