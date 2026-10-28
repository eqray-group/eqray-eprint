package request

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/types"
)

// MemberSaveRequest 是新增 / 编辑会员的入参。
// 会员不再维护积分、经验、标签、分组与等级，相关字段已随体系下线移除。
type MemberSaveRequest struct {
	Nickname string `form:"nickname" json:"nickname"`
	Mobile   string `form:"mobile" json:"mobile"`
	Avatar   string `form:"avatar" json:"avatar"`
	// Name 真实姓名
	Name     string          `form:"name" json:"name"`
	Sex      int8            `form:"sex" json:"sex"`
	Birthday *types.DateTime `form:"birthday" json:"birthday"`
	AreaID   int             `form:"areaId" json:"areaId"`
	Mark     string          `form:"mark" json:"mark"`
	// Status 会员状态：新增时不传默认启用；编辑时不传表示保持原状态
	Status *int8 `form:"status" json:"status"`
}

// Rules 返回校验规则，由 Controller 交给 param.Bind 执行。
func (r *MemberSaveRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"nickname": "required|max:64",
		"mobile":   "required|max:20",
		"avatar":   "max:512",
		"name":     "max:30",
		"sex":      "integer|in:0,1,2",
		"areaId":   "integer|min:0",
		"mark":     "max:255",
		// 不传 status 时跳过校验：新增默认启用，编辑保持原状态
		"status": "integer|in:0,1",
	}
}

// Messages 覆盖框架默认的英文提示，保证接口文案统一为中文。
func (r *MemberSaveRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nickname.required": "昵称不能为空",
		"nickname.max":      "昵称不能超过 64 个字符",
		"mobile.required":   "手机号不能为空",
		"mobile.max":        "手机号不能超过 20 个字符",
		"avatar.max":        "头像地址不能超过 512 个字符",
		"name.max":          "真实姓名不能超过 30 个字符",
		"sex.integer":       "性别取值不合法",
		"sex.in":            "性别取值不合法",
		"areaId.integer":    "所在地取值不合法",
		"areaId.min":        "所在地取值不合法",
		"mark.max":          "用户标记不能超过 255 个字符",
		"status.integer":    "会员状态取值不合法",
		"status.in":         "会员状态取值不合法",
	}
}

// MemberStatusRequest 是启用 / 禁用会员的入参。
type MemberStatusRequest struct {
	Status int8 `form:"status" json:"status"`
}

// Rules 返回校验规则。
func (r *MemberStatusRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"status": "required|integer|in:0,1",
	}
}

// Messages 覆盖默认提示。
func (r *MemberStatusRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"status.required": "会员状态不能为空",
		"status.integer":  "会员状态取值不合法",
		"status.in":       "会员状态取值不合法",
	}
}
