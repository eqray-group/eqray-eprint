package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/pay/admin/request"
	"eqray-eprint/app/modules/pay/service"
)

// ChannelConfigController 提供支付渠道配置的后台维护接口。
type ChannelConfigController struct {
	service *service.ChannelConfigService
}

// NewChannelConfigController 构造渠道配置控制器。
func NewChannelConfigController() *ChannelConfigController {
	return &ChannelConfigController{service: service.NewChannelConfigService()}
}

// Index 分页查询渠道配置；列表不返回密钥，避免日常浏览就暴露私钥。
// 入参：code、channel、status、pageNo、pageSize（查询串）；返回配置列表与总数。
func (c *ChannelConfigController) Index(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	q := service.ChannelConfigListQuery{
		Code:     ctx.Request().Query("code"),
		Channel:  ctx.Request().Query("channel"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := int8(status)
		q.Status = &v
	}

	list, total, err := c.service.Page(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询单条配置并回显密钥，仅供后台编辑页使用。
// 入参：路径 id；返回配置详情（含 privateKey、apiKey、platformPublicKey）。
func (c *ChannelConfigController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	detail, err := c.service.Detail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, detail)
}

// Store 新增渠道配置。
// 入参见 ChannelConfigSaveRequest；返回新配置 ID。
func (c *ChannelConfigController) Store(ctx http.Context) http.Response {
	var req request.ChannelConfigSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	id, err := c.service.Create(ctx.Context(), toChannelConfigInput(req))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 编辑渠道配置；密钥字段留空表示保持原值。
// 入参：路径 id + 配置字段；返回 null。
func (c *ChannelConfigController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.ChannelConfigSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, toChannelConfigInput(req)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Destroy 删除渠道配置（软删）；历史支付单仍按 channelCode 追溯。
// 入参：路径 id；返回 null。
func (c *ChannelConfigController) Destroy(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Delete(ctx.Context(), id); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// toChannelConfigInput 把请求体转换为服务层入参。
func toChannelConfigInput(req request.ChannelConfigSaveRequest) service.ChannelConfigInput {
	return service.ChannelConfigInput{
		Code:              req.Code,
		Name:              req.Name,
		Channel:           req.Channel,
		Scene:             req.Scene,
		AppID:             req.AppID,
		MchID:             req.MchID,
		SerialNo:          req.SerialNo,
		PrivateKey:        req.PrivateKey,
		APIKey:            req.APIKey,
		PlatformPublicKey: req.PlatformPublicKey,
		Gateway:           req.Gateway,
		NotifyURL:         req.NotifyURL,
		ReturnURL:         req.ReturnURL,
		Extra:             req.Extra,
		Status:            req.Status,
		Sort:              req.Sort,
		Remark:            req.Remark,
	}
}
