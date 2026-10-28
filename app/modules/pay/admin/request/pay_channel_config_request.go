package request

import (
	"github.com/goravel/framework/contracts/http"
)

// ChannelConfigSaveRequest 是后台新增 / 编辑支付渠道配置的入参。
// 三个密钥字段编辑时留空表示不修改，避免前端没回填就覆盖原值。
type ChannelConfigSaveRequest struct {
	Code              string `form:"code" json:"code"`
	Name              string `form:"name" json:"name"`
	Channel           string `form:"channel" json:"channel"`
	Scene             string `form:"scene" json:"scene"`
	AppID             string `form:"appId" json:"appId"`
	MchID             string `form:"mchId" json:"mchId"`
	SerialNo          string `form:"serialNo" json:"serialNo"`
	PrivateKey        string `form:"privateKey" json:"privateKey"`
	APIKey            string `form:"apiKey" json:"apiKey"`
	PlatformPublicKey string `form:"platformPublicKey" json:"platformPublicKey"`
	Gateway           string `form:"gateway" json:"gateway"`
	NotifyURL         string `form:"notifyUrl" json:"notifyUrl"`
	ReturnURL         string `form:"returnUrl" json:"returnUrl"`
	Extra             string `form:"extra" json:"extra"`
	Status            int8   `form:"status" json:"status"`
	Sort              int    `form:"sort" json:"sort"`
	Remark            string `form:"remark" json:"remark"`
}

// Rules 返回校验规则；渠道与场景的枚举关系由 service 校验。
func (r *ChannelConfigSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"code":      "required|max:32",
		"name":      "required|max:64",
		"channel":   "required|in:wechat,alipay",
		"scene":     "required|max:16",
		"appId":     "max:64",
		"mchId":     "max:64",
		"notifyUrl": "max:255",
		"returnUrl": "max:255",
		"remark":    "max:255",
		"status":    "integer|in:0,1",
		"sort":      "integer",
	}
}

// Messages 覆盖默认英文提示。
func (r *ChannelConfigSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"code.required":    "配置编码不能为空",
		"code.max":         "配置编码不能超过 32 个字符",
		"name.required":    "配置名称不能为空",
		"name.max":         "配置名称不能超过 64 个字符",
		"channel.required": "请选择渠道",
		"channel.in":       "渠道只能是 wechat 或 alipay",
		"scene.required":   "请选择支付场景",
		"appId.max":        "应用 ID 不能超过 64 个字符",
		"mchId.max":        "商户号不能超过 64 个字符",
		"notifyUrl.max":    "回调地址不能超过 255 个字符",
		"returnUrl.max":    "回跳地址不能超过 255 个字符",
		"remark.max":       "备注不能超过 255 个字符",
		"status.in":        "状态取值不合法",
		"status.integer":   "状态取值不合法",
		"sort.integer":     "排序必须是整数",
	}
}
