using System;
using System.Configuration;
using System.Reflection;

namespace eprint.Services
{
    /// <summary>
    /// 客户端运行配置：服务端地址、版本号与打印引擎路径。
    /// 服务端地址以独立配置文件（<see cref="ClientConfig"/>）为准，App.config 只提供首次启动的默认值。
    /// </summary>
    public static class AppConfig
    {
        /// <summary>
        /// 服务端根地址，不带结尾斜杠：读写独立配置文件里的值，改完由
        /// <see cref="ClientConfig.Save"/> 落盘，下次启动直接复用。
        /// </summary>
        public static string ApiBaseUrl
        {
            get { return ClientConfig.Current.ApiBaseUrl; }
            set { ClientConfig.Current.ApiBaseUrl = Normalize(value); }
        }

        /// <summary>
        /// App.config 里配置的初始服务端地址，仅在独立配置文件没有记录时作为兜底。
        /// </summary>
        public static string DefaultApiBaseUrl
        {
            get
            {
                string value = ConfigurationManager.AppSettings["ApiBaseUrl"];
                if (string.IsNullOrWhiteSpace(value))
                {
                    return "http://127.0.0.1:11111";
                }

                return Normalize(value);
            }
        }

        /// <summary>
        /// 客户端版本号，随心跳与设备登记上报。
        /// </summary>
        public static string ClientVersion
        {
            get { return Assembly.GetExecutingAssembly().GetName().Version.ToString(3); }
        }

        /// <summary>
        /// SumatraPDF 便携版路径；相对路径基于客户端目录。
        /// 留空时按 tools/SumatraPDF.exe → 同目录 → Program Files 的顺序自动查找。
        /// </summary>
        public static string SumatraPdfPath { get; set; } =
            (ConfigurationManager.AppSettings["SumatraPdfPath"] ?? string.Empty).Trim();

        /// <summary>
        /// 规范化服务端地址：去空白、补全 http:// 前缀、去掉结尾斜杠；空输入返回空串。
        /// </summary>
        public static string Normalize(string value)
        {
            string text = (value ?? string.Empty).Trim();
            if (text.Length == 0)
            {
                return string.Empty;
            }

            if (!text.StartsWith("http://", StringComparison.OrdinalIgnoreCase) &&
                !text.StartsWith("https://", StringComparison.OrdinalIgnoreCase))
            {
                text = "http://" + text;
            }

            return text.TrimEnd('/');
        }
    }
}
