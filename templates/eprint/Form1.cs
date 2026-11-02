using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;
using eprint.Models;
using eprint.Services;

namespace eprint
{
    /// <summary>
    /// 主窗体：登记设备、上报打印机能力，并在后台轮询任务、展示自动打印状态。
    /// </summary>
    public partial class Form1 : Form
    {
        private readonly BindingList<PrinterRow> _rows = new BindingList<PrinterRow>();
        private readonly HashSet<string> _taskHistory = new HashSet<string>(StringComparer.Ordinal);
        private readonly Timer _taskTimer = new Timer();
        /// <summary>任务日志最多保留的条数，超出后丢弃最早的记录。</summary>
        private const int MaxLogLines = 200;

        private bool _initializing;
        private bool _polling;
        private bool _loadingPrinters;
        /// <summary>勾选或能力改过但还没上报时为 true，用于区分「待上报 / 已上报」。</summary>
        private bool _pendingReport;
        private string _printerLoadError;
        /// <summary>上一条日志内容与时间，用于抑制短时间内重复的提示。</summary>
        private string _lastLogMessage;
        private DateTime _lastLogAt;

        /// <summary>
        /// 用户点击「退出登录」时为 true，Program 据此重新打开登录窗体。
        /// </summary>
        public bool LogoutRequested { get; private set; }

        /// <summary>
        /// 构造主窗体并绑定打印机列表。
        /// </summary>
        public Form1()
        {
            InitializeComponent();
            gridPrinters.AutoGenerateColumns = false;
            gridPrinters.DataSource = _rows;
            _rows.ListChanged += (sender, args) => UpdateSummary();

            // 每 5 秒拉一次待打印任务；同一时间只允许一个轮询 / 打印流程
            _taskTimer.Interval = 5000;
            _taskTimer.Tick += taskTimer_Tick;
            FormClosed += (sender, args) => _taskTimer.Stop();

            // 初始未开始轮询，先展示暂停状态，设备登记成功后切换为运行中
            SetAutoState(false);
            ShowEmptyState("正在读取本机打印机…");
        }

        /// <summary>
        /// 窗体加载：回填账号、门店下拉框，然后登记设备并拉取打印机。
        /// </summary>
        private async void Form1_Load(object sender, EventArgs e)
        {
            lblUser.Text = "账号：" + AppSession.Username;

            _initializing = true;
            cboShop.DataSource = AppSession.Shops;
            cboShop.DisplayMember = "DisplayName";
            cboShop.ValueMember = "Id";
            if (AppSession.CurrentShop != null)
            {
                cboShop.SelectedValue = AppSession.CurrentShop.Id;
            }
            else if (cboShop.Items.Count > 0)
            {
                cboShop.SelectedIndex = 0;
            }
            _initializing = false;

            await InitializeDeviceAsync();
            _taskTimer.Start();
            SetAutoState(AppSession.Device != null);

            AppendLog(
                "已登录：" + AppSession.Username + "，门店：" +
                (AppSession.CurrentShop != null ? AppSession.CurrentShop.DisplayName : "未选择"),
                0);
        }

        /// <summary>
        /// 切换门店：重新登记设备（设备归属随之更新），并刷新该设备的打印机。
        /// </summary>
        private async void cboShop_SelectedIndexChanged(object sender, EventArgs e)
        {
            if (_initializing)
            {
                return;
            }

            var shop = cboShop.SelectedItem as ShopOption;
            if (shop == null)
            {
                return;
            }
            if (AppSession.CurrentShop != null && AppSession.CurrentShop.Id == shop.Id)
            {
                return;
            }

            await InitializeDeviceAsync();
        }

        /// <summary>
        /// 登记设备并加载打印机；设备编码来自本机指纹，重复登录会复用同一设备记录。
        /// </summary>
        private async Task InitializeDeviceAsync()
        {
            var shop = cboShop.SelectedItem as ShopOption;
            if (shop == null)
            {
                SetStatus("请选择门店", true);
                if (_rows.Count == 0)
                {
                    ShowEmptyState("请选择门店后读取打印机");
                }
                return;
            }

            SetBusy(true, "正在登记设备…");
            try
            {
                DeviceInfo device = await AppSession.Backend.RegisterDeviceAsync(
                    shop.Id,
                    DeviceIdentity.GetDeviceCode(),
                    Environment.MachineName,
                    AppConfig.ClientVersion);

                AppSession.Device = device;
                AppSession.CurrentShop = shop;
                lblDevice.Text = "设备：" + device.Name + "（" + device.Code + "）";
                SetAutoState(true);

                AppendLog("设备已登记：" + device.Name + "（" + device.Code + "），门店：" + shop.DisplayName, 0);

                await LoadPrintersAsync();
            }
            catch (Exception ex)
            {
                SetStatus("设备登记失败：" + ex.Message, true);
                if (_rows.Count == 0)
                {
                    ShowEmptyState("设备登记失败\n请检查服务端连接后点击「刷新列表」");
                }
            }
            finally
            {
                SetBusy(false, null);
            }
        }

