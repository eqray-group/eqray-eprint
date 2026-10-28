package controller

import (
	"github.com/goravel/framework/contracts/http"

	baseEnums "eqray-eprint/app/modules/base/enums"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/member/admin/request"
	"eqray-eprint/app/modules/member/service"
)

// MemberController 提供会员档案的增删改查与地址、第三方绑定只读接口。
type MemberController struct {
	service *service.MemberService
	socials *service.MemberSocialService
}

// NewMemberController 构造会员档案控制器，并装配会员与社交服务。
func NewMemberController() *MemberController {
	return &MemberController{
		service: service.NewMemberService(),
		socials: service.NewMemberSocialService(),
	}
}

// Index 分页查询会员档案，支持昵称、手机号与状态筛选。
// 入参：nickname、mobile、status、pageNo、pageSize（查询串）；返回会员列表与总数。
func (c *MemberController) Index(ctx http.Context) http.Response {
	page := param.Page{
		PageNo:   ctx.Request().QueryInt("pageNo", 1),
		PageSize: ctx.Request().QueryInt("pageSize", param.DefaultPageSize),
	}
	if err := page.Normalize(); err != nil {
		return result.Fail(ctx, err)
	}

	q := service.MemberListQuery{
		Nickname: ctx.Request().Query("nickname"),
		Mobile:   ctx.Request().Query("mobile"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if status := ctx.Request().QueryInt("status", -1); status >= 0 {
		v := baseEnums.CommonStatusEnum(status)
		q.Status = &v
	}

	list, total, err := c.service.Page(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Show 查询会员详情。
// 入参：路径 id；返回会员详情，不存在返回业务错误码。
func (c *MemberController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	m, err := c.service.Detail(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, m)
}

// Store 后台新增会员档案；手机号唯一，重复时返回业务错误码。
// 入参见 MemberSaveRequest；返回新会员 ID。
func (c *MemberController) Store(ctx http.Context) http.Response {
	var req request.MemberSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	id, err := c.service.Create(ctx.Context(), toMemberSaveInput(req))
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, map[string]any{"id": id})
}

// Update 编辑会员档案。
// 入参：路径 id + 会员字段；返回 null。
func (c *MemberController) Update(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.MemberSaveRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Update(ctx.Context(), id, toMemberSaveInput(req)); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// UpdateStatus 启用或禁用会员，禁用后 C 端接口会拒绝该会员。
// 入参：路径 id + status（0 禁用 / 1 启用）；返回 null。
func (c *MemberController) UpdateStatus(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	var req request.MemberStatusRequest
	if err := param.Bind(ctx, &req); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.UpdateStatus(ctx.Context(), id, req.Status); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// Addresses 查询指定会员的收货地址列表（只读），默认地址排最前。
// 入参：路径 id；返回地址列表。
func (c *MemberController) Addresses(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	list, err := c.service.Addresses(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, list)
}

// SocialBindings 管理端只读查看会员绑定的第三方账号，不支持在此解绑。
// 入参：路径 id；返回绑定记录列表。
func (c *MemberController) SocialBindings(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	list, err := c.socials.Bindings(ctx.Context(), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, list)
}

// Destroy 删除会员档案（软删），已产生的会员卡与第三方绑定保留。
// 入参：路径 id；返回 null。
func (c *MemberController) Destroy(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Delete(ctx.Context(), id); err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, nil)
}

// toMemberSaveInput 把请求体转换为服务层入参。
func toMemberSaveInput(req request.MemberSaveRequest) service.MemberSaveInput {
	in := service.MemberSaveInput{
		Nickname: req.Nickname,
		Mobile:   req.Mobile,
		Avatar:   req.Avatar,
		Name:     req.Name,
		Sex:      req.Sex,
		AreaID:   req.AreaID,
		Mark:     req.Mark,
		Status:   req.Status,
	}
	if req.Birthday != nil && !req.Birthday.IsZero() {
		birthday := *req.Birthday
		in.Birthday = &birthday
	}

	return in
}
