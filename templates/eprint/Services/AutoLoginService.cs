using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using eprint.Models;

namespace eprint.Services
{
    /// <summary>
    /// 自动登录：用独立配置文件里的刷新令牌换一对新令牌，成功后直接填充 AppSession，
    /// 免去启动时输入账号、密码与验证码；令牌失效则清空本地令牌并回到登录页。
    /// </summary>
    public static class AutoLoginService
    {
        /// <summary>
        /// 尝试恢复上次登录态：成功返回 true，此时 AppSession 已可支撑主界面；
        /// 返回 false 时 message 为登录页要展示的提示（无令牌时为空）。
        /// </summary>
        public static bool TryRestore(out string message)
        {
            message = null;
            ClientConfig config = ClientConfig.Current;
            if (string.IsNullOrEmpty(config.RefreshToken))
            {
                return false;
            }

            try
            {
                // 放到线程池执行，避免在界面线程上同步等待异步请求
                return Task.Run(() => RestoreAsync(config)).GetAwaiter().GetResult();
            }
            catch (ApiException ex)
            {
                // 刷新令牌已失效（过期 / 被服务端注销）：清掉本地令牌，下次不再重试
                config.RefreshToken = string.Empty;
                config.Save();
                message = "自动登录已失效（" + ex.Message + "），请重新登录";
                return false;
            }
            catch (Exception ex)
            {
                // 网络不可达等临时故障：保留令牌，下次启动继续尝试自动登录
                message = "自动登录暂不可用（" + ex.Message + "），请手动登录";
                return false;
            }
        }

        /// <summary>
        /// 续期令牌并拉取门店，成功则写入 AppSession 与配置文件；任一步失败由调用方统一处理。
        /// </summary>
        private static async Task<bool> RestoreAsync(ClientConfig config)
        {
            var client = new ApiClient(config.ApiBaseUrl);
            var backend = new BackendApi(client);

            TokenPair tokens = await backend.RefreshAsync(config.RefreshToken);

            // 刷新令牌是一次性的，续期成功后立刻回写：后面拉门店失败也不会丢掉新令牌
            config.RefreshToken = tokens.RefreshToken;
            config.Save();

            List<ShopOption> shops = await backend.GetMyShopsAsync();
            if (shops == null || shops.Count == 0)
            {
                throw new ApiException(-1, "该账号未绑定任何打印门店");
            }

            AppSession.Api = client;
            AppSession.Backend = backend;
            AppSession.Tokens = tokens;
            AppSession.Username = config.Username;
            AppSession.Shops = shops;
            AppSession.CurrentShop = shops[0];
            return true;
        }
    }
}
