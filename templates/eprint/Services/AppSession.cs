using System.Collections.Generic;
using eprint.Models;

namespace eprint.Services
{
    /// <summary>
    /// 本次运行期的登录态：HTTP 客户端、令牌、门店与设备。
    /// 只保存在内存，进程退出后需要重新登录。
    /// </summary>
    public static class AppSession
    {
        public static ApiClient Api { get; set; }

        public static BackendApi Backend { get; set; }

        public static TokenPair Tokens { get; set; }

        public static string Username { get; set; }

        public static List<ShopOption> Shops { get; set; }

        public static ShopOption CurrentShop { get; set; }

        public static DeviceInfo Device { get; set; }

        /// <summary>
        /// 退出登录时清空全部运行期状态，避免下一个账号复用旧令牌。
        /// </summary>
        public static void Reset()
        {
            Api = null;
            Backend = null;
            Tokens = null;
            Username = null;
            Shops = null;
            CurrentShop = null;
            Device = null;
        }
    }
}