        /// <summary>
        /// 读取本机打印机，并与服务端已上报记录合并：服务端状态覆盖勾选与能力。
        /// </summary>
        private async Task LoadPrintersAsync()
        {
            if (AppSession.Device == null)
            {
                SetStatus("设备尚未登记，请重新选择门店后重试", true);
                if (_rows.Count == 0)
                {
                    ShowEmptyState("设备尚未登记\n请重新选择门店或点击「刷新列表」");
                }
                return;
            }

            SetBusy(true, "正在读取打印机…");
            _loadingPrinters = true;
            _printerLoadError = null;
            ShowEmptyState("正在读取本机打印机…");
            try
            {
                List<PrinterInfo> localPrinters = await Task.Run(() => PrinterService.GetInstalledPrinters());
                List<PrinterDto> reportedPrinters = await AppSession.Backend.GetPrintersAsync(AppSession.Device.Id);
                Dictionary<string, PrinterDto> reported = (reportedPrinters ?? new List<PrinterDto>())
                    .GroupBy(item => item.SystemName, StringComparer.OrdinalIgnoreCase)
                    .ToDictionary(group => group.Key, group => group.First(), StringComparer.OrdinalIgnoreCase);

                _rows.Clear();
                foreach (PrinterInfo printer in localPrinters)
                {
                    var row = new PrinterRow
                    {
                        Name = printer.Name,
                        SystemName = printer.SystemName,
                        Model = printer.Model,
                        A4 = printer.SupportsA4,
                        A3 = printer.SupportsA3,
                        Color = printer.SupportsColor,
                        Duplex = printer.SupportsDuplex,
                        IsDefault = printer.IsDefault
                    };

                    PrinterDto saved;
                    if (reported.TryGetValue(printer.SystemName, out saved))
                    {
                        row.Selected = saved.Status == 1;
                        row.A4 = HasPaper(saved.PaperTypes, "A4");
                        row.A3 = HasPaper(saved.PaperTypes, "A3");
                        row.Color = saved.ColorSupport == 1;
                        row.Duplex = saved.DuplexSupport == 1;
                    }

                    _rows.Add(row);
                }

                // 刚和服务端对齐过，重新加载后不存在待上报改动
                _pendingReport = false;
                AppendLog("已读取本机打印机 " + localPrinters.Count + " 台，与服务端同步完成", 0);

                // 绑定后默认会选中第一行，清除选中避免用户误以为该行已被勾选
                gridPrinters.ClearSelection();
                gridPrinters.CurrentCell = null;
            }
            catch (Exception ex)
            {
                _printerLoadError = "读取打印机失败\n请检查服务端连接后点击「刷新列表」";
                SetStatus("读取打印机失败：" + ex.Message, true);
            }
            finally
            {
                _loadingPrinters = false;
                UpdateSummary();
                SetBusy(false, null);
                if (_printerLoadError == null)
                {
                    SetStatus("就绪", false);
                }
            }
        }

        /// <summary>
        /// 「刷新列表」：重新读取本机打印机与服务端配置。
        /// </summary>
        private async void btnRefresh_Click(object sender, EventArgs e)
        {
            await LoadPrintersAsync();
        }

        /// <summary>
        /// 「全选」：勾选当前列表里的全部打印机，并立即刷新统计与表格显示。
        /// </summary>
        private void btnSelectAll_Click(object sender, EventArgs e)
        {
            if (_rows.Count == 0)
            {
                return;
            }

            foreach (PrinterRow row in _rows)
            {
                row.Selected = true;
            }

            _pendingReport = true;
            _rows.ResetBindings();
            gridPrinters.Refresh();
            UpdateSummary();
        }

