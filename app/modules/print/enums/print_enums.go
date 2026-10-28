// Package enums 定义打印模块的枚举常量，数据库存数值或英文 code，代码里只用常量。
package enums

// 订单状态：单线推进 0 → 1 → 2 → 3 → 4，旁支 5/6/7/8。
const (
	OrderStatusUnpaid    int8 = 0 // 待支付
	OrderStatusPending   int8 = 1 // 待打印
	OrderStatusPrinting  int8 = 2 // 打印中
	OrderStatusPickup    int8 = 3 // 待取件
	OrderStatusCompleted int8 = 4 // 已完成
	OrderStatusCanceled  int8 = 5 // 已取消
	OrderStatusRefunding int8 = 6 // 退款中
	OrderStatusRefunded  int8 = 7 // 已退款
	OrderStatusExpired   int8 = 8 // 超期关闭
)

// 支付状态：只表达钱，履约状态看 OrderStatus。
const (
	PayStatusUnpaid    int8 = 0 // 未支付
	PayStatusPaid      int8 = 1 // 已支付
	PayStatusRefunding int8 = 2 // 退款中
	PayStatusRefunded  int8 = 3 // 已退款
)

// 支付方式，与 pay 模块的 channel 对应。
const (
	PayTypeWechat = "wechat"
	PayTypeAlipay = "alipay"
	PayTypeWallet = "wallet"
)

// 明细任务状态。
const (
	TaskStatusPending  int8 = 0 // 待分配
	TaskStatusAssigned int8 = 1 // 已分配
	TaskStatusPrinting int8 = 2 // 打印中
	TaskStatusDone     int8 = 3 // 已完成
	TaskStatusFailed   int8 = 4 // 失败
	TaskStatusCanceled int8 = 5 // 已取消
)

// 打印参数：纸张、色彩、单双面、方向。
const (
	PaperA4 = "A4"
	PaperA3 = "A3"

	ColorBW    = "bw"
	ColorColor = "color"

	ModeSingle = "single"
	ModeDouble = "double"

	OrientationPortrait  = "portrait"
	OrientationLandscape = "landscape"
)

// 文件转换状态。
const (
	ConvertStatusPending    int8 = 0 // 待处理
	ConvertStatusConverting int8 = 1 // 转换中
	ConvertStatusReady      int8 = 2 // 可打印
	ConvertStatusFailed     int8 = 3 // 转换失败
)

// 门店营业状态与审核状态。
const (
	ShopStatusClosed int8 = 0 // 停业
	ShopStatusOpen   int8 = 1 // 营业

	AuditStatusPending  int8 = 0 // 待审核
	AuditStatusApproved int8 = 1 // 审核通过
	AuditStatusRejected int8 = 2 // 审核拒绝
)

// 设备在线状态。
const (
	DeviceOffline int8 = 0 // 离线
	DeviceOnline  int8 = 1 // 在线
)

// 门店成员角色。
const (
	StaffRoleOwner = "owner" // 店主，可管理设备与价格
	StaffRoleStaff = "staff" // 店员，只能核销与查看订单
)

// 通用启用状态。
const (
	StatusDisabled int8 = 0 // 停用
	StatusEnabled  int8 = 1 // 启用
)

// 轮播图跳转类型。
const (
	BannerLinkNone = "none"
	BannerLinkURL  = "url"
	BannerLinkPage = "page"
)

// 任务失败错误码，客户端回执时上报。
const (
	ErrorPaperUnsupported  = "PAPER_UNSUPPORTED"
	ErrorColorUnsupported  = "COLOR_UNSUPPORTED"
	ErrorDuplexUnsupported = "DUPLEX_UNSUPPORTED"
	ErrorPrinterOffline    = "PRINTER_OFFLINE"
	ErrorFileBroken        = "FILE_BROKEN"
	ErrorPaperJam          = "PAPER_JAM"
	ErrorUnknown           = "UNKNOWN"
)
