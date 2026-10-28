namespace eprint.Models
{
    /// <summary>
    /// 本机打印机条目：系统队列名、默认标记与从驱动探测到的能力。
    /// </summary>
    public sealed class PrinterInfo
    {
        /// <summary>
        /// 打印机显示名称，与 Windows「设备和打印机」中显示的名称一致。
        /// </summary>
        public string Name { get; set; }

        /// <summary>
        /// Windows 打印队列名，静默打印与上报服务端都以它为准。
        /// </summary>
        public string SystemName { get; set; }

        /// <summary>
        /// 型号，当前版本留空，后续可通过 WMI 补充。
        /// </summary>
        public string Model { get; set; }

        /// <summary>
        /// 驱动是否报告支持 A4 纸张。
        /// </summary>
        public bool SupportsA4 { get; set; }

        /// <summary>
        /// 驱动是否报告支持 A3 纸张。
        /// </summary>
        public bool SupportsA3 { get; set; }

        /// <summary>
        /// 驱动是否报告支持彩色打印。
        /// </summary>
        public bool SupportsColor { get; set; }

        /// <summary>
        /// 驱动是否报告支持双面打印。
        /// </summary>
        public bool SupportsDuplex { get; set; }

        /// <summary>
        /// 是否为系统默认打印机。
        /// </summary>
        public bool IsDefault { get; set; }
    }
}
