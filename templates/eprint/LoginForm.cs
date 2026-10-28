using System;
using System.Drawing;
using System.IO;
using System.Threading.Tasks;
using System.Windows.Forms;
using eprint.Models;
using eprint.Services;

namespace eprint
{
    /// <summary>
    /// 登录窗体：填写账号密码与图形验证码，登录成功后加载当前账号的门店。
    /// 服务端地址不在本窗体填写，统一走右上角「设置」，取值为独立配置文件里的值。
    /// </summary>
    public partial class LoginForm : Form
    {
        private CaptchaInfo _captcha;
        private readonly string _notice;

        /// <summary>
        /// 无启动提示的登录窗体：退出登录后重新登录、以及设计器实例化时使用。
        /// </summary>
        public LoginForm()
            : this(null)
        {
        }

        /// <summary>
        /// 构造登录窗体：记住启动提示，服务端地址在加载时从配置文件读取。
        /// notice 为自动登录失败等提示文案，显示在状态栏；为空表示无提示。
        /// </summary>
        public LoginForm(string notice)
        {
            InitializeComponent();
            _notice = notice;
        }

        /// <summary>
        /// 窗体加载：先回填上次的账号 / 密码与自动登录勾选，再获取验证码；
        /// 自动登录失败的提示放在最后显示，避免被验证码刷新的状态覆盖。
        /// </summary>
        private async void LoginForm_Load(object sender, EventArgs e)
        {
            ApplySavedSettings();
            await RefreshCaptchaAsync();

            if (!string.IsNullOrEmpty(_notice))
            {
                SetStatus(_notice, true);
            }
        }

        /// <summary>
        /// 点击验证码图片时换一张，同时清空已输入的验证码。
        /// </summary>
        private async void picCaptcha_Click(object sender, EventArgs e)
        {
            await RefreshCaptchaAsync();
        }

        /// <summary>
        /// 「看不清？换一张」链接：重新获取验证码。
        /// </summary>
        private async void lnkRefreshCaptcha_LinkClicked(object sender, LinkLabelLinkClickedEventArgs e)
        {
            await RefreshCaptchaAsync();
        }

        /// <summary>
        /// 「设置」齿轮：打开设置窗口修改服务端地址；保存后按新地址重新获取验证码，
        /// 避免拿旧地址的验证码去提交新地址。
        /// </summary>
        private async void btnSettings_Click(object sender, EventArgs e)
        {
            using (var form = new SettingsForm())
            {
                if (form.ShowDialog(this) != DialogResult.OK)
                {
                    return;
                }
            }

            await RefreshCaptchaAsync();

            if (_captcha != null)
            {
                SetStatus("服务端地址已更新，请重新输入验证码", false);
            }
        }

        /// <summary>
        /// 登录：校验输入 → 调登录接口 → 查询绑定门店 → 成功则关闭登录窗体进入主界面。
        /// 任一环节失败都会刷新验证码，避免验证码已被服务端消费后继续使用。
        /// </summary>
        private async void btnLogin_Click(object sender, EventArgs e)
        {
            string server = AppConfig.ApiBaseUrl;
            if (server.Length == 0)
            {
                SetStatus("未配置服务端地址，请点击右上角设置填写", true);
                return;
            }

            if (AppSession.Api == null ||
                !string.Equals(server, AppSession.Api.BaseUrl, StringComparison.OrdinalIgnoreCase))
            {
                await RefreshCaptchaAsync();
                SetStatus("服务端地址已变更，请重新输入验证码", true);
                return;
            }

            string username = txtUsername.Text.Trim();
            string password = txtPassword.Text;
            string captchaCode = txtCaptcha.Text.Trim();
            if (username.Length == 0) { SetStatus("请输入账号", true); txtUsername.Focus(); return; }
            if (password.Length == 0) { SetStatus("请输入密码", true); txtPassword.Focus(); return; }
            if (captchaCode.Length == 0) { SetStatus("请输入验证码", true); txtCaptcha.Focus(); return; }
            if (_captcha == null) { SetStatus("验证码未加载，请点击图片刷新", true); return; }

            SetBusy(true, "正在登录…");
            try
            {
                TokenPair tokens = await AppSession.Backend.LoginAsync(username, password, _captcha.CaptchaId, captchaCode);
                var shops = await AppSession.Backend.GetMyShopsAsync();
                if (shops == null || shops.Count == 0)
                {
                    // 先换验证码再提示，避免 RefreshCaptchaAsync 把错误文案清空
                    await RefreshCaptchaAsync();
                    SetStatus("该账号未绑定任何打印门店，请联系管理员", true);
                    return;
                }

                AppSession.Tokens = tokens;
                AppSession.Username = username;
                AppSession.Shops = shops;
                AppSession.CurrentShop = shops[0];
                AppConfig.ApiBaseUrl = server;
                PersistLoginSettings(server, username, password, tokens);

                DialogResult = DialogResult.OK;
                Close();
            }
            catch (ApiException ex)
            {
                string message = "登录失败：" + ex.Message;
                await RefreshCaptchaAsync();
                SetStatus(message, true);
            }
            catch (Exception ex)
            {
                string message = "登录失败：" + ex.Message;
                await RefreshCaptchaAsync();
                SetStatus(message, true);
            }
            finally
            {
                SetBusy(false, null);
            }
        }

