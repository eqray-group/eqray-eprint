// Package errs 定义 member 模块的业务错误码，编码规则见 docs 8.3(2)。
package errs

import baseErrs "eqray-eprint/app/modules/base/errs"

// member 模块业务错误码：1-003-{子模块}-{错误码}，
// user 占用 001，social 占用 005，auth 占用 006，address 占用 007，card 占用 009。
// 等级（002）、标签（003）、分组（004）、签到（008）体系已下线，对应错误码一并移除。
var (
	MemberNotExists    = baseErrs.ErrorCode{Code: 1_003_001_001, Msg: "会员【%v】不存在"}
	MemberMobileExists = baseErrs.ErrorCode{Code: 1_003_001_002, Msg: "手机号【%v】已被其他会员使用"}

	// social 占用 005 子模块：第三方登录渠道、授权码与绑定关系
	MemberSocialClientUnavailable = baseErrs.ErrorCode{Code: 1_003_005_001, Msg: "第三方登录渠道【%v】未配置或已停用"}
	MemberSocialCodeInvalid       = baseErrs.ErrorCode{Code: 1_003_005_002, Msg: "第三方授权码无效或已过期"}
	MemberSocialAccountNotBound   = baseErrs.ErrorCode{Code: 1_003_005_003, Msg: "该第三方账号尚未绑定会员"}
	MemberSocialAlreadyBound      = baseErrs.ErrorCode{Code: 1_003_005_004, Msg: "该第三方账号已绑定其他会员"}
	MemberSocialBindingNotExists  = baseErrs.ErrorCode{Code: 1_003_005_005, Msg: "绑定关系【%v】不存在"}
	MemberSocialLastBinding       = baseErrs.ErrorCode{Code: 1_003_005_006, Msg: "至少保留一个第三方绑定，不能全部解绑"}
	MemberSocialRequestFailed     = baseErrs.ErrorCode{Code: 1_003_005_007, Msg: "第三方接口调用失败：%v"}
	MemberDisabled                = baseErrs.ErrorCode{Code: 1_003_005_008, Msg: "会员账号已被禁用"}
	MemberSocialTypeDuplicated    = baseErrs.ErrorCode{Code: 1_003_005_009, Msg: "该会员已绑定【%v】，请先解绑"}
	MemberSocialClientNotExists   = baseErrs.ErrorCode{Code: 1_003_005_010, Msg: "第三方登录配置【%v】不存在"}
	MemberSocialClientTypeExists  = baseErrs.ErrorCode{Code: 1_003_005_011, Msg: "第三方登录渠道【%v】已配置"}

	// auth 占用 006 子模块：手机号短信验证码登录 / 注册
	MemberMobileInvalid         = baseErrs.ErrorCode{Code: 1_003_006_001, Msg: "手机号格式不正确"}
	MemberSmsCodeTooFrequent    = baseErrs.ErrorCode{Code: 1_003_006_002, Msg: "验证码发送过于频繁，请稍后再试"}
	MemberSmsCodeDailyLimit     = baseErrs.ErrorCode{Code: 1_003_006_003, Msg: "该手机号今日验证码发送次数已达上限"}
	MemberSmsSendFailed         = baseErrs.ErrorCode{Code: 1_003_006_004, Msg: "短信验证码发送失败"}
	MemberSmsCodeWrong          = baseErrs.ErrorCode{Code: 1_003_006_005, Msg: "短信验证码不正确或已过期"}
	MemberSmsChannelNotExists   = baseErrs.ErrorCode{Code: 1_003_006_006, Msg: "短信渠道【%v】不存在或已停用"}
	MemberSmsTemplateNotExists  = baseErrs.ErrorCode{Code: 1_003_006_007, Msg: "短信模板【%v】不存在或已停用"}
	MemberSmsChannelCodeExists  = baseErrs.ErrorCode{Code: 1_003_006_008, Msg: "短信渠道代码【%v】已存在"}
	MemberSmsTemplateCodeExists = baseErrs.ErrorCode{Code: 1_003_006_009, Msg: "短信模板代码【%v】已存在"}
	MemberSmsChannelInUse       = baseErrs.ErrorCode{Code: 1_003_006_010, Msg: "渠道【%v】下仍有模板，不可删除"}
	// MemberSmsSceneInvalid 用于发送验证码时传入了白名单之外的场景标识
	MemberSmsSceneInvalid = baseErrs.ErrorCode{Code: 1_003_006_011, Msg: "不支持的短信场景【%v】"}

	// address 占用 007 子模块：C 端会员自助维护收货地址
	MemberAddressNotExists     = baseErrs.ErrorCode{Code: 1_003_007_001, Msg: "收货地址【%v】不存在"}
	MemberAddressLimitExceeded = baseErrs.ErrorCode{Code: 1_003_007_002, Msg: "收货地址最多保存 %v 条"}

	// card 占用 009 子模块：会员卡发卡、改卡与折扣
	MemberCardNotExists      = baseErrs.ErrorCode{Code: 1_003_009_001, Msg: "会员【%v】还没有会员卡"}
	MemberCardAlreadyExists  = baseErrs.ErrorCode{Code: 1_003_009_002, Msg: "会员【%v】已有会员卡，请直接修改"}
	MemberCardDiscountInvalid = baseErrs.ErrorCode{Code: 1_003_009_003, Msg: "折扣百分比取值不合法：%v"}
)
