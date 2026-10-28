using System.Collections.Generic;

namespace eprint.Models
{
    /// <summary>
    /// 服务端统一响应体 { code, msg, data }；code 为 0 表示业务成功。
    /// </summary>
    public sealed class ApiEnvelope<T>
    {
        public int Code { get; set; }

        public string Msg { get; set; }

        public T Data { get; set; }
    }

    /// <summary>
    /// 登录成功后签发的访问令牌 / 刷新令牌。
    /// </summary>
    public sealed class TokenPair
    {
        public string AccessToken { get; set; }

        public string RefreshToken { get; set; }

        public long ExpiresIn { get; set; }

        public long RefreshExpiresIn { get; set; }
    }

    /// <summary>
    /// 图形验证码：captchaId 用于登录时回传，image 为 base64 图片。
    /// </summary>
    public sealed class CaptchaInfo
    {
        public string CaptchaId { get; set; }

        public string Image { get; set; }
    }

    /// <summary>
    /// 当前账号绑定的门店，客户端登录后据此选择设备归属。
    /// </summary>
    public sealed class ShopOption
    {
        public long Id { get; set; }

        public string Name { get; set; }

        public string Role { get; set; }

        public int Status { get; set; }

        /// <summary>
        /// 下拉框显示文案：停业门店加后缀提示，不影响选择。
        /// </summary>
        public string DisplayName
        {
            get { return Status == 1 ? Name : Name + "（停业）"; }
        }
    }

    /// <summary>
    /// 登记成功后的设备记录，ID 用于后续心跳、打印机与任务接口。
    /// </summary>
    public sealed class DeviceInfo
    {
        public long Id { get; set; }

        public long PrintShopId { get; set; }

        public string Code { get; set; }

        public string Name { get; set; }

        public string ClientVersion { get; set; }

        public int Status { get; set; }

        public int OnlineStatus { get; set; }
    }

    /// <summary>
    /// 服务端保存的一台打印机，包含门店后台 / 客户端配置的能力。
    /// </summary>
    public sealed class PrinterDto
    {
        public long Id { get; set; }

        public long PrintShopId { get; set; }

        public long DeviceId { get; set; }

        public string Name { get; set; }

        public string SystemName { get; set; }

        public string Model { get; set; }

        public string PaperTypes { get; set; }

        public int ColorSupport { get; set; }

        public int DuplexSupport { get; set; }

        public int Status { get; set; }
    }

    /// <summary>
    /// 客户端上报的一台打印机；SystemName 为 Windows 打印队列名。
    /// </summary>
    public sealed class PrinterSyncItem
    {
        public string SystemName { get; set; }

        public string Name { get; set; }

        public string Model { get; set; }

        public string PaperTypes { get; set; }

        public int ColorSupport { get; set; }

        public int DuplexSupport { get; set; }
    }

    /// <summary>
    /// 服务端下发的打印任务：纸张、色彩、单双面等参数以服务端为准，客户端只读。
    /// FileUrl 可能是绝对地址，也可能是 /storage 开头的相对地址。
    /// </summary>
    public sealed class TaskDto
    {
        public string TaskNo { get; set; }

        public string OrderNo { get; set; }

        public long FileId { get; set; }

        public string FileUrl { get; set; }

        public string FileName { get; set; }

        public int FilePageCount { get; set; }

        public string PaperType { get; set; }

        public string ColorType { get; set; }

        public string PrintMode { get; set; }

        public string Orientation { get; set; }

        public string PageRange { get; set; }

        public int PageCount { get; set; }

        public int Copies { get; set; }

        public string SystemName { get; set; }
    }

    /// <summary>
    /// 列表里的打印机行：勾选状态 + 本地探测能力，绑定到 DataGridView。
    /// </summary>
    public sealed class PrinterRow
    {
        public bool Selected { get; set; }

        public string Name { get; set; }

        public string SystemName { get; set; }

        public string Model { get; set; }

        public bool A4 { get; set; }

        public bool A3 { get; set; }

        public bool Color { get; set; }

        public bool Duplex { get; set; }

        public bool IsDefault { get; set; }

        public string DefaultText
        {
            get { return IsDefault ? "是" : string.Empty; }
        }
    }
}