        /// <summary>
        /// 退出登录窗体，不进入主界面。
        /// </summary>
        private void btnCancel_Click(object sender, EventArgs e)
        {
            DialogResult = DialogResult.Cancel;
            Close();
        }

        /// <summary>
        /// 按配置文件里的服务端地址重新获取验证码，并重建 HTTP 客户端与业务接口。
        /// </summary>
        private async Task RefreshCaptchaAsync()
        {
            string server = AppConfig.ApiBaseUrl;
            if (server.Length == 0)
            {
                SetStatus("未配置服务端地址，请点击右上角设置填写", true);
                return;
            }

            var client = new ApiClient(server);
            var backend = new BackendApi(client);
            SetBusy(true, "正在获取验证码…");

            try
            {
                CaptchaInfo captcha = await backend.GetCaptchaAsync();
                _captcha = captcha;
                SetCaptchaImage(captcha.Image);
                txtCaptcha.Clear();

                AppSession.Api = client;
                AppSession.Backend = backend;
                SetStatus(string.Empty, false);
            }
            catch (Exception ex)
            {
                _captcha = null;
                SetCaptchaImage(null);
                SetStatus("验证码加载失败：" + ex.Message, true);
            }
            finally
            {
                SetBusy(false, null);
            }
        }

        /// <summary>
        /// 把 base64 图片显示到验证码框；传入空值时只清空图片。
        /// </summary>
        private void SetCaptchaImage(string base64)
        {
            Image oldImage = picCaptcha.Image;
            picCaptcha.Image = null;
            if (oldImage != null)
            {
                oldImage.Dispose();
            }

            if (string.IsNullOrWhiteSpace(base64))
            {
                return;
            }

            int comma = base64.IndexOf(',');
            if (base64.StartsWith("data:", StringComparison.OrdinalIgnoreCase) && comma >= 0)
            {
                base64 = base64.Substring(comma + 1);
            }

            using (var stream = new MemoryStream(Convert.FromBase64String(base64)))
            using (Image image = Image.FromStream(stream))
            {
                picCaptcha.Image = new Bitmap(image);
            }
        }

        /// <summary>
        /// 用独立配置回填上次的账号、密码与自动登录勾选，减少重复输入。
        /// 密码在本机用 DPAPI 解密，解密失败（换用户 / 换机器）时按空值处理。
        /// </summary>
        private void ApplySavedSettings()
        {
            ClientConfig config = ClientConfig.Current;
            txtUsername.Text = config.Username ?? string.Empty;
            txtPassword.Text = config.Password ?? string.Empty;
            chkAutoLogin.Checked = config.AutoLogin;
        }

        /// <summary>
        /// 登录成功后持久化配置：服务器地址与账号总是记住；
        /// 勾选自动登录时额外保存密码与刷新令牌，供下次启动免输入直接进入主界面。
        /// </summary>
        private void PersistLoginSettings(string server, string username, string password, TokenPair tokens)
        {
            ClientConfig config = ClientConfig.Current;
            config.ApiBaseUrl = server;
            config.Username = username;
            config.AutoLogin = chkAutoLogin.Checked;

            if (chkAutoLogin.Checked)
            {
                config.Password = password;
                config.RefreshToken = tokens.RefreshToken;
            }
            else
            {
                // 未勾选自动登录：清掉旧凭据，避免下次启动拿旧令牌直接进入
                config.Password = string.Empty;
                config.RefreshToken = string.Empty;
            }

            config.Save();
        }

        /// <summary>
        /// 登录过程中禁用输入与按钮，切换等待光标。
        /// </summary>
        private void SetBusy(bool busy, string message)
        {
            UseWaitCursor = busy;
            txtUsername.Enabled = !busy;
            txtPassword.Enabled = !busy;
            txtCaptcha.Enabled = !busy;
            picCaptcha.Enabled = !busy;
            chkAutoLogin.Enabled = !busy;
            btnSettings.Enabled = !busy;
            btnLogin.Enabled = !busy;
            btnCancel.Enabled = !busy;

            if (message != null)
            {
                SetStatus(message, false);
            }
        }

        /// <summary>
        /// 更新底部状态文案与颜色。
        /// </summary>
        private void SetStatus(string message, bool isError)
        {
            lblStatus.ForeColor = isError ? Color.Firebrick : Color.DimGray;
            lblStatus.Text = message ?? string.Empty;
        }
    }
}
