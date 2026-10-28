using System;
using System.Collections.Generic;
using System.Drawing.Printing;
using System.Linq;
using System.Management;
using eprint.Models;

namespace eprint.Services
{
    /// <summary>
    /// 本机打印机查询服务：读取 Windows 已安装打印机，并从驱动探测纸张 / 彩色 / 双面能力。
    /// 探测结果只作为界面默认值，用户可以在列表中调整后再上报。
    /// </summary>
    public static class PrinterService
    {
        // A4 / A3 在 Windows 打印驱动中的标准尺寸（1/100 英寸）
        private const int A4Width = 827;
        private const int A4Height = 1169;
        private const int A3Width = 1169;
        private const int A3Height = 1654;

        // 部分驱动上报的尺寸与标准值有少量偏差，允许 20/100 英寸误差
        private const int SizeTolerance = 20;

        /// <summary>
        /// 读取本机已安装的打印机列表（本地打印机与已添加的网络打印机都会返回）。
        /// 返回：默认打印机排在最前，其余按名称升序排列；读取失败时抛出异常，由调用方提示用户。
        /// </summary>
        public static List<PrinterInfo> GetInstalledPrinters()
        {
            string defaultName = TryGetDefaultPrinterName();
            Dictionary<string, string> drivers = TryLoadDriverNames();

            var printers = new List<PrinterInfo>();
            foreach (string name in PrinterSettings.InstalledPrinters)
            {
                var info = new PrinterInfo
                {
                    Name = name,
                    SystemName = name,
                    Model = drivers.ContainsKey(name) ? drivers[name] : string.Empty,
                    IsDefault = string.Equals(name, defaultName, StringComparison.OrdinalIgnoreCase),
                    SupportsA4 = true
                };

                try
                {
                    var settings = new PrinterSettings { PrinterName = name };
                    info.SupportsA4 = SupportsPaper(settings, PaperKind.A4, A4Width, A4Height);
                    info.SupportsA3 = SupportsPaper(settings, PaperKind.A3, A3Width, A3Height);
                    info.SupportsDuplex = settings.CanDuplex;
                    info.SupportsColor = settings.SupportsColor;
                }
                catch
                {
                    // 驱动异常时按最常见的 A4 能力兜底，用户仍可手动调整
                }

                // 后端只支持 A4 / A3，两者都没探测到时至少给出 A4 默认值
                if (!info.SupportsA4 && !info.SupportsA3)
                {
                    info.SupportsA4 = true;
                }

                printers.Add(info);
            }

            return printers
                .OrderByDescending(p => p.IsDefault)
                .ThenBy(p => p.Name, StringComparer.CurrentCultureIgnoreCase)
                .ToList();
        }

        /// <summary>
        /// 读取系统默认打印机名称。
        /// 返回：默认打印机名称；未设置默认打印机或读取失败时返回 null。
        /// </summary>
        public static string TryGetDefaultPrinterName()
        {
            try
            {
                string name = new PrinterSettings().PrinterName;
                return string.IsNullOrWhiteSpace(name) ? null : name;
            }
            catch
            {
                // 打印后台服务不可用等情况不阻断列表读取，默认打印机按“无”处理
                return null;
            }
        }

        /// <summary>
        /// 判断指定打印驱动是否支持某种纸张：优先比对标准 PaperKind，其次比对毫米尺寸。
        /// </summary>
        private static bool SupportsPaper(PrinterSettings settings, PaperKind kind, int width, int height)
        {
            foreach (PaperSize size in settings.PaperSizes)
            {
                if (size.RawKind == (int)kind)
                {
                    return true;
                }

                if (Math.Abs(size.Width - width) <= SizeTolerance &&
                    Math.Abs(size.Height - height) <= SizeTolerance)
                {
                    return true;
                }
            }

            return false;
        }

        /// <summary>
        /// 通过 WMI 一次性读取打印机驱动名称，用作界面上的“型号”展示。
        /// 返回：队列名到驱动名的映射；WMI 不可用时返回空字典。
        /// </summary>
        private static Dictionary<string, string> TryLoadDriverNames()
        {
            var names = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);

            try
            {
                using (var searcher = new ManagementObjectSearcher("SELECT Name, DriverName FROM Win32_Printer"))
                {
                    foreach (ManagementObject item in searcher.Get())
                    {
                        try
                        {
                            string name = Convert.ToString(item["Name"]);
                            string driver = Convert.ToString(item["DriverName"]);
                            if (!string.IsNullOrWhiteSpace(name))
                            {
                                names[name] = driver ?? string.Empty;
                            }
                        }
                        finally
                        {
                            item.Dispose();
                        }
                    }
                }
            }
            catch
            {
                // WMI 查询失败不影响打印机列表，只是型号列为空
            }

            return names;
        }
    }
}
