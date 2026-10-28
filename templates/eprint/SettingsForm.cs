using System;
using System.Drawing;
using System.Windows.Forms;
using eprint.Services;

namespace eprint
{
    /// <summary>
    /// 设置窗口：修改服务端地址并立即写入独立配置文件，下次启动直接生效。
    /// 只负责地址这一项，其他配置仍由 client.json 自行维护。
    /// </summary>
    public partial class SettingsForm : Form
    {
        /// <summary>
        /// 构造设置窗口，回填当前生效的服务端地址与配置文件路径。
        /// </summary>
        public SettingsForm()
        {
            InitializeComponent();
            txtServer.Text = AppConfig.ApiBaseUrl;
            lblPath.Text = "配置文件：" + ClientConfig.FilePath;
        }

        /// <summary>
        /// 保存：规范化地址（补 http://、去结尾斜杠）后写入独立配置文件；
        /// 地址为空或写盘失败时停留在本窗口并提示原因。
        /// </summary>
        private void btnSave_Click(object sender, EventArgs e)
        {
            string server = AppConfig.Normalize(txtServer.Text);
            if (server.Length == 0)
            {
                SetStatus("请填写服务端地址", true);
                txtServer.Focus();
                return;
            }

            ClientConfig config = ClientConfig.Current;
            string previous = config.ApiBaseUrl;
            config.ApiBaseUrl = server;
            if (!config.Save())
            {
                // 写盘失败就回滚内存值，避免本次运行用了一个并没保存下来的地址
                config.ApiBaseUrl = previous;
                SetStatus("配置保存失败，请检查目录写入权限", true);
                return;
            }

            SetStatus(string.Empty, false);
            DialogResult = DialogResult.OK;
            Close();
        }

        /// <summary>
        /// 取消：不做任何修改直接关闭。
        /// </summary>
        private void btnCancel_Click(object sender, EventArgs e)
        {
            DialogResult = DialogResult.Cancel;
            Close();
        }

        /// <summary>
        /// 更新窗口内提示文案；isError 为 true 时用红色，否则用灰色。
        /// </summary>
        private void SetStatus(string message, bool isError)
        {
            lblStatus.ForeColor = isError ? Color.Firebrick : Color.DimGray;
            lblStatus.Text = message ?? string.Empty;
        }
    }
}
