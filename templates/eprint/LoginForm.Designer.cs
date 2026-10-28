namespace eprint
{
    partial class LoginForm
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
            this.panelForm = new System.Windows.Forms.Panel();
            this.lblLoginTitle = new System.Windows.Forms.Label();
            this.lblLoginSubtitle = new System.Windows.Forms.Label();
            this.lblUsername = new System.Windows.Forms.Label();
            this.txtUsername = new System.Windows.Forms.TextBox();
            this.lblPassword = new System.Windows.Forms.Label();
            this.txtPassword = new System.Windows.Forms.TextBox();
            this.lblCaptcha = new System.Windows.Forms.Label();
            this.txtCaptcha = new System.Windows.Forms.TextBox();
            this.picCaptcha = new System.Windows.Forms.PictureBox();
            this.lnkRefreshCaptcha = new System.Windows.Forms.LinkLabel();
            this.chkAutoLogin = new System.Windows.Forms.CheckBox();
            this.btnLogin = new System.Windows.Forms.Button();
            this.btnCancel = new System.Windows.Forms.Button();
            this.lblStatus = new System.Windows.Forms.Label();
            this.btnSettings = new System.Windows.Forms.Button();
            this.panelBrand = new System.Windows.Forms.Panel();
            this.lblBrandTitle = new System.Windows.Forms.Label();
            this.lblBrandSubtitle = new System.Windows.Forms.Label();
            this.lblBrandDesc = new System.Windows.Forms.Label();
            this.lblBrandFeatures = new System.Windows.Forms.Label();
            this.lblBrandVersion = new System.Windows.Forms.Label();
            this.panelBrandAccent = new System.Windows.Forms.Panel();
            this.panelForm.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.picCaptcha)).BeginInit();
            this.panelBrand.SuspendLayout();
            this.SuspendLayout();
            // 
            // panelForm
            // 
            this.panelForm.BackColor = System.Drawing.Color.White;
            this.panelForm.Controls.Add(this.lblLoginTitle);
            this.panelForm.Controls.Add(this.lblLoginSubtitle);
            this.panelForm.Controls.Add(this.lblUsername);
            this.panelForm.Controls.Add(this.txtUsername);
            this.panelForm.Controls.Add(this.lblPassword);
            this.panelForm.Controls.Add(this.txtPassword);
            this.panelForm.Controls.Add(this.lblCaptcha);
            this.panelForm.Controls.Add(this.txtCaptcha);
            this.panelForm.Controls.Add(this.picCaptcha);
            this.panelForm.Controls.Add(this.lnkRefreshCaptcha);
            this.panelForm.Controls.Add(this.chkAutoLogin);
            this.panelForm.Controls.Add(this.btnLogin);
            this.panelForm.Controls.Add(this.btnCancel);
            this.panelForm.Controls.Add(this.lblStatus);
            this.panelForm.Controls.Add(this.btnSettings);
            this.panelForm.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panelForm.Location = new System.Drawing.Point(330, 0);
            this.panelForm.Name = "panelForm";
            this.panelForm.Size = new System.Drawing.Size(510, 520);
            this.panelForm.TabIndex = 1;
            // 
            // lblLoginTitle
            // 
            this.lblLoginTitle.AutoSize = true;
            this.lblLoginTitle.Font = new System.Drawing.Font("Microsoft YaHei", 18F, System.Drawing.FontStyle.Bold);
            this.lblLoginTitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.lblLoginTitle.Location = new System.Drawing.Point(44, 42);
            this.lblLoginTitle.Name = "lblLoginTitle";
            this.lblLoginTitle.Size = new System.Drawing.Size(126, 31);
            this.lblLoginTitle.TabIndex = 0;
            this.lblLoginTitle.Text = "登录客户端";
            // 
            // lblLoginSubtitle
            // 
            this.lblLoginSubtitle.AutoSize = true;
            this.lblLoginSubtitle.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblLoginSubtitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.lblLoginSubtitle.Location = new System.Drawing.Point(46, 80);
            this.lblLoginSubtitle.Name = "lblLoginSubtitle";
            this.lblLoginSubtitle.Size = new System.Drawing.Size(260, 17);
            this.lblLoginSubtitle.TabIndex = 1;
            this.lblLoginSubtitle.Text = "使用后台账号登录，登录后自动接收打印任务";
            // 
            // btnSettings
            // 
            this.btnSettings.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Right)));
            this.btnSettings.BackColor = System.Drawing.Color.White;
            this.btnSettings.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnSettings.FlatAppearance.BorderSize = 0;
            this.btnSettings.FlatAppearance.MouseDownBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(236)))), ((int)(((byte)(240)))));
            this.btnSettings.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(242)))), ((int)(((byte)(244)))), ((int)(((byte)(247)))));
            this.btnSettings.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnSettings.Font = new System.Drawing.Font("Segoe MDL2 Assets", 13F);
            this.btnSettings.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(102)))), ((int)(((byte)(112)))), ((int)(((byte)(133)))));
            this.btnSettings.Location = new System.Drawing.Point(438, 32);
            this.btnSettings.Name = "btnSettings";
            this.btnSettings.Size = new System.Drawing.Size(36, 36);
            this.btnSettings.TabIndex = 16;
            this.btnSettings.Text = "\uE713";
            this.btnSettings.UseVisualStyleBackColor = false;
            this.btnSettings.Click += new System.EventHandler(this.btnSettings_Click);
            // 
            // lblUsername
            // 
            this.lblUsername.AutoSize = true;
            this.lblUsername.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblUsername.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(52)))), ((int)(((byte)(64)))), ((int)(((byte)(84)))));
            this.lblUsername.Location = new System.Drawing.Point(44, 132);
            this.lblUsername.Name = "lblUsername";
            this.lblUsername.Size = new System.Drawing.Size(40, 17);
            this.lblUsername.TabIndex = 4;
            this.lblUsername.Text = "账号";
            // 
            // txtUsername
            // 
            this.txtUsername.AutoSize = false;
            this.txtUsername.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.txtUsername.Font = new System.Drawing.Font("Microsoft YaHei", 9.5F);
            this.txtUsername.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.txtUsername.Location = new System.Drawing.Point(44, 154);
            this.txtUsername.Name = "txtUsername";
            this.txtUsername.Size = new System.Drawing.Size(422, 30);
            this.txtUsername.TabIndex = 5;
            // 
            // lblPassword
            // 
            this.lblPassword.AutoSize = true;
            this.lblPassword.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblPassword.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(52)))), ((int)(((byte)(64)))), ((int)(((byte)(84)))));
            this.lblPassword.Location = new System.Drawing.Point(44, 202);
            this.lblPassword.Name = "lblPassword";
            this.lblPassword.Size = new System.Drawing.Size(40, 17);
            this.lblPassword.TabIndex = 6;
            this.lblPassword.Text = "密码";
            // 
            // txtPassword
            // 
            this.txtPassword.AutoSize = false;
            this.txtPassword.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.txtPassword.Font = new System.Drawing.Font("Microsoft YaHei", 9.5F);
            this.txtPassword.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.txtPassword.Location = new System.Drawing.Point(44, 224);
            this.txtPassword.Name = "txtPassword";
            this.txtPassword.Size = new System.Drawing.Size(422, 30);
            this.txtPassword.TabIndex = 7;
            this.txtPassword.UseSystemPasswordChar = true;
            // 
            // lblCaptcha
            // 
            this.lblCaptcha.AutoSize = true;
            this.lblCaptcha.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblCaptcha.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(52)))), ((int)(((byte)(64)))), ((int)(((byte)(84)))));
            this.lblCaptcha.Location = new System.Drawing.Point(44, 272);
            this.lblCaptcha.Name = "lblCaptcha";
            this.lblCaptcha.Size = new System.Drawing.Size(56, 17);
            this.lblCaptcha.TabIndex = 8;
            this.lblCaptcha.Text = "验证码";
            // 
            // txtCaptcha
            // 
            this.txtCaptcha.AutoSize = false;
            this.txtCaptcha.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.txtCaptcha.Font = new System.Drawing.Font("Microsoft YaHei", 9.5F);
            this.txtCaptcha.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(16)))), ((int)(((byte)(24)))), ((int)(((byte)(40)))));
            this.txtCaptcha.Location = new System.Drawing.Point(44, 294);
            this.txtCaptcha.Name = "txtCaptcha";
            this.txtCaptcha.Size = new System.Drawing.Size(200, 30);
            this.txtCaptcha.TabIndex = 9;
            // 
            // picCaptcha
            // 
            this.picCaptcha.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            this.picCaptcha.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.picCaptcha.Cursor = System.Windows.Forms.Cursors.Hand;
            this.picCaptcha.Location = new System.Drawing.Point(262, 290);
            this.picCaptcha.Name = "picCaptcha";
            this.picCaptcha.Size = new System.Drawing.Size(204, 38);
            this.picCaptcha.SizeMode = System.Windows.Forms.PictureBoxSizeMode.Zoom;
            this.picCaptcha.TabIndex = 10;
            this.picCaptcha.TabStop = false;
            this.picCaptcha.Click += new System.EventHandler(this.picCaptcha_Click);
            // 
            // lnkRefreshCaptcha
            // 
            this.lnkRefreshCaptcha.AutoSize = true;
            this.lnkRefreshCaptcha.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.lnkRefreshCaptcha.LinkBehavior = System.Windows.Forms.LinkBehavior.HoverUnderline;
            this.lnkRefreshCaptcha.LinkColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.lnkRefreshCaptcha.Location = new System.Drawing.Point(262, 334);
            this.lnkRefreshCaptcha.Name = "lnkRefreshCaptcha";
            this.lnkRefreshCaptcha.Size = new System.Drawing.Size(103, 17);
            this.lnkRefreshCaptcha.TabIndex = 12;
            this.lnkRefreshCaptcha.TabStop = true;
            this.lnkRefreshCaptcha.Text = "看不清？换一张";
            this.lnkRefreshCaptcha.LinkClicked += new System.Windows.Forms.LinkLabelLinkClickedEventHandler(this.lnkRefreshCaptcha_LinkClicked);
            // 
            // chkAutoLogin
            // 
            this.chkAutoLogin.AutoSize = true;
            this.chkAutoLogin.Cursor = System.Windows.Forms.Cursors.Hand;
            this.chkAutoLogin.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.chkAutoLogin.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.chkAutoLogin.Location = new System.Drawing.Point(42, 332);
            this.chkAutoLogin.Name = "chkAutoLogin";
            this.chkAutoLogin.Size = new System.Drawing.Size(198, 21);
            this.chkAutoLogin.TabIndex = 11;
            this.chkAutoLogin.Text = "自动登录（下次启动免输入）";
            this.chkAutoLogin.UseVisualStyleBackColor = true;
            // 
            // btnLogin
            // 
            this.btnLogin.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.btnLogin.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnLogin.FlatAppearance.BorderSize = 0;
            this.btnLogin.FlatAppearance.MouseDownBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(30)))), ((int)(((byte)(64)))), ((int)(((byte)(175)))));
            this.btnLogin.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(29)))), ((int)(((byte)(78)))), ((int)(((byte)(216)))));
            this.btnLogin.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnLogin.Font = new System.Drawing.Font("Microsoft YaHei", 10F, System.Drawing.FontStyle.Bold);
            this.btnLogin.ForeColor = System.Drawing.Color.White;
            this.btnLogin.Location = new System.Drawing.Point(44, 372);
            this.btnLogin.Name = "btnLogin";
            this.btnLogin.Size = new System.Drawing.Size(200, 42);
            this.btnLogin.TabIndex = 13;
            this.btnLogin.Text = "登录";
            this.btnLogin.UseVisualStyleBackColor = false;
            this.btnLogin.Click += new System.EventHandler(this.btnLogin_Click);
            // 
            // btnCancel
            // 
            this.btnCancel.BackColor = System.Drawing.Color.White;
            this.btnCancel.Cursor = System.Windows.Forms.Cursors.Hand;
            this.btnCancel.FlatAppearance.BorderColor = System.Drawing.Color.FromArgb(((int)(((byte)(208)))), ((int)(((byte)(213)))), ((int)(((byte)(221)))));
            this.btnCancel.FlatAppearance.MouseOverBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(249)))), ((int)(((byte)(250)))), ((int)(((byte)(251)))));
            this.btnCancel.FlatStyle = System.Windows.Forms.FlatStyle.Flat;
            this.btnCancel.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.btnCancel.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(71)))), ((int)(((byte)(84)))), ((int)(((byte)(103)))));
            this.btnCancel.Location = new System.Drawing.Point(256, 372);
            this.btnCancel.Name = "btnCancel";
            this.btnCancel.Size = new System.Drawing.Size(100, 42);
            this.btnCancel.TabIndex = 14;
            this.btnCancel.Text = "退出";
            this.btnCancel.UseVisualStyleBackColor = false;
            this.btnCancel.Click += new System.EventHandler(this.btnCancel_Click);
            // 
            // lblStatus
            // 
            this.lblStatus.AutoEllipsis = true;
            this.lblStatus.AutoSize = false;
            this.lblStatus.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.lblStatus.ForeColor = System.Drawing.Color.DimGray;
            this.lblStatus.Location = new System.Drawing.Point(44, 424);
            this.lblStatus.Name = "lblStatus";
            this.lblStatus.Size = new System.Drawing.Size(422, 24);
            this.lblStatus.TabIndex = 15;
            this.lblStatus.TextAlign = System.Drawing.ContentAlignment.MiddleLeft;
            // 
            // panelBrand
            // 
            this.panelBrand.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(15)))), ((int)(((byte)(23)))), ((int)(((byte)(42)))));
            this.panelBrand.Controls.Add(this.lblBrandTitle);
            this.panelBrand.Controls.Add(this.lblBrandSubtitle);
            this.panelBrand.Controls.Add(this.lblBrandDesc);
            this.panelBrand.Controls.Add(this.lblBrandFeatures);
            this.panelBrand.Controls.Add(this.lblBrandVersion);
            this.panelBrand.Controls.Add(this.panelBrandAccent);
            this.panelBrand.Dock = System.Windows.Forms.DockStyle.Left;
            this.panelBrand.Location = new System.Drawing.Point(0, 0);
            this.panelBrand.Name = "panelBrand";
            this.panelBrand.Size = new System.Drawing.Size(330, 520);
            this.panelBrand.TabIndex = 0;
            // 
            // lblBrandTitle
            // 
            this.lblBrandTitle.AutoSize = true;
            this.lblBrandTitle.Font = new System.Drawing.Font("Microsoft YaHei", 28F, System.Drawing.FontStyle.Bold);
            this.lblBrandTitle.ForeColor = System.Drawing.Color.White;
            this.lblBrandTitle.Location = new System.Drawing.Point(34, 48);
            this.lblBrandTitle.Name = "lblBrandTitle";
            this.lblBrandTitle.Size = new System.Drawing.Size(132, 50);
            this.lblBrandTitle.TabIndex = 0;
            this.lblBrandTitle.Text = "易云印";
            // 
            // lblBrandSubtitle
            // 
            this.lblBrandSubtitle.AutoSize = true;
            this.lblBrandSubtitle.Font = new System.Drawing.Font("Microsoft YaHei", 14F);
            this.lblBrandSubtitle.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(147)))), ((int)(((byte)(197)))), ((int)(((byte)(253)))));
            this.lblBrandSubtitle.Location = new System.Drawing.Point(36, 105);
            this.lblBrandSubtitle.Name = "lblBrandSubtitle";
            this.lblBrandSubtitle.Size = new System.Drawing.Size(126, 25);
            this.lblBrandSubtitle.TabIndex = 1;
            this.lblBrandSubtitle.Text = "云打印客户端";
            // 
            // lblBrandDesc
            // 
            this.lblBrandDesc.AutoSize = true;
            this.lblBrandDesc.Font = new System.Drawing.Font("Microsoft YaHei", 10F);
            this.lblBrandDesc.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(148)))), ((int)(((byte)(163)))), ((int)(((byte)(184)))));
            this.lblBrandDesc.Location = new System.Drawing.Point(36, 166);
            this.lblBrandDesc.Name = "lblBrandDesc";
            this.lblBrandDesc.Size = new System.Drawing.Size(212, 38);
            this.lblBrandDesc.TabIndex = 2;
            this.lblBrandDesc.Text = "让打印任务自动抵达门店\r\n连接后台与本地打印机";
            // 
            // lblBrandFeatures
            // 
            this.lblBrandFeatures.AutoSize = true;
            this.lblBrandFeatures.Font = new System.Drawing.Font("Microsoft YaHei", 9.5F);
            this.lblBrandFeatures.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(203)))), ((int)(((byte)(213)))), ((int)(((byte)(225)))));
            this.lblBrandFeatures.Location = new System.Drawing.Point(36, 240);
            this.lblBrandFeatures.Name = "lblBrandFeatures";
            this.lblBrandFeatures.Size = new System.Drawing.Size(248, 68);
            this.lblBrandFeatures.TabIndex = 3;
            this.lblBrandFeatures.Text = "· 自动登记设备并上报打印机能力\r\n· 支持 A4 / A3、彩色、双面参数\r\n· 每 5 秒拉取任务并静默打印 PDF";
            // 
            // lblBrandVersion
            // 
            this.lblBrandVersion.AutoSize = true;
            this.lblBrandVersion.Font = new System.Drawing.Font("Microsoft YaHei", 8.5F);
            this.lblBrandVersion.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(100)))), ((int)(((byte)(116)))), ((int)(((byte)(139)))));
            this.lblBrandVersion.Location = new System.Drawing.Point(36, 482);
            this.lblBrandVersion.Name = "lblBrandVersion";
            this.lblBrandVersion.Size = new System.Drawing.Size(140, 17);
            this.lblBrandVersion.TabIndex = 4;
            this.lblBrandVersion.Text = "本地打印客户端 · Windows";
            // 
            // panelBrandAccent
            // 
            this.panelBrandAccent.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(37)))), ((int)(((byte)(99)))), ((int)(((byte)(235)))));
            this.panelBrandAccent.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelBrandAccent.Location = new System.Drawing.Point(0, 0);
            this.panelBrandAccent.Name = "panelBrandAccent";
            this.panelBrandAccent.Size = new System.Drawing.Size(330, 4);
            this.panelBrandAccent.TabIndex = 5;
            // 
            // LoginForm
            // 
            this.AcceptButton = this.btnLogin;
            this.AutoScaleDimensions = new System.Drawing.SizeF(8F, 17F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(244)))), ((int)(((byte)(246)))), ((int)(((byte)(250)))));
            this.CancelButton = this.btnCancel;
            this.ClientSize = new System.Drawing.Size(840, 520);
            this.Controls.Add(this.panelForm);
            this.Controls.Add(this.panelBrand);
            this.Font = new System.Drawing.Font("Microsoft YaHei", 9F);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedDialog;
            this.MaximizeBox = false;
            this.MinimizeBox = false;
            this.Name = "LoginForm";
            this.StartPosition = System.Windows.Forms.FormStartPosition.CenterScreen;
            this.Text = "易云印 云打印客户端 - 登录";
            this.Load += new System.EventHandler(this.LoginForm_Load);
            this.panelForm.ResumeLayout(false);
            this.panelForm.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.picCaptcha)).EndInit();
            this.panelBrand.ResumeLayout(false);
            this.panelBrand.PerformLayout();
            this.ResumeLayout(false);
        }

        #endregion

        private System.Windows.Forms.Panel panelForm;
        private System.Windows.Forms.Label lblLoginTitle;
        private System.Windows.Forms.Label lblLoginSubtitle;
        private System.Windows.Forms.Label lblUsername;
        private System.Windows.Forms.TextBox txtUsername;
        private System.Windows.Forms.Label lblPassword;
        private System.Windows.Forms.TextBox txtPassword;
        private System.Windows.Forms.Label lblCaptcha;
        private System.Windows.Forms.TextBox txtCaptcha;
        private System.Windows.Forms.PictureBox picCaptcha;
        private System.Windows.Forms.LinkLabel lnkRefreshCaptcha;
        private System.Windows.Forms.CheckBox chkAutoLogin;
        private System.Windows.Forms.Button btnSettings;
        private System.Windows.Forms.Button btnLogin;
        private System.Windows.Forms.Button btnCancel;
        private System.Windows.Forms.Label lblStatus;
        private System.Windows.Forms.Panel panelBrand;
        private System.Windows.Forms.Label lblBrandTitle;
        private System.Windows.Forms.Label lblBrandSubtitle;
        private System.Windows.Forms.Label lblBrandDesc;
        private System.Windows.Forms.Label lblBrandFeatures;
        private System.Windows.Forms.Label lblBrandVersion;
        private System.Windows.Forms.Panel panelBrandAccent;
    }
}