        /// <summary>
        /// 「全不选」：清空全部勾选；再次上报会把该设备下的打印机全部停用。
        /// </summary>
        private void btnSelectNone_Click(object sender, EventArgs e)
        {
            if (_rows.Count == 0)
            {
                return;
            }

            foreach (PrinterRow row in _rows)
            {
                row.Selected = false;
            }

            _pendingReport = true;
            _rows.ResetBindings();
            gridPrinters.Refresh();
            UpdateSummary();
        }

        /// <summary>
        /// 「上报选中」：把勾选的打印机与能力覆盖到服务端；未勾选的会被服务端停用。
        /// </summary>
        private async void btnSync_Click(object sender, EventArgs e)
        {
            if (AppSession.Device == null)
            {
                SetStatus("设备尚未登记，请刷新列表后重试", true);
                return;
            }

            List<PrinterRow> selected = _rows.Where(row => row.Selected).ToList();
            foreach (PrinterRow row in selected)
            {
                if (!row.A4 && !row.A3)
                {
                    SetStatus("打印机【" + row.Name + "】至少需要勾选 A4 或 A3", true);
                    return;
                }
            }

            if (selected.Count == 0 &&
                MessageBox.Show(this,
                    "当前未勾选任何打印机，上报后会停用该设备下的全部打印机。是否继续？",
                    "易云印", MessageBoxButtons.YesNo, MessageBoxIcon.Warning) != DialogResult.Yes)
            {
                return;
            }

            SetBusy(true, "正在上报打印机…");
            try
            {
                List<PrinterSyncItem> items = selected.Select(row => new PrinterSyncItem
                {
                    SystemName = row.SystemName,
                    Name = row.Name,
                    Model = row.Model ?? string.Empty,
                    PaperTypes = BuildPaperTypes(row),
                    ColorSupport = row.Color ? 1 : 0,
                    DuplexSupport = row.Duplex ? 1 : 0
                }).ToList();

                await AppSession.Backend.SyncPrintersAsync(AppSession.Device.Id, items);

                // 上报成功后本地与服务端一致，卡片从「待上报」切到「已上报」
                _pendingReport = false;
                UpdateSummary();

                SetStatus("已上报 " + items.Count + " 台打印机", 1);
            }
            catch (Exception ex)
            {
                SetStatus("上报失败：" + ex.Message, true);
            }
            finally
            {
                SetBusy(false, null);
            }
        }

        /// <summary>
        /// 定时器入口：每 5 秒尝试拉取一次待打印任务，轮询失败不打断下一轮。
        /// </summary>
        private async void taskTimer_Tick(object sender, EventArgs e)
        {
            await PollTasksAsync();
        }

