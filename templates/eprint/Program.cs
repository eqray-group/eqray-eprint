using System;
using System.Net;
using System.Windows.Forms;
using eprint.Services;

namespace eprint
{
    /// <summary>
    /// 程序入口：优先按独立配置文件自动登录，失败再进登录窗体；退出登录后回到登录窗体，关闭窗体则退出进程。
    /// </summary>
    internal static class Program
    {
        /// <summary>
        /// 启动 WinForms 消息循环，并串起「自动登录 → 主界面 → 退出登录 → 登录窗体」的循环。
        /// </summary>
        [STAThread]
        private static void Main()
        {
            ServicePointManager.SecurityProtocol |= SecurityProtocolType.Tls12;

            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.SetUnhandledExceptionMode(UnhandledExceptionMode.CatchException);
            Application.ThreadException += OnThreadException;

            while (true)
            {
                string notice = null;
                bool restored = ClientConfig.Current.AutoLogin && AutoLoginService.TryRestore(out notice);

                if (!restored)
                {
                    using (var login = new LoginForm(notice))
                    {
                        if (login.ShowDialog() != DialogResult.OK)
                        {
                            return;
                        }
                    }
                }

                using (var main = new Form1())
                {
                    Application.Run(main);
                    if (!main.LogoutRequested)
                    {
                        return;
                    }
                }
            }
        }

        /// <summary>
        /// 兜底处理界面线程未捕获异常，避免程序直接闪退。
        /// </summary>
        private static void OnThreadException(object sender, System.Threading.ThreadExceptionEventArgs e)
        {
            MessageBox.Show("程序发生未处理异常：" + e.Exception.Message, "易云印",
                MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
    }
}
