// Package errs 定义 print 模块的业务错误码。
package errs

import baseErrs "eqray-eprint/app/modules/base/errs"

// print 模块错误码：1-004-{子模块}-{序号}。
// shop 001，staff 002，device 003，printer 004，price 005，order 006，task 007，file 008。
var (
	ShopNotExists    = baseErrs.ErrorCode{Code: 1_004_001_001, Msg: "门店【%v】不存在"}
	ShopClosed       = baseErrs.ErrorCode{Code: 1_004_001_002, Msg: "门店已停业，暂不可下单"}
	ShopNotAudited   = baseErrs.ErrorCode{Code: 1_004_001_003, Msg: "门店未通过审核，暂不可下单"}
	ShopAuditInvalid = baseErrs.ErrorCode{Code: 1_004_001_004, Msg: "当前审核状态不允许该操作"}

	StaffExists    = baseErrs.ErrorCode{Code: 1_004_002_001, Msg: "该账号已绑定此门店"}
	StaffNotExists = baseErrs.ErrorCode{Code: 1_004_002_002, Msg: "门店成员不存在"}
	StaffNotBelong = baseErrs.ErrorCode{Code: 1_004_002_003, Msg: "当前账号不属于该门店"}

	DeviceNotExists  = baseErrs.ErrorCode{Code: 1_004_003_001, Msg: "设备【%v】不存在"}
	DeviceDisabled   = baseErrs.ErrorCode{Code: 1_004_003_002, Msg: "设备已禁用，请联系管理员"}
	DeviceOffline    = baseErrs.ErrorCode{Code: 1_004_003_003, Msg: "门店打印设备不在线，请稍后再试"}
	DeviceCodeExists = baseErrs.ErrorCode{Code: 1_004_003_004, Msg: "设备编码【%v】已被占用"}

	PrinterNotExists  = baseErrs.ErrorCode{Code: 1_004_004_001, Msg: "打印机【%v】不存在"}
	PrinterDisabled   = baseErrs.ErrorCode{Code: 1_004_004_002, Msg: "打印机已停用"}
	PrinterNoMatch    = baseErrs.ErrorCode{Code: 1_004_004_003, Msg: "本店暂不支持该打印参数：%v"}
	PrinterNameExists = baseErrs.ErrorCode{Code: 1_004_004_004, Msg: "同一设备下打印队列【%v】已登记"}

	PriceNotExists  = baseErrs.ErrorCode{Code: 1_004_005_001, Msg: "价格规则不存在"}
	PriceRuleExists = baseErrs.ErrorCode{Code: 1_004_005_002, Msg: "该参数组合的价格已配置"}
	PriceInvalid    = baseErrs.ErrorCode{Code: 1_004_005_003, Msg: "单价必须大于 0"}

	OrderNotExists     = baseErrs.ErrorCode{Code: 1_004_006_001, Msg: "订单【%v】不存在"}
	OrderStatusInvalid = baseErrs.ErrorCode{Code: 1_004_006_002, Msg: "订单当前状态不允许该操作"}
	OrderItemEmpty     = baseErrs.ErrorCode{Code: 1_004_006_003, Msg: "订单至少需要一个文件"}
	OrderPickupInvalid = baseErrs.ErrorCode{Code: 1_004_006_004, Msg: "取件码不正确或订单已完成"}

	TaskNotExists  = baseErrs.ErrorCode{Code: 1_004_007_001, Msg: "打印任务不存在"}
	TaskRetryLimit = baseErrs.ErrorCode{Code: 1_004_007_002, Msg: "任务重试次数已达上限"}

	FileNotExists   = baseErrs.ErrorCode{Code: 1_004_008_001, Msg: "打印文件不存在"}
	FileNotReady    = baseErrs.ErrorCode{Code: 1_004_008_002, Msg: "文件尚未转换完成，请稍后重试"}
	FileUnsupported = baseErrs.ErrorCode{Code: 1_004_008_003, Msg: "不支持的文件类型：%v"}

	BannerNotExists = baseErrs.ErrorCode{Code: 1_004_009_001, Msg: "轮播图不存在"}
)