        /// <summary>
        /// 拉取分配给当前设备的任务并逐条打印。
        /// 先回执「打印中」再下载，避免同一任务被多轮轮询重复执行。
        /// </summary>
        private async Task PollTasksAsync()
        {
            if (_polling || AppSession.Device == null || AppSession.Backend == null)
            {
                return;
            }

            _polling = true;
            try
            {
                List<TaskDto> tasks = await AppSession.Backend.GetTasksAsync(AppSession.Device.Id, 5);
                if (tasks == null)
                {
                    return;
                }

                foreach (TaskDto task in tasks)
                {
                    if (task == null || string.IsNullOrWhiteSpace(task.TaskNo) || _taskHistory.Contains(task.TaskNo))
                    {
                        continue;
                    }

                    _taskHistory.Add(task.TaskNo);
                    bool reportedPrinting = false;
                    try
                    {
                        // 回执打印中成功后才开始下载与打印，失败则下一轮重试
                        await AppSession.Backend.ReportTaskAsync(
                            AppSession.Device.Id, task.TaskNo, 2, 0, string.Empty, string.Empty);
                        reportedPrinting = true;
                        SetStatus("正在打印任务 " + task.TaskNo + "…", false);

                        int printedPages = await ExecuteTaskAsync(task);
                        try
                        {
                            await AppSession.Backend.ReportTaskAsync(
                                AppSession.Device.Id, task.TaskNo, 3, printedPages, string.Empty, string.Empty);
                            SetStatus("任务 " + task.TaskNo + " 已打印完成", 1);
                        }
                        catch (Exception ex)
                        {
                            // 已实际打印成功，不能再回执失败，避免后台重复派单
                            SetStatus("任务 " + task.TaskNo + " 已打印，但完成回执失败：" + ex.Message, true);
                        }
                    }
                    catch (PrintException ex)
                    {
                        if (reportedPrinting)
                        {
                            await ReportFailureSafeAsync(task, ex.Code, ex.Message);
                        }
                        else
                        {
                            _taskHistory.Remove(task.TaskNo);
                            SetStatus("任务 " + task.TaskNo + " 无法开始打印：" + ex.Message, true);
                        }
                    }
                    catch (Exception ex)
                    {
                        if (reportedPrinting)
                        {
                            await ReportFailureSafeAsync(task, "UNKNOWN", ex.Message);
                        }
                        else
                        {
                            _taskHistory.Remove(task.TaskNo);
                            SetStatus("任务 " + task.TaskNo + " 无法开始打印：" + ex.Message, true);
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                // 网络异常不打断定时器，下一轮继续
                SetStatus("拉取打印任务失败：" + ex.Message, true);
            }
            finally
            {
                _polling = false;
            }
        }

        /// <summary>
        /// 下载任务 PDF 并调用 SumatraPDF 静默打印，返回打印页数（单份页数 × 份数）。
        /// 无论成功失败都清理临时文件。
        /// </summary>
        private static async Task<int> ExecuteTaskAsync(TaskDto task)
        {
            string directory = Path.Combine(Path.GetTempPath(), "eprint", "tasks");
            string pdfPath = null;
            try
            {
                pdfPath = await PrintService.DownloadAsync(task, directory);
                await PrintService.PrintAsync(task, pdfPath);

                int copies = task.Copies > 0 ? task.Copies : 1;
                int pages = task.PageCount > 0 ? task.PageCount : 1;
                return pages * copies;
            }
            finally
            {
                if (!string.IsNullOrEmpty(pdfPath))
                {
                    try
                    {
                        File.Delete(pdfPath);
                    }
                    catch
                    {
                        // 临时文件清理失败不影响打印结果回执
                    }
                }
            }
        }

        /// <summary>
        /// 尽力回执任务失败：原因截断到服务端 255 字符上限，回执失败只影响状态栏提示。
        /// </summary>
        private async Task ReportFailureSafeAsync(TaskDto task, string errorCode, string reason)
        {
            string message = string.IsNullOrWhiteSpace(reason) ? "打印失败" : reason.Trim();
            if (message.Length > 255)
            {
                message = message.Substring(0, 255);
            }
            if (string.IsNullOrWhiteSpace(errorCode))
            {
                errorCode = "UNKNOWN";
            }

            try
            {
                await AppSession.Backend.ReportTaskAsync(
                    AppSession.Device.Id, task.TaskNo, 4, 0, errorCode, message);
                SetStatus("任务 " + task.TaskNo + " 打印失败：" + message, true);
            }
            catch (Exception ex)
            {
                SetStatus("任务 " + task.TaskNo + " 打印失败，且失败回执未送达：" + ex.Message, true);
            }
        }

        /// <summary>
        /// 退出登录：清理服务端令牌后回到登录窗体，令牌清理失败不阻断退出。
        /// </summary>
        private async void btnLogout_Click(object sender, EventArgs e)
        {
            if (_polling)
            {
                MessageBox.Show(this, "正在处理打印任务，请等待任务完成后再退出", "易云印",
                    MessageBoxButtons.OK, MessageBoxIcon.Information);
                return;
            }

            if (MessageBox.Show(this, "确定退出登录吗？", "易云印",
                    MessageBoxButtons.YesNo, MessageBoxIcon.Question) != DialogResult.Yes)
            {
                return;
            }

            _taskTimer.Stop();
            SetAutoState(false);
            SetBusy(true, "正在退出…");
            try
            {
                if (AppSession.Backend != null)
                {
                    await AppSession.Backend.LogoutAsync();
                }
            }
            catch
            {
                // 令牌清理失败也允许本地退出，避免用户被卡在主界面
            }

            // 退出登录后清掉本地刷新令牌，避免下次启动还按旧身份自动登录
            ClientConfig config = ClientConfig.Current;
            config.RefreshToken = string.Empty;
            config.Save();

            LogoutRequested = true;
            AppSession.Reset();
            Close();
        }

        /// <summary>
        /// 复选框点击后立即提交，保证勾选状态马上写入行对象。
        /// </summary>
        private void gridPrinters_CurrentCellDirtyStateChanged(object sender, EventArgs e)
        {
            if (gridPrinters.IsCurrentCellDirty && gridPrinters.CurrentCell is DataGridViewCheckBoxCell)
            {
                gridPrinters.CommitEdit(DataGridViewDataErrorContexts.Commit);
            }
        }

        /// <summary>
        /// 勾选或能力列变化时标记「待上报」并刷新顶部统计；
        /// 「默认」列只是展示本机默认打印机，不参与上报，忽略。
        /// </summary>
        private void gridPrinters_CellValueChanged(object sender, DataGridViewCellEventArgs e)
        {
            if (e.RowIndex < 0 || e.ColumnIndex < 0)
            {
                return;
            }

            if (e.ColumnIndex == colDefault.Index)
            {
                return;
            }

            _pendingReport = true;
            UpdateSummary();
        }

        /// <summary>
        /// 更新顶部统计：本机总数与已勾选数量。
        /// </summary>
        private void UpdateSummary()
        {
            if (lblInstalledValue.IsDisposed)
            {
                return;
            }

            int total = _rows.Count;
            int selected = _rows.Count(row => row.Selected);
            lblInstalledValue.Text = total.ToString();
            lblSelectedValue.Text = selected.ToString();
            lblSelectedHint.Text = BuildSelectedHint(selected);
            lblSelectedHint.ForeColor = selected > 0 && !_pendingReport
                ? Color.FromArgb(16, 185, 129)
                : Color.FromArgb(152, 162, 179);

            if (total == 0)
            {
                string message;
                if (_loadingPrinters)
                {
                    message = "正在读取本机打印机…";
                }
                else if (!string.IsNullOrEmpty(_printerLoadError))
                {
                    message = _printerLoadError;
                }
                else
                {
                    message = "本机暂未检测到打印机\n请安装打印机后点击「刷新列表」";
                }

                ShowEmptyState(message);
            }
            else
            {
                HideEmptyState();
            }
        }

        /// <summary>
        /// 生成「已勾选上报」卡片的副标题：还有改动未上报显示「待上报」，
        /// 已和服务端一致显示「已上报」。参数 selected 为当前勾选数量。
        /// </summary>
        private string BuildSelectedHint(int selected)
        {
            if (selected == 0)
            {
                return _pendingReport ? "已全部取消，待上报" : "尚未勾选打印机";
            }

            return _pendingReport ? "台已勾选，待上报" : "台已勾选，已上报";
        }

        /// <summary>
        /// 显示表格上方的空状态提示，并把它提到最上层遮住空表格。
        /// </summary>
        private void ShowEmptyState(string message)
        {
            lblEmpty.Text = message;
            lblEmpty.Visible = true;
            gridPrinters.Visible = false;
            lblEmpty.BringToFront();
        }

        /// <summary>
        /// 有打印机数据时隐藏空状态提示并恢复表格。
        /// </summary>
        private void HideEmptyState()
        {
            lblEmpty.Visible = false;
            gridPrinters.Visible = true;
        }

        /// <summary>
        /// 更新自动打印状态卡片与底部状态栏，运行中显示绿色，停止时显示灰色。
        /// </summary>
        private void SetAutoState(bool running)
        {
            Color green = Color.FromArgb(16, 185, 129);
            Color gray = Color.FromArgb(152, 162, 179);

            lblAutoValue.Text = running ? "已开启" : "已暂停";
            lblAutoValue.ForeColor = running ? green : gray;
            lblAutoHint.Text = running ? "每 5 秒拉取打印任务" : "等待登录 / 选择门店";
            lblFooter.Text = running ? "● 自动接收打印任务中" : "● 自动打印未开启";
            lblFooter.ForeColor = running ? green : gray;
        }

        /// <summary>
        /// 按行拼出纸张能力：A4 / A3，逗号分隔。
        /// </summary>
        private static string BuildPaperTypes(PrinterRow row)
        {
            var papers = new List<string>(2);
            if (row.A4)
            {
                papers.Add("A4");
            }
            if (row.A3)
            {
                papers.Add("A3");
            }

            return string.Join(",", papers);
        }

        /// <summary>
        /// 判断逗号分隔的纸张能力是否包含目标纸张。
        /// </summary>
        private static bool HasPaper(string paperTypes, string target)
        {
            if (string.IsNullOrWhiteSpace(paperTypes))
            {
                return false;
            }

            return paperTypes.Split(',')
                .Any(item => string.Equals(item.Trim(), target, StringComparison.OrdinalIgnoreCase));
        }

        /// <summary>
        /// 请求期间禁用操作按钮，避免重复提交。
        /// </summary>
        private void SetBusy(bool busy, string message)
        {
            UseWaitCursor = busy;
            btnRefresh.Enabled = !busy;
            btnSync.Enabled = !busy;
            btnLogout.Enabled = !busy;
            btnSelectAll.Enabled = !busy;
            btnSelectNone.Enabled = !busy;
            cboShop.Enabled = !busy && cboShop.Items.Count > 1;

            if (message != null)
            {
                SetStatus(message, false);
            }
        }

        /// <summary>
        /// 更新底部状态栏：isError 为 true 时文字与圆点转红，同时写一条错误日志。
        /// </summary>
        private void SetStatus(string message, bool isError)
        {
            SetStatus(message, isError ? 2 : 0);
        }

        /// <summary>
        /// 更新底部状态栏并按级别写任务日志（0 普通 / 1 成功 / 2 错误）；
        /// message 为空时恢复「就绪」。
        /// </summary>
        private void SetStatus(string message, int logLevel)
        {
            if (string.IsNullOrEmpty(message))
            {
                message = "就绪";
            }

            bool isError = logLevel == 2;
            lblStatus.Text = message;
            lblStatus.ForeColor = isError
                ? Color.Firebrick
                : Color.FromArgb(71, 84, 103);
            panelStatusDot.BackColor = isError
                ? Color.Firebrick
                : (message == "就绪" ? Color.Silver : Color.FromArgb(37, 99, 235));

            AppendLog(message, logLevel);
        }

        /// <summary>
        /// 追加一条任务日志：过滤「就绪」与「正在…」这类临时提示，
        /// 相同内容一分钟内只记一次，避免轮询失败时刷屏。
        /// 参数 message 为日志内容，level 0 普通 / 1 成功 / 2 错误。
        /// </summary>
        private void AppendLog(string message, int level)
        {
            if (txtLog == null || txtLog.IsDisposed || string.IsNullOrEmpty(message))
            {
                return;
            }
            if (message == "就绪" || message.StartsWith("正在"))
            {
                return;
            }
            if (message == _lastLogMessage && (DateTime.Now - _lastLogAt).TotalSeconds < 60)
            {
                return;
            }

            _lastLogMessage = message;
            _lastLogAt = DateTime.Now;

            Color color = level == 2
                ? Color.Firebrick
                : (level == 1 ? Color.FromArgb(16, 185, 129) : Color.FromArgb(71, 84, 103));

            // 追加前判断是否停在底部：用户翻看历史日志时不强行滚到底
            bool follow = txtLog.SelectionStart >= txtLog.TextLength - 1;

            txtLog.SelectionStart = txtLog.TextLength;
            txtLog.SelectionLength = 0;
            txtLog.SelectionColor = color;
            txtLog.AppendText(DateTime.Now.ToString("HH:mm:ss") + "  " + message + Environment.NewLine);

            // 超出上限时从最早的行开始删：按 '\n' 定位行尾，逐行删除可以保留已有日志的着色
            while (txtLog.Lines.Length > MaxLogLines)
            {
                int firstLineEnd = txtLog.Text.IndexOf('\n');
                if (firstLineEnd < 0)
                {
                    break;
                }

                txtLog.Select(0, firstLineEnd + 1);
                txtLog.SelectedText = string.Empty;
            }

            txtLog.SelectionColor = txtLog.ForeColor;
            if (follow)
            {
                txtLog.SelectionStart = txtLog.TextLength;
                txtLog.ScrollToCaret();
            }
        }

        /// <summary>
        /// 「清空」：只清掉界面上显示的任务日志，不影响服务端与正在打印的任务。
        /// </summary>
        private void btnClearLog_Click(object sender, EventArgs e)
        {
            txtLog.Clear();
            _lastLogMessage = null;
        }

        /// <summary>
        /// 把状态圆点重绘成实心圆：Panel 本身是方块，自绘才能和右侧文字对齐且边缘平滑。
        /// </summary>
        private void panelStatusDot_Paint(object sender, PaintEventArgs e)
        {
            e.Graphics.SmoothingMode = System.Drawing.Drawing2D.SmoothingMode.AntiAlias;
            using (var brush = new SolidBrush(panelStatusDot.BackColor))
            {
                e.Graphics.FillEllipse(brush, 0, 0, panelStatusDot.Width - 1, panelStatusDot.Height - 1);
            }
        }
    }
}
