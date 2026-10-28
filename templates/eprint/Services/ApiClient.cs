using System;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading.Tasks;
using eprint.Models;

namespace eprint.Services
{
    /// <summary>
    /// 与服务端通信的 HTTP 客户端：统一处理 Bearer 令牌、JSON 序列化与业务错误码。
    /// 非 0 的 code 统一抛出 ApiException，由界面层决定如何提示。
    /// </summary>
    public sealed class ApiClient
    {
        private readonly HttpClient _http;

        /// <summary>
        /// 服务端根地址，不带结尾斜杠。
        /// </summary>
        public string BaseUrl { get; private set; }

        /// <summary>
        /// 当前访问令牌，为空表示未登录。
        /// </summary>
        public string AccessToken { get; set; }

        /// <summary>
        /// 刷新令牌，退出登录时随请求提交。
        /// </summary>
        public string RefreshToken { get; set; }

        /// <summary>
        /// 用服务端地址构造客户端；超时 30 秒，避免界面长时间无响应。
        /// </summary>
        public ApiClient(string baseUrl)
        {
            BaseUrl = (baseUrl ?? string.Empty).Trim().TrimEnd('/');
            _http = new HttpClient { Timeout = TimeSpan.FromSeconds(30) };
        }

        /// <summary>
        /// 发起 GET 请求并返回 data 部分。
        /// </summary>
        public Task<T> GetAsync<T>(string path)
        {
            return SendAsync<T>(HttpMethod.Get, path, null);
        }

        /// <summary>
        /// 发起 POST 请求并返回 data 部分；body 为 null 时不写请求体。
        /// </summary>
        public Task<T> PostAsync<T>(string path, object body)
        {
            return SendAsync<T>(HttpMethod.Post, path, body);
        }

        /// <summary>
        /// 发送请求并解析 { code, msg, data }：HTTP 401 提示重新登录，业务 code 非 0 抛 ApiException。
        /// </summary>
        private async Task<T> SendAsync<T>(HttpMethod method, string path, object body)
        {
            using (var request = new HttpRequestMessage(method, BaseUrl + path))
            {
                if (!string.IsNullOrEmpty(AccessToken))
                {
                    request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", AccessToken);
                }

                if (body != null)
                {
                    request.Content = new StringContent(Json.Serialize(body), Encoding.UTF8, "application/json");
                }

                using (HttpResponseMessage response = await _http.SendAsync(request).ConfigureAwait(false))
                {
                    string text = await response.Content.ReadAsStringAsync().ConfigureAwait(false);

                    if (response.StatusCode == HttpStatusCode.Unauthorized)
                    {
                        throw new ApiException(401, "登录已过期，请重新登录");
                    }
                    if (!response.IsSuccessStatusCode)
                    {
                        throw new ApiException((int)response.StatusCode, "服务端返回 HTTP " + (int)response.StatusCode);
                    }

                    ApiEnvelope<T> envelope;
                    try
                    {
                        envelope = Json.Deserialize<ApiEnvelope<T>>(text);
                    }
                    catch (Exception)
                    {
                        throw new ApiException(-1, "服务端响应格式不正确");
                    }

                    if (envelope == null)
                    {
                        throw new ApiException(-1, "服务端响应为空");
                    }
                    if (envelope.Code != 0)
                    {
                        throw new ApiException(envelope.Code, envelope.Msg);
                    }

                    return envelope.Data;
                }
            }
        }
    }

    /// <summary>
    /// 业务接口错误：Code 为服务端业务码或 HTTP 状态码，Message 可直接展示给用户。
    /// </summary>
    public sealed class ApiException : Exception
    {
        public int Code { get; private set; }

        public ApiException(int code, string message)
            : base(message)
        {
            Code = code;
        }
    }
}
