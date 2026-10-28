using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using eprint.Models;

namespace eprint.Services
{
    /// <summary>
    /// 云打印业务接口封装：登录、门店、设备登记、打印机上报与任务回执。
    /// 只负责拼请求与解包，不做界面提示。
    /// </summary>
    public sealed class BackendApi
    {
        private readonly ApiClient _client;

        /// <summary>
        /// 用已配置的 HTTP 客户端构造业务接口。
        /// </summary>
        public BackendApi(ApiClient client)
        {
            _client = client;
        }

        /// <summary>
        /// 获取图形验证码；登录接口必须携带 captchaId 与 captchaCode。
        /// </summary>
        public Task<CaptchaInfo> GetCaptchaAsync()
        {
            return _client.GetAsync<CaptchaInfo>("/api/v1/system/auth/captcha");
        }

        /// <summary>
        /// 后台账号密码登录，成功后把双令牌写回客户端并返回令牌对。
        /// </summary>
        public async Task<TokenPair> LoginAsync(string username, string password, string captchaId, string captchaCode)
        {
            var tokens = await _client.PostAsync<TokenPair>("/api/v1/system/auth/login", new
            {
                username = username,
                password = password,
                captchaId = captchaId,
                captchaCode = captchaCode
            });

            if (tokens == null)
            {
                throw new ApiException(-1, "登录响应缺少令牌");
            }

            _client.AccessToken = tokens.AccessToken;
            _client.RefreshToken = tokens.RefreshToken;

            return tokens;
        }

        /// <summary>
        /// 用刷新令牌换一对新令牌（服务端会作废旧刷新令牌），并写回客户端；
        /// 自动登录时调用，成功后必须持久化新的刷新令牌。
        /// </summary>
        public async Task<TokenPair> RefreshAsync(string refreshToken)
        {
            var tokens = await _client.PostAsync<TokenPair>("/api/v1/system/auth/refresh", new
            {
                refreshToken = refreshToken
            });

            if (tokens == null)
            {
                throw new ApiException(-1, "续期响应缺少令牌");
            }

            _client.AccessToken = tokens.AccessToken;
            _client.RefreshToken = tokens.RefreshToken;

            return tokens;
        }

        /// <summary>
        /// 查询当前账号绑定的门店，用于登录后选择设备归属。
        /// </summary>
        public Task<List<ShopOption>> GetMyShopsAsync()
        {
            return _client.GetAsync<List<ShopOption>>("/api/v1/print/device/shops");
        }

        /// <summary>
        /// 登记设备：同一机器码重复调用只更新归属与心跳，不会重复创建设备。
        /// </summary>
        public Task<DeviceInfo> RegisterDeviceAsync(long shopId, string code, string name, string version)
        {
            return _client.PostAsync<DeviceInfo>("/api/v1/print/device/devices/register", new
            {
                printShopId = shopId,
                code = code,
                name = name,
                version = version
            });
        }

        /// <summary>
        /// 查询设备已上报的打印机（含停用），用于回填勾选状态与能力。
        /// </summary>
        public Task<List<PrinterDto>> GetPrintersAsync(long deviceId)
        {
            return _client.GetAsync<List<PrinterDto>>(
                "/api/v1/print/device/devices/" + deviceId + "/printers");
        }

        /// <summary>
        /// 用当前勾选的完整集合覆盖设备打印机：新增的创建、已有的更新、未出现的停用。
        /// </summary>
        public Task<List<PrinterDto>> SyncPrintersAsync(long deviceId, IEnumerable<PrinterSyncItem> printers)
        {
            var payload = new
            {
                printers = printers.Select(item => new
                {
                    systemName = item.SystemName,
                    name = item.Name,
                    model = item.Model,
                    paperTypes = item.PaperTypes,
                    colorSupport = item.ColorSupport,
                    duplexSupport = item.DuplexSupport
                }).ToList()
            };

            return _client.PostAsync<List<PrinterDto>>(
                "/api/v1/print/device/devices/" + deviceId + "/printers/sync", payload);
        }

        /// <summary>
        /// 拉取分配给当前设备的待打印任务；limit 为单次上限（服务端最多 20）。
        /// </summary>
        public Task<List<TaskDto>> GetTasksAsync(long deviceId, int limit)
        {
            return _client.GetAsync<List<TaskDto>>(
                "/api/v1/print/device/devices/" + deviceId + "/tasks?limit=" + limit);
        }

        /// <summary>
        /// 回执任务状态：2 打印中 / 3 完成 / 4 失败；失败时带错误码与原因。
        /// </summary>
        public Task ReportTaskAsync(long deviceId, string taskNo, int status, int printedPages, string errorCode, string reason)
        {
            return _client.PostAsync<object>(
                "/api/v1/print/device/devices/" + deviceId + "/tasks/report",
                new
                {
                    taskNo = taskNo,
                    status = status,
                    printedPages = printedPages,
                    errorCode = errorCode,
                    reason = reason
                });
        }

        /// <summary>
        /// 退出登录：清理服务端令牌；失败不阻断客户端退出流程。
        /// </summary>
        public Task LogoutAsync()
        {
            return _client.PostAsync<object>("/api/v1/system/auth/logout", new
            {
                refreshToken = _client.RefreshToken
            });
        }
    }
}
