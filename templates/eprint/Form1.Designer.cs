namespace eprint
{
    partial class Form1
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }

            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle1 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle2 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle3 = new System.Windows.Forms.DataGridViewCellStyle();
            this.panelHeader = new System.Windows.Forms.Panel();
            this.panelAccent = new System.Windows.Forms.Panel();
            this.lblAppTitle = new System.Windows.Forms.Label();
            this.lblAppSubtitle = new System.Windows.Forms.Label();
            this.lblUser = new System.Windows.Forms.Label();
            this.lblDevice = new System.Windows.Forms.Label();
            this.btnLogout = new System.Windows.Forms.Button();
            this.panelHeaderLine = new System.Windows.Forms.Panel();
            this.panelBody = new System.Windows.Forms.Panel();
            this.panelGridCard = new System.Windows.Forms.Panel();
            this.panelGridBody = new System.Windows.Forms.Panel();
            this.gridPrinters = new System.Windows.Forms.DataGridView();
            this.colSelected = new System.Windows.Forms.DataGridViewCheckBoxColumn();
            this.colName = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.colModel = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.colA4 = new System.Windows.Forms.DataGridViewCheckBoxColumn();
            this.colA3 = new System.Windows.Forms.DataGridViewCheckBoxColumn();
            this.colColor = new System.Windows.Forms.DataGridViewCheckBoxColumn();
            this.colDuplex = new System.Windows.Forms.DataGridViewCheckBoxColumn();
            this.colDefault = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.lblEmpty = new System.Windows.Forms.Label();
            this.panelGridLine = new System.Windows.Forms.Panel();
            this.panelGridHeader = new System.Windows.Forms.Panel();
            this.lblGridTitle = new System.Windows.Forms.Label();
            this.lblGridHint = new System.Windows.Forms.Label();
            this.btnSelectAll = new System.Windows.Forms.Button();
            this.btnSelectNone = new System.Windows.Forms.Button();
            this.tblStats = new System.Windows.Forms.TableLayoutPanel();
            this.panelCardInstalled = new System.Windows.Forms.Panel();
            this.lblInstalledTitle = new System.Windows.Forms.Label();
            this.lblInstalledValue = new System.Windows.Forms.Label();
            this.lblInstalledHint = new System.Windows.Forms.Label();
            this.panelCardSelected = new System.Windows.Forms.Panel();
            this.lblSelectedTitle = new System.Windows.Forms.Label();
            this.lblSelectedValue = new System.Windows.Forms.Label();
            this.lblSelectedHint = new System.Windows.Forms.Label();
            this.panelCardAuto = new System.Windows.Forms.Panel();
            this.lblAutoTitle = new System.Windows.Forms.Label();
            this.lblAutoValue = new System.Windows.Forms.Label();
            this.lblAutoHint = new System.Windows.Forms.Label();
            this.panelToolbar = new System.Windows.Forms.Panel();
            this.lblShop = new System.Windows.Forms.Label();
            this.cboShop = new System.Windows.Forms.ComboBox();
            this.btnRefresh = new System.Windows.Forms.Button();
            this.btnSync = new System.Windows.Forms.Button();
            this.lblToolbarHint = new System.Windows.Forms.Label();
            this.panelStatus = new System.Windows.Forms.Panel();
            this.panelStatusLeft = new System.Windows.Forms.Panel();
            this.lblStatusIcon = new System.Windows.Forms.Label();
            this.lblStatus = new System.Windows.Forms.Label();
            this.lblFooter = new System.Windows.Forms.Label();
            this.panelStatusLine = new System.Windows.Forms.Panel();
            this.panelHeader.SuspendLayout();
            this.panelBody.SuspendLayout();
            this.panelGridCard.SuspendLayout();
            this.panelGridBody.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.gridPrinters)).BeginInit();
            this.panelGridHeader.SuspendLayout();
            this.tblStats.SuspendLayout();
            this.panelCardInstalled.SuspendLayout();
            this.panelCardSelected.SuspendLayout();
            this.panelCardAuto.SuspendLayout();
            this.panelToolbar.SuspendLayout();
            this.panelStatus.SuspendLayout();
            this.panelStatusLeft.SuspendLayout();
            this.SuspendLayout();
            // 
            // panelHeader
            // 
            this.panelHeader.BackColor = System.Drawing.Color.White;
            this.panelHeader.Controls.Add(this.panelAccent);
            this.panelHeader.Controls.Add(this.lblAppTitle);
            this.panelHeader.Controls.Add(this.lblAppSubtitle);
            this.panelHeader.Controls.Add(this.lblUser);
            this.panelHeader.Controls.Add(this.lblDevice);
            this.panelHeader.Controls.Add(this.btnLogout);
            this.panelHeader.Controls.Add(this.panelHeaderLine);
            this.panelHeader.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelHeader.Location = new System.Drawing.Point(0, 0);
            this.panelHeader.Name = "panelHeader";
            this.panelHeader.Size = new System.Drawing.Size(1120, 82);
            this.panelHeader.TabIndex = 0;
            // 
            // panelAccent
            // 
            this.panelAccent.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.panelAccent.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelAccent.Location = new System.Drawing.Point(0, 0);
            this.panelAccent.Name = "panelAccent";
            this.panelAccent.Size = new System.Drawing.Size(1120, 4);
            this.panelAccent.TabIndex = 0;
            // 
            // lblAppTitle
            // 
            this.lblAppTitle.AutoSize = true;
            this.lblAppTitle.Font = new System.Drawing.Font("Microsoft YaHei", 16F, System.Drawing.FontStyle.Bold);
            this.lblAppTitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.lblAppTitle.Location = new System.Drawing.Point(24, 20);
            this.lblAppTitle.Name = "lblAppTitle";
            this.lblAppTitle.Size = new System.Drawing.Size(164, 30);
            this.lblAppTitle.TabIndex = 1;
            this.lblAppTitle.Text = "易云印 云打印客户端";
            // 
            // lblAppSubtitle
            // 
            this.lblAppSubtitle.AutoSize = true;
            this.lblAppSubtitle.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblAppSubtitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.lblAppSubtitle.Location = new System.Drawing.Point(26, 52);
            this.lblAppSubtitle.Name = "lblAppSubtitle";
            this.lblAppSubtitle.Size = new System.Drawing.Size(220, 17);
            this.lblAppSubtitle.TabIndex = 2;
            this.lblAppSubtitle.Text = "打印机能力上报 · 自动接单打印";
            // 
            // lblUser
            // 
            this.lblUser.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.lblUser.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblUser.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.lblUser.Location = new System.Drawing.Point(756, 20);
            this.lblUser.Name = "lblUser";
            this.lblUser.Size = new System.Drawing.Size(220, 20);
            this.lblUser.TabIndex = 3;
            this.lblUser.Text = "账号：";
            this.lblUser.TextAlign = System.Drawing.ContentAlignment.MiddleRight;
            // 
            // lblDevice
            // 
            this.lblDevice.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.lblDevice.AutoEllipsis = true;
            this.lblDevice.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.lblDevice.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(152)))), ((int)(((byte)(162)))), ((int)(((byte)(179)))));
            this.lblDevice.Location = new System.Drawing.Point(500, 46);
            this.lblDevice.Name = "lblDevice";
            this.lblDevice.Size = new System.Drawing.Size(476, 20);
            this.lblDevice.TabIndex = 4;
            this.lblDevice.Text = "设备：未登记";
            this.lblDevice.TextAlign = System.Drawing.ContentAlignment.MiddleRight;
            // 
            // btnLogout
            // 
            this.btnLogout.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.btnLogout.BackColor = System.Drawing.Color.White;
            this.btnLogout.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnLogout.FlatAppearance.BorderColor = System.Drawing.Color.FromArgb(((int)(((byte)(208)))), ((int)(((byte)(213)))), ((int)(((byte)(221)))));
            this.btnLogout.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            this.btnLogout.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnLogout.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.btnLogout.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.btnLogout.Location = new System.Drawing.Point(988, 24);
            this.btnLogout.Name = "btnLogout";
            this.btnLogout.Size = new System.Drawing.Size(108, 34);
            this.btnLogout.TabIndex = 5;
            this.btnLogout.Text = "退出登录";
            this.btnLogout.UseVisualStyleBackColor = false;
            this.btnLogout.Click += new System.EventHandler(this.btnLogout_Click);
            // 
            // panelHeaderLine
            // 
            this.panelHeaderLine.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(236)))), ((int)(((byte)(240)))));
            this.panelHeaderLine.Dock = System.Windows.Forms.DockStyle.Bottom;
            this.panelHeaderLine.Location = new System.Drawing.Point(0, 81);
            this.panelHeaderLine.Name = "panelHeaderLine";
            this.panelHeaderLine.Size = new System.Drawing.Size(1120, 1);
            this.panelHeaderLine.TabIndex = 6;
            // 
            // panelBody
            // 
            this.panelBody.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(244)))), ((int)(((byte)(246)))), ((int)(((byte)(250)))));
            this.panelBody.Controls.Add(this.panelGridCard);
            this.panelBody.Controls.Add(this.tblStats);
            this.panelBody.Controls.Add(this.panelToolbar);
            this.panelBody.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelBody.Location = new System.Drawing.Point(0, 82);
            this.panelBody.Name = "panelBody";
            this.panelBody.Padding = new System.Windows.Forms.Padding(20, 16, 20, 12);
            this.panelBody.Size = new System.Drawing.Size(1120, 596);
            this.panelBody.TabIndex = 1;
            // 
            // panelGridCard
            // 
            this.panelGridCard.BackColor = System.Drawing.Color.White;
            this.panelGridCard.Controls.Add(this.panelGridBody);
            this.panelGridCard.Controls.Add(this.panelGridLine);
            this.panelGridCard.Controls.Add(this.panelGridHeader);
            this.panelGridCard.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelGridCard.Location = new System.Drawing.Point(20, 172);
            this.panelGridCard.Name = "panelGridCard";
            this.panelGridCard.Size = new System.Drawing.Size(1080, 412);
            this.panelGridCard.TabIndex = 2;
            // 
            // panelGridBody
            // 
            this.panelGridBody.BackColor = System.Drawing.Color.White;
            this.panelGridBody.Controls.Add(this.gridPrinters);
            this.panelGridBody.Controls.Add(this.lblEmpty);
            this.panelGridBody.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelGridBody.Location = new System.Drawing.Point(0, 59);
            this.panelGridBody.Name = "panelGridBody";
            this.panelGridBody.Padding = new System.Windows.Forms.Padding(16, 0, 16, 12);
            this.panelGridBody.Size = new System.Drawing.Size(1080, 353);
            this.panelGridBody.TabIndex = 1;
            // 
            // gridPrinters
            // 
            this.gridPrinters.AllowUserToAddRows = false;
            this.gridPrinters.AllowUserToDeleteRows = false;
            this.gridPrinters.AllowUserToOrderColumns = false;
            this.gridPrinters.AllowUserToResizeRows = false;
            this.gridPrinters.AutoSizeColumnsMode = System.Windows.Forms.DataGridViewAutoSizeColumnsMode.Fill;
            this.gridPrinters.BackgroundColor = System.Drawing.Color.White;
            this.gridPrinters.BorderStyle = System.Windows.Forms.BorderStyle.None;
            this.gridPrinters.CellBorderStyle = System.Windows.Forms.DataGridViewCellBorderStyle.SingleHorizontal;
            this.gridPrinters.ColumnHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle1.BackColor = System.Drawing.Color.White;
            dataGridViewCellStyle1.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(52)))), ((int)(((byte)(64)))), ((int)(((byte)(84)))));
            dataGridViewCellStyle1.Padding = new System.Windows.Forms.Padding(4, 0, 4, 0);
            dataGridViewCellStyle1.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(239)))), ((int)(((byte)(246)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle1.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(29)))), ((int)(((byte)(78)))), ((int)(((byte)(216)))));
            dataGridViewCellStyle2.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle2.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            dataGridViewCellStyle2.Font = new System.Drawing.Font("Microsoft YaHei", 9F, System.Drawing.FontStyle.Bold);
            dataGridViewCellStyle2.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            dataGridViewCellStyle2.Padding = new System.Windows.Forms.Padding(8, 0, 8, 0);
            dataGridViewCellStyle2.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            dataGridViewCellStyle2.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            dataGridViewCellStyle2.WrapMode = System.Windows.Forms.DataGridViewTriState.False;
            dataGridViewCellStyle3.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(252)))), ((int)(((byte)(253)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle3.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(239)))), ((int)(((byte)(246)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle3.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(29)))), ((int)(((byte)(78)))), ((int)(((byte)(216)))));
            this.gridPrinters.ColumnHeadersDefaultCellStyle = dataGridViewCellStyle2;
            this.gridPrinters.DefaultCellStyle = dataGridViewCellStyle1;
            this.gridPrinters.AlternatingRowsDefaultCellStyle = dataGridViewCellStyle3;
            this.gridPrinters.ColumnHeadersHeight = 42;
            this.gridPrinters.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.DisableResizing;
            this.gridPrinters.Columns.AddRange(new System.Windows.Forms.DataGridViewColumn[] {
            this.colSelected,
            this.colName,
            this.colModel,
            this.colA4,
            this.colA3,
            this.colColor,
            this.colDuplex,
            this.colDefault});
            this.gridPrinters.Dock = System.Windows.Forms.DockStyle.Fill;
            this.gridPrinters.EnableHeadersVisualStyles = false;
            this.gridPrinters.GridColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(236)))), ((int)(((byte)(240)))));
            this.gridPrinters.Location = new System.Drawing.Point(16, 0);
            this.gridPrinters.MultiSelect = false;
            this.gridPrinters.Name = "gridPrinters";
            this.gridPrinters.RowHeadersVisible = false;
            this.gridPrinters.RowTemplate.Height = 38;
            this.gridPrinters.SelectionMode = System.Windows.Forms.DataGridViewSelectionMode.FullRowSelect;
            this.gridPrinters.Size = new System.Drawing.Size(1048, 341);
            this.gridPrinters.TabIndex = 0;
            this.gridPrinters.CellValueChanged += new System.Windows.Forms.DataGridViewCellEventHandler(this.gridPrinters_CellValueChanged);
            this.gridPrinters.CurrentCellDirtyStateChanged += new System.EventHandler(this.gridPrinters_CurrentCellDirtyStateChanged);
            // 
            // colSelected
            // 
            this.colSelected.DataPropertyName = "Selected";
            this.colSelected.DefaultCellStyle.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colSelected.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.colSelected.HeaderCell.Style.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colSelected.HeaderText = "上报";
            this.colSelected.Name = "colSelected";
            this.colSelected.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colSelected.Width = 60;
            // 
            // colName
            // 
            this.colName.AutoSizeMode = System.Windows.Forms.DataGridViewAutoSizeColumnMode.Fill;
            this.colName.DataPropertyName = "Name";
            this.colName.DefaultCellStyle.Font = new System.Drawing.Font("Microsoft YaHei", 9F, System.Drawing.FontStyle.Bold);
            this.colName.DefaultCellStyle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.colName.HeaderText = "打印机";
            this.colName.MinimumWidth = 260;
            this.colName.Name = "colName";
            this.colName.ReadOnly = true;
            this.colName.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            // 
            // colModel
            // 
            this.colModel.DataPropertyName = "Model";
            this.colModel.DefaultCellStyle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.colModel.HeaderText = "型号";
            this.colModel.Name = "colModel";
            this.colModel.ReadOnly = true;
            this.colModel.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colModel.Width = 180;
            // 
            // colA4
            // 
            this.colA4.DataPropertyName = "A4";
            this.colA4.DefaultCellStyle.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colA4.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.colA4.HeaderCell.Style.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colA4.HeaderText = "A4";
            this.colA4.Name = "colA4";
            this.colA4.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colA4.Width = 64;
            // 
            // colA3
            // 
            this.colA3.DataPropertyName = "A3";
            this.colA3.DefaultCellStyle.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colA3.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.colA3.HeaderCell.Style.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colA3.HeaderText = "A3";
            this.colA3.Name = "colA3";
            this.colA3.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colA3.Width = 64;
            // 
            // colColor
            // 
            this.colColor.DataPropertyName = "Color";
            this.colColor.DefaultCellStyle.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colColor.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.colColor.HeaderCell.Style.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colColor.HeaderText = "彩色";
            this.colColor.Name = "colColor";
            this.colColor.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colColor.Width = 72;
            // 
            // colDuplex
            // 
            this.colDuplex.DataPropertyName = "Duplex";
            this.colDuplex.DefaultCellStyle.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colDuplex.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.colDuplex.HeaderCell.Style.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colDuplex.HeaderText = "双面";
            this.colDuplex.Name = "colDuplex";
            this.colDuplex.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colDuplex.Width = 72;
            // 
            // colDefault
            // 
            this.colDefault.DataPropertyName = "DefaultText";
            this.colDefault.DefaultCellStyle.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colDefault.DefaultCellStyle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.colDefault.HeaderCell.Style.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleCenter;
            this.colDefault.HeaderText = "默认";
            this.colDefault.Name = "colDefault";
            this.colDefault.ReadOnly = true;
            this.colDefault.SortMode = System.Windows.Forms.DataGridViewColumnSortMode.NotSortable;
            this.colDefault.Width = 72;
            // 
            // lblEmpty
            // 
            this.lblEmpty.BackColor = System.Drawing.Color.White;
            this.lblEmpty.Dock = System.Windows.Forms.DockStyle.Fill;
            this.lblEmpty.Font = new System.Drawing.Font("Microsoft YaHei", 10F);
            this.lblEmpty.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(152)))), ((int)(((byte)(162)))), ((int)(((byte)(179)))));
            this.lblEmpty.Location = new System.Drawing.Point(16, 0);
            this.lblEmpty.Name = "lblEmpty";
            this.lblEmpty.Size = new System.Drawing.Size(1048, 341);
            this.lblEmpty.TabIndex = 1;
            this.lblEmpty.Text = "正在读取本机打印机…";
            this.lblEmpty.TextAlign = System.Drawing.ContentAlignment.MiddleCenter;
            // 
            // panelGridLine
            // 
            this.panelGridLine.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(236)))), ((int)(((byte)(240)))));
            this.panelGridLine.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelGridLine.Location = new System.Drawing.Point(0, 58);
            this.panelGridLine.Name = "panelGridLine";
            this.panelGridLine.Size = new System.Drawing.Size(1080, 1);
            this.panelGridLine.TabIndex = 2;
            // 
            // panelGridHeader
            // 
            this.panelGridHeader.Controls.Add(this.lblGridTitle);
            this.panelGridHeader.Controls.Add(this.lblGridHint);
            this.panelGridHeader.Controls.Add(this.btnSelectAll);
            this.panelGridHeader.Controls.Add(this.btnSelectNone);
            this.panelGridHeader.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelGridHeader.Location = new System.Drawing.Point(0, 0);
            this.panelGridHeader.Name = "panelGridHeader";
            this.panelGridHeader.Size = new System.Drawing.Size(1080, 58);
            this.panelGridHeader.TabIndex = 0;
            // 
            // lblGridTitle
            // 
            this.lblGridTitle.AutoSize = true;
            this.lblGridTitle.Font = new System.Drawing.Font("Microsoft YaHei", 11F, System.Drawing.FontStyle.Bold);
            this.lblGridTitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.lblGridTitle.Location = new System.Drawing.Point(16, 10);
            this.lblGridTitle.Name = "lblGridTitle";
            this.lblGridTitle.Size = new System.Drawing.Size(84, 20);
            this.lblGridTitle.TabIndex = 0;
            this.lblGridTitle.Text = "打印机列表";
            // 
            // lblGridHint
            // 
            this.lblGridHint.AutoSize = true;
            this.lblGridHint.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.lblGridHint.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.lblGridHint.Location = new System.Drawing.Point(16, 33);
            this.lblGridHint.Name = "lblGridHint";
            this.lblGridHint.Size = new System.Drawing.Size(376, 17);
            this.lblGridHint.TabIndex = 1;
            this.lblGridHint.Text = "勾选需要上报的打印机，并确认纸张 / 彩色 / 双面能力";
            // 
            // btnSelectAll
            // 
            this.btnSelectAll.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.btnSelectAll.BackColor = System.Drawing.Color.White;
            this.btnSelectAll.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnSelectAll.FlatAppearance.BorderColor = System.Drawing.Color.FromArgb(((int)(((byte)(208)))), ((int)(((byte)(213)))), ((int)(((byte)(221)))));
            this.btnSelectAll.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            this.btnSelectAll.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnSelectAll.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.btnSelectAll.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.btnSelectAll.Location = new System.Drawing.Point(924, 15);
            this.btnSelectAll.Name = "btnSelectAll";
            this.btnSelectAll.Size = new System.Drawing.Size(64, 28);
            this.btnSelectAll.TabIndex = 2;
            this.btnSelectAll.Text = "全选";
            this.btnSelectAll.UseVisualStyleBackColor = false;
            this.btnSelectAll.Click += new System.EventHandler(this.btnSelectAll_Click);
            // 
            // btnSelectNone
            // 
            this.btnSelectNone.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.btnSelectNone.BackColor = System.Drawing.Color.White;
            this.btnSelectNone.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnSelectNone.FlatAppearance.BorderColor = System.Drawing.Color.FromArgb(((int)(((byte)(208)))), ((int)(((byte)(213)))), ((int)(((byte)(221)))));
            this.btnSelectNone.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            this.btnSelectNone.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnSelectNone.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.btnSelectNone.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.btnSelectNone.Location = new System.Drawing.Point(996, 15);
            this.btnSelectNone.Name = "btnSelectNone";
            this.btnSelectNone.Size = new System.Drawing.Size(64, 28);
            this.btnSelectNone.TabIndex = 3;
            this.btnSelectNone.Text = "全不选";
            this.btnSelectNone.UseVisualStyleBackColor = false;
            this.btnSelectNone.Click += new System.EventHandler(this.btnSelectNone_Click);
            // 
            // tblStats
            // 
            this.tblStats.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(244)))), ((int)(((byte)(246)))), ((int)(((byte)(250)))));
            this.tblStats.ColumnCount = 3;
            this.tblStats.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.333F));
            this.tblStats.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.333F));
            this.tblStats.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 33.334F));
            this.tblStats.Controls.Add(this.panelCardInstalled, 0, 0);
            this.tblStats.Controls.Add(this.panelCardSelected, 1, 0);
            this.tblStats.Controls.Add(this.panelCardAuto, 2, 0);
            this.tblStats.Dock = System.Windows.Forms.DockStyle.Top;
            this.tblStats.Location = new System.Drawing.Point(20, 80);
            this.tblStats.Name = "tblStats";
            this.tblStats.Padding = new System.Windows.Forms.Padding(0, 12, 0, 4);
            this.tblStats.RowCount = 1;
            this.tblStats.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 100F));
            this.tblStats.Size = new System.Drawing.Size(1080, 92);
            this.tblStats.TabIndex = 1;
            // 
            // panelCardInstalled
            // 
            this.panelCardInstalled.BackColor = System.Drawing.Color.White;
            this.panelCardInstalled.Controls.Add(this.lblInstalledTitle);
            this.panelCardInstalled.Controls.Add(this.lblInstalledValue);
            this.panelCardInstalled.Controls.Add(this.lblInstalledHint);
            this.panelCardInstalled.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelCardInstalled.Location = new System.Drawing.Point(0, 12);
            this.panelCardInstalled.Margin = new System.Windows.Forms.Padding(0, 0, 12, 0);
            this.panelCardInstalled.Name = "panelCardInstalled";
            this.panelCardInstalled.Size = new System.Drawing.Size(348, 76);
            this.panelCardInstalled.TabIndex = 0;
            // 
            // lblInstalledTitle
            // 
            this.lblInstalledTitle.AutoSize = true;
            this.lblInstalledTitle.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblInstalledTitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.lblInstalledTitle.Location = new System.Drawing.Point(16, 10);
            this.lblInstalledTitle.Name = "lblInstalledTitle";
            this.lblInstalledTitle.Size = new System.Drawing.Size(92, 17);
            this.lblInstalledTitle.TabIndex = 0;
            this.lblInstalledTitle.Text = "本机打印机";
            // 
            // lblInstalledValue
            // 
            this.lblInstalledValue.AutoSize = true;
            this.lblInstalledValue.Font = new System.Drawing.Font("Microsoft YaHei", 18F, System.Drawing.FontStyle.Bold);
            this.lblInstalledValue.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.lblInstalledValue.Location = new System.Drawing.Point(16, 31);
            this.lblInstalledValue.Name = "lblInstalledValue";
            this.lblInstalledValue.Size = new System.Drawing.Size(29, 31);
            this.lblInstalledValue.TabIndex = 1;
            this.lblInstalledValue.Text = "0";
            // 
            // lblInstalledHint
            // 
            this.lblInstalledHint.AutoSize = true;
            this.lblInstalledHint.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblInstalledHint.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(152)))), ((int)(((byte)(162)))), ((int)(((byte)(179)))));
            this.lblInstalledHint.Location = new System.Drawing.Point(52, 43);
            this.lblInstalledHint.Name = "lblInstalledHint";
            this.lblInstalledHint.Size = new System.Drawing.Size(80, 17);
            this.lblInstalledHint.TabIndex = 2;
            this.lblInstalledHint.Text = "台已安装";
            // 
            // panelCardSelected
            // 
            this.panelCardSelected.BackColor = System.Drawing.Color.White;
            this.panelCardSelected.Controls.Add(this.lblSelectedTitle);
            this.panelCardSelected.Controls.Add(this.lblSelectedValue);
            this.panelCardSelected.Controls.Add(this.lblSelectedHint);
            this.panelCardSelected.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelCardSelected.Location = new System.Drawing.Point(360, 12);
            this.panelCardSelected.Margin = new System.Windows.Forms.Padding(0, 0, 12, 0);
            this.panelCardSelected.Name = "panelCardSelected";
            this.panelCardSelected.Size = new System.Drawing.Size(348, 76);
            this.panelCardSelected.TabIndex = 1;
            // 
            // lblSelectedTitle
            // 
            this.lblSelectedTitle.AutoSize = true;
            this.lblSelectedTitle.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblSelectedTitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.lblSelectedTitle.Location = new System.Drawing.Point(16, 10);
            this.lblSelectedTitle.Name = "lblSelectedTitle";
            this.lblSelectedTitle.Size = new System.Drawing.Size(92, 17);
            this.lblSelectedTitle.TabIndex = 0;
            this.lblSelectedTitle.Text = "已勾选上报";
            // 
            // lblSelectedValue
            // 
            this.lblSelectedValue.AutoSize = true;
            this.lblSelectedValue.Font = new System.Drawing.Font("Microsoft YaHei", 18F, System.Drawing.FontStyle.Bold);
            this.lblSelectedValue.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(185)))), ((int)(((byte)(129)))));
            this.lblSelectedValue.Location = new System.Drawing.Point(16, 31);
            this.lblSelectedValue.Name = "lblSelectedValue";
            this.lblSelectedValue.Size = new System.Drawing.Size(29, 31);
            this.lblSelectedValue.TabIndex = 1;
            this.lblSelectedValue.Text = "0";
            // 
            // lblSelectedHint
            // 
            this.lblSelectedHint.AutoSize = true;
            this.lblSelectedHint.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblSelectedHint.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(152)))), ((int)(((byte)(162)))), ((int)(((byte)(179)))));
            this.lblSelectedHint.Location = new System.Drawing.Point(52, 43);
            this.lblSelectedHint.Name = "lblSelectedHint";
            this.lblSelectedHint.Size = new System.Drawing.Size(110, 17);
            this.lblSelectedHint.TabIndex = 2;
            this.lblSelectedHint.Text = "台已勾选，待上报";
            // 
            // panelCardAuto
            // 
            this.panelCardAuto.BackColor = System.Drawing.Color.White;
            this.panelCardAuto.Controls.Add(this.lblAutoTitle);
            this.panelCardAuto.Controls.Add(this.lblAutoValue);
            this.panelCardAuto.Controls.Add(this.lblAutoHint);
            this.panelCardAuto.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelCardAuto.Location = new System.Drawing.Point(720, 12);
            this.panelCardAuto.Margin = new System.Windows.Forms.Padding(0);
            this.panelCardAuto.Name = "panelCardAuto";
            this.panelCardAuto.Size = new System.Drawing.Size(360, 76);
            this.panelCardAuto.TabIndex = 2;
            // 
            // lblAutoTitle
            // 
            this.lblAutoTitle.AutoSize = true;
            this.lblAutoTitle.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblAutoTitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.lblAutoTitle.Location = new System.Drawing.Point(16, 10);
            this.lblAutoTitle.Name = "lblAutoTitle";
            this.lblAutoTitle.Size = new System.Drawing.Size(92, 17);
            this.lblAutoTitle.TabIndex = 0;
            this.lblAutoTitle.Text = "自动打印";
            // 
            // lblAutoValue
            // 
            this.lblAutoValue.AutoSize = true;
            this.lblAutoValue.Font = new System.Drawing.Font("Microsoft YaHei", 14F, System.Drawing.FontStyle.Bold);
            this.lblAutoValue.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(185)))), ((int)(((byte)(129)))));
            this.lblAutoValue.Location = new System.Drawing.Point(16, 34);
            this.lblAutoValue.Name = "lblAutoValue";
            this.lblAutoValue.Size = new System.Drawing.Size(69, 25);
            this.lblAutoValue.TabIndex = 1;
            this.lblAutoValue.Text = "已开启";
            // 
            // lblAutoHint
            // 
            this.lblAutoHint.AutoSize = true;
            this.lblAutoHint.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblAutoHint.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(152)))), ((int)(((byte)(162)))), ((int)(((byte)(179)))));
            this.lblAutoHint.Location = new System.Drawing.Point(94, 40);
            this.lblAutoHint.Name = "lblAutoHint";
            this.lblAutoHint.Size = new System.Drawing.Size(140, 17);
            this.lblAutoHint.TabIndex = 2;
            this.lblAutoHint.Text = "每 5 秒拉取打印任务";
            // 
            // panelToolbar
            // 
            this.panelToolbar.BackColor = System.Drawing.Color.White;
            this.panelToolbar.Controls.Add(this.lblShop);
            this.panelToolbar.Controls.Add(this.cboShop);
            this.panelToolbar.Controls.Add(this.btnRefresh);
            this.panelToolbar.Controls.Add(this.btnSync);
            this.panelToolbar.Controls.Add(this.lblToolbarHint);
            this.panelToolbar.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelToolbar.Location = new System.Drawing.Point(20, 16);
            this.panelToolbar.Name = "panelToolbar";
            this.panelToolbar.Size = new System.Drawing.Size(1080, 64);
            this.panelToolbar.TabIndex = 0;
            // 
            // lblShop
            // 
            this.lblShop.AutoSize = true;
            this.lblShop.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblShop.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.lblShop.Location = new System.Drawing.Point(18, 22);
            this.lblShop.Name = "lblShop";
            this.lblShop.Size = new System.Drawing.Size(56, 17);
            this.lblShop.TabIndex = 0;
            this.lblShop.Text = "所属门店";
            // 
            // cboShop
            // 
            this.cboShop.DropDownStyle = System.Windows.Forms.ComboBoxStyle.DropDownList;
            this.cboShop.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.cboShop.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.cboShop.FormattingEnabled = true;
            this.cboShop.Location = new System.Drawing.Point(86, 17);
            this.cboShop.Name = "cboShop";
            this.cboShop.Size = new System.Drawing.Size(260, 25);
            this.cboShop.TabIndex = 1;
            this.cboShop.SelectedIndexChanged += new System.EventHandler(this.cboShop_SelectedIndexChanged);
            // 
            // btnRefresh
            // 
            this.btnRefresh.BackColor = System.Drawing.Color.White;
            this.btnRefresh.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnRefresh.FlatAppearance.BorderColor = System.Drawing.Color.FromArgb(((int)(((byte)(208)))), ((int)(((byte)(213)))), ((int)(((byte)(221)))));
            this.btnRefresh.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            this.btnRefresh.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnRefresh.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.btnRefresh.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.btnRefresh.Location = new System.Drawing.Point(360, 14);
            this.btnRefresh.Name = "btnRefresh";
            this.btnRefresh.Size = new System.Drawing.Size(100, 34);
            this.btnRefresh.TabIndex = 2;
            this.btnRefresh.Text = "刷新列表";
            this.btnRefresh.UseVisualStyleBackColor = false;
            this.btnRefresh.Click += new System.EventHandler(this.btnRefresh_Click);
            // 
            // btnSync
            // 
            this.btnSync.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.btnSync.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnSync.FlatAppearance.BorderSize = 0;
            this.btnSync.FlatAppearance.MouseDownBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(30)))), ((int)(((byte)(64)))), ((int)(((byte)(175)))));
            this.btnSync.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(29)))), ((int)(((byte)(78)))), ((int)(((byte)(216)))));
            this.btnSync.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnSync.Font = new System.Drawing.Font("Microsoft YaHei", 9F, System.Drawing.FontStyle.Bold);
            this.btnSync.ForeColor = System.Drawing.Color.White;
            this.btnSync.Location = new System.Drawing.Point(472, 14);
            this.btnSync.Name = "btnSync";
            this.btnSync.Size = new System.Drawing.Size(112, 34);
            this.btnSync.TabIndex = 3;
            this.btnSync.Text = "上报选中";
            this.btnSync.UseVisualStyleBackColor = false;
            this.btnSync.Click += new System.EventHandler(this.btnSync_Click);
            // 
            // lblToolbarHint
            // 
            this.lblToolbarHint.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.lblToolbarHint.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.lblToolbarHint.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(152)))), ((int)(((byte)(162)))), ((int)(((byte)(179)))));
            this.lblToolbarHint.Location = new System.Drawing.Point(780, 22);
            this.lblToolbarHint.Name = "lblToolbarHint";
            this.lblToolbarHint.Size = new System.Drawing.Size(280, 20);
            this.lblToolbarHint.TabIndex = 4;
            this.lblToolbarHint.Text = "A4 / A3 · 彩色 · 双面能力可调整";
            this.lblToolbarHint.TextAlign = System.Drawing.ContentAlignment.MiddleRight;
            // 
            // panelStatus
            // 
            this.panelStatus.BackColor = System.Drawing.Color.White;
            this.panelStatus.Controls.Add(this.panelStatusLeft);
            this.panelStatus.Controls.Add(this.lblFooter);
            this.panelStatus.Controls.Add(this.panelStatusLine);
            this.panelStatus.Dock = System.Windows.Forms.DockStyle.Bottom;
            this.panelStatus.Location = new System.Drawing.Point(0, 678);
            this.panelStatus.Name = "panelStatus";
            this.panelStatus.Size = new System.Drawing.Size(1120, 42);
            this.panelStatus.TabIndex = 2;
            // 
            // panelStatusLeft
            // 
            this.panelStatusLeft.Controls.Add(this.lblStatusIcon);
            this.panelStatusLeft.Controls.Add(this.lblStatus);
            this.panelStatusLeft.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelStatusLeft.Location = new System.Drawing.Point(0, 1);
            this.panelStatusLeft.Name = "panelStatusLeft";
            this.panelStatusLeft.Size = new System.Drawing.Size(820, 41);
            this.panelStatusLeft.TabIndex = 0;
            // 
            // lblStatusIcon
            // 
            this.lblStatusIcon.AutoSize = true;
            this.lblStatusIcon.Font = new System.Drawing.Font("Microsoft YaHei", 10F);
            this.lblStatusIcon.ForeColor = System.Drawing.Color.Silver;
            this.lblStatusIcon.Location = new System.Drawing.Point(20, 13);
            this.lblStatusIcon.Name = "lblStatusIcon";
            this.lblStatusIcon.Size = new System.Drawing.Size(18, 19);
            this.lblStatusIcon.TabIndex = 0;
            this.lblStatusIcon.Text = "●";
            // 
            // lblStatus
            // 
            this.lblStatus.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.lblStatus.AutoEllipsis = true;
            this.lblStatus.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblStatus.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.lblStatus.Location = new System.Drawing.Point(40, 11);
            this.lblStatus.Name = "lblStatus";
            this.lblStatus.Size = new System.Drawing.Size(760, 20);
            this.lblStatus.TabIndex = 1;
            this.lblStatus.Text = "就绪";
            this.lblStatus.TextAlign = System.Drawing.ContentAlignment.MiddleLeft;
            // 
            // lblFooter
            // 
            this.lblFooter.Dock = System.Windows.Forms.DockStyle.Right;
            this.lblFooter.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.lblFooter.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(185)))), ((int)(((byte)(129)))));
            this.lblFooter.Location = new System.Drawing.Point(820, 1);
            this.lblFooter.Name = "lblFooter";
            this.lblFooter.Padding = new System.Windows.Forms.Padding(0, 0, 20, 0);
            this.lblFooter.Size = new System.Drawing.Size(300, 41);
            this.lblFooter.TabIndex = 1;
            this.lblFooter.Text = "● 自动接收打印任务中";
            this.lblFooter.TextAlign = System.Drawing.ContentAlignment.MiddleRight;
            // 
            // panelStatusLine
            // 
            this.panelStatusLine.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(236)))), ((int)(((byte)(240)))));
            this.panelStatusLine.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelStatusLine.Location = new System.Drawing.Point(0, 0);
            this.panelStatusLine.Name = "panelStatusLine";
            this.panelStatusLine.Size = new System.Drawing.Size(1120, 1);
            this.panelStatusLine.TabIndex = 2;
            // 
            // Form1
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(8F, 17F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(244)))), ((int)(((byte)(246)))), ((int)(((byte)(250)))));
            this.ClientSize = new System.Drawing.Size(1120, 720);
            this.Controls.Add(this.panelBody);
            this.Controls.Add(this.panelStatus);
            this.Controls.Add(this.panelHeader);
            this.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.MinimumSize = new System.Drawing.Size(1000, 640);
            this.Name = "Form1";
            this.StartPosition = System.Windows.Forms.FormStartPosition.CenterScreen;
            this.Text = "易云印 云打印客户端";
            this.Load += new System.EventHandler(this.Form1_Load);
            this.panelHeader.ResumeLayout(false);
            this.panelHeader.PerformLayout();
            this.panelBody.ResumeLayout(false);
            this.panelGridCard.ResumeLayout(false);
            this.panelGridBody.ResumeLayout(false);
            ((System.ComponentModel.ISupportInitialize)(this.gridPrinters)).EndInit();
            this.panelGridHeader.ResumeLayout(false);
            this.panelGridHeader.PerformLayout();
            this.tblStats.ResumeLayout(false);
            this.panelCardInstalled.ResumeLayout(false);
            this.panelCardInstalled.PerformLayout();
            this.panelCardSelected.ResumeLayout(false);
            this.panelCardSelected.PerformLayout();
            this.panelCardAuto.ResumeLayout(false);
            this.panelCardAuto.PerformLayout();
            this.panelToolbar.ResumeLayout(false);
            this.panelToolbar.PerformLayout();
            this.panelStatus.ResumeLayout(false);
            this.panelStatusLeft.ResumeLayout(false);
            this.panelStatusLeft.PerformLayout();
            this.ResumeLayout(false);
        }

        #endregion

        private System.Windows.Forms.Panel panelHeader;
        private System.Windows.Forms.Panel panelAccent;
        private System.Windows.Forms.Label lblAppTitle;
        private System.Windows.Forms.Label lblAppSubtitle;
        private System.Windows.Forms.Label lblUser;
        private System.Windows.Forms.Label lblDevice;
        private System.Windows.Forms.Button btnLogout;
        private System.Windows.Forms.Panel panelHeaderLine;
        private System.Windows.Forms.Panel panelBody;
        private System.Windows.Forms.Panel panelGridCard;
        private System.Windows.Forms.Panel panelGridBody;
        private System.Windows.Forms.DataGridView gridPrinters;
        private System.Windows.Forms.DataGridViewCheckBoxColumn colSelected;
        private System.Windows.Forms.DataGridViewTextBoxColumn colName;
        private System.Windows.Forms.DataGridViewTextBoxColumn colModel;
        private System.Windows.Forms.DataGridViewCheckBoxColumn colA4;
        private System.Windows.Forms.DataGridViewCheckBoxColumn colA3;
        private System.Windows.Forms.DataGridViewCheckBoxColumn colColor;
        private System.Windows.Forms.DataGridViewCheckBoxColumn colDuplex;
        private System.Windows.Forms.DataGridViewTextBoxColumn colDefault;
        private System.Windows.Forms.Label lblEmpty;
        private System.Windows.Forms.Panel panelGridLine;
        private System.Windows.Forms.Panel panelGridHeader;
        private System.Windows.Forms.Label lblGridTitle;
        private System.Windows.Forms.Label lblGridHint;
        private System.Windows.Forms.Button btnSelectAll;
        private System.Windows.Forms.Button btnSelectNone;
        private System.Windows.Forms.TableLayoutPanel tblStats;
        private System.Windows.Forms.Panel panelCardInstalled;
        private System.Windows.Forms.Label lblInstalledTitle;
        private System.Windows.Forms.Label lblInstalledValue;
        private System.Windows.Forms.Label lblInstalledHint;
        private System.Windows.Forms.Panel panelCardSelected;
        private System.Windows.Forms.Label lblSelectedTitle;
        private System.Windows.Forms.Label lblSelectedValue;
        private System.Windows.Forms.Label lblSelectedHint;
        private System.Windows.Forms.Panel panelCardAuto;
        private System.Windows.Forms.Label lblAutoTitle;
        private System.Windows.Forms.Label lblAutoValue;
        private System.Windows.Forms.Label lblAutoHint;
        private System.Windows.Forms.Panel panelToolbar;
        private System.Windows.Forms.Label lblShop;
        private System.Windows.Forms.ComboBox cboShop;
        private System.Windows.Forms.Button btnRefresh;
        private System.Windows.Forms.Button btnSync;
        private System.Windows.Forms.Label lblToolbarHint;
        private System.Windows.Forms.Panel panelStatus;
        private System.Windows.Forms.Panel panelStatusLeft;
        private System.Windows.Forms.Label lblStatusIcon;
        private System.Windows.Forms.Label lblStatus;
        private System.Windows.Forms.Label lblFooter;
        private System.Windows.Forms.Panel panelStatusLine;
    }
}
