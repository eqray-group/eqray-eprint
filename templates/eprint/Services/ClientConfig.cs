using System;
using System.IO;
using System.Security.Cryptography;
using System.Text;

namespace eprint.Services
{
    /// <summary>
    /// 客户端独立配置文件（%LOCALAPPDATA%\yiyunyin\client.json）。
    /// 保存服务端地址、上次登录账号与自动登录所需的刷新令牌：登录页改完立即写盘，下次启动直接复用，
    /// 不再依赖 App.config。密码与刷新令牌落盘前经 Windows DPAPI 按当前用户加密，
    /// 换账号或把文件拷到别的机器都无法解密。
    /// </summary>
    public sealed class ClientConfig
    {
        private const string FolderName = "yiyunyin";
        private const string FileName = "client.json";
        private static ClientConfig _current;

        /// <summary>
        /// 服务端根地址，不带结尾斜杠；登录页修改后写回本属性并落盘。
        /// </summary>
        public string ApiBaseUrl { get; set; }

        /// <summary>
        /// 下次启动是否用刷新令牌自动登录。
        /// </summary>
        public bool AutoLogin { get; set; }

        /// <summary>
        /// 上次登录的账号，用于登录页预填。
        /// </summary>
        public string Username { get; set; }

        /// <summary>
        /// 上次登录的密码（仅在勾选自动登录时保存），用于登录页预填。
        /// </summary>
        public string Password { get; set; }

        /// <summary>
        /// 自动登录用的刷新令牌；每次续期成功后服务端会轮换旧令牌，必须立刻回写本属性并保存。
        /// </summary>
        public string RefreshToken { get; set; }

        /// <summary>
        /// 当前生效的配置：首次访问时从独立配置文件读取，文件不存在则用 App.config 里的默认值。
        /// </summary>
        public static ClientConfig Current
        {
            get
            {
                if (_current == null)
                {
                    _current = Load();
                }

                return _current;
            }
        }

        /// <summary>
        /// 配置文件绝对路径：%LOCALAPPDATA%\yiyunyin\client.json。
        /// </summary>
        public static string FilePath
        {
            get
            {
                string folder = Path.Combine(
                    Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), FolderName);
                return Path.Combine(folder, FileName);
            }
        }

        /// <summary>
        /// 把当前配置写回独立配置文件，返回是否写入成功；
        /// 目录不存在会自动创建，写入失败只影响「下次启动是否记住」，不抛异常打断界面。
        /// </summary>
        public bool Save()
        {
            try
            {
                string path = FilePath;
                string folder = Path.GetDirectoryName(path);
                if (!string.IsNullOrEmpty(folder) && !Directory.Exists(folder))
                {
                    Directory.CreateDirectory(folder);
                }

                var file = new ClientConfigFile
                {
                    apiBaseUrl = AppConfig.Normalize(ApiBaseUrl),
                    autoLogin = AutoLogin,
                    username = Username ?? string.Empty,
                    password = Protect(Password),
                    refreshToken = Protect(RefreshToken)
                };

                File.WriteAllText(path, Json.Serialize(file), new UTF8Encoding(false));
                return true;
            }
            catch (Exception)
            {
                return false;
            }
        }

        /// <summary>
        /// 读取配置文件；文件缺失、内容损坏或解密失败时回落 App.config 的默认地址，
        /// 保证客户端总能以「未登录」状态启动。
        /// </summary>
        private static ClientConfig Load()
        {
            var config = new ClientConfig
            {
                ApiBaseUrl = AppConfig.DefaultApiBaseUrl,
                Username = string.Empty,
                Password = string.Empty,
                RefreshToken = string.Empty
            };

            try
            {
                string path = FilePath;
                if (!File.Exists(path))
                {
                    return config;
                }

                string text = File.ReadAllText(path, Encoding.UTF8);
                if (string.IsNullOrWhiteSpace(text))
                {
                    return config;
                }

                ClientConfigFile file = Json.Deserialize<ClientConfigFile>(text);
                if (file == null)
                {
                    return config;
                }

                if (!string.IsNullOrWhiteSpace(file.apiBaseUrl))
                {
                    config.ApiBaseUrl = AppConfig.Normalize(file.apiBaseUrl);
                }
                config.AutoLogin = file.autoLogin;
                config.Username = file.username ?? string.Empty;
                config.Password = Unprotect(file.password);
                config.RefreshToken = Unprotect(file.refreshToken);
                return config;
            }
            catch (Exception)
            {
                return config;
            }
        }

        /// <summary>
        /// 用 DPAPI（当前用户范围）加密待写入的敏感文本，空值原样返回空串。
        /// </summary>
        private static string Protect(string value)
        {
            if (string.IsNullOrEmpty(value))
            {
                return string.Empty;
            }

            try
            {
                byte[] cipher = ProtectedData.Protect(
                    Encoding.UTF8.GetBytes(value), null, DataProtectionScope.CurrentUser);
                return Convert.ToBase64String(cipher);
            }
            catch (Exception)
            {
                return string.Empty;
            }
        }

        /// <summary>
        /// 解密配置文件里的敏感文本；解密失败（换用户 / 换机器）返回空串，由用户重新登录。
        /// </summary>
        private static string Unprotect(string value)
        {
            if (string.IsNullOrEmpty(value))
            {
                return string.Empty;
            }

            try
            {
                byte[] plain = ProtectedData.Unprotect(
                    Convert.FromBase64String(value), null, DataProtectionScope.CurrentUser);
                return Encoding.UTF8.GetString(plain);
            }
            catch (Exception)
            {
                return string.Empty;
            }
        }
    }

    /// <summary>
    /// client.json 的落盘结构：属性名即 JSON 键名，password 与 refreshToken 存的是 DPAPI 密文。
    /// 单独建类是为了让内存中的明文与磁盘上的密文互不干扰。
    /// </summary>
    public sealed class ClientConfigFile
    {
        /// <summary>服务端根地址。</summary>
        public string apiBaseUrl { get; set; }

        /// <summary>下次启动是否自动登录。</summary>
        public bool autoLogin { get; set; }

        /// <summary>上次登录账号。</summary>
        public string username { get; set; }

        /// <summary>DPAPI 加密后的密码。</summary>
        public string password { get; set; }

        /// <summary>DPAPI 加密后的刷新令牌。</summary>
        public string refreshToken { get; set; }
    }
}
