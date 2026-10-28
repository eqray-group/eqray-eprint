-- =====================================================================
-- 云打印平台 · 表结构（范围内：到店打印 + 取件码核销，不含快递发货）
-- ---------------------------------------------------------------------
-- 全局约定：
--   1. 主键 id 为 bigint（有符号），由应用侧雪花 ID 生成，不使用 AUTO_INCREMENT；
--   2. 引擎 InnoDB，字符集 utf8mb4，排序规则 utf8mb4_unicode_ci；
--   3. 金额统一 bigint，单位「分」；页数（page）与印张（sheet = 页数 × 份数）分开存；
--   4. 启停类状态 tinyint：1 启用 / 0 禁用；业务状态值见各字段注释；
--   5. 审计字段统一 creator / created_at / updater / updated_at / deleted_at；
--   6. 唯一键直接建在业务键上，软删除后不复用（与 system_users 的做法一致）；
--      门店成员绑定（print_shop_staff）与价格规则属于可重建数据，不建 deleted_at，
--      删除即物理删除，避免「软删占着唯一键」导致的重新绑定冲突；
--   7. 鉴权口径：打印客户端不设独立设备令牌，由店主/店员用 system_users 账号登录，
--      走后台令牌 + 权限点；shopId 一律取自 print_shop_staff 绑定，不接受前端传入；
--   8. 打印执行口径：客户端只调用 Windows 打印子系统（spooler）静默打印，
--      打印机通过 system_name（Windows 队列名）指定，只有 Windows 能识别的打印机才建记录；
--   9. 支付与钱包：金额只信服务端，余额变动必须与流水同事务；
--      支付单用 biz_type + biz_id 泛化，不反向依赖打印模块。
-- ---------------------------------------------------------------------
-- 表清单（16 张）：
--   打印业务 9 张：
--   print_shop            门店
--   print_shop_staff      门店成员（后台账号 ↔ 门店）
--   print_device          打印设备（Windows 客户端实例）
--   print_printer         打印机（装在设备上的 Windows 打印队列）
--   print_price_config    打印价格配置
--   print_file            打印文件（infra_file 的打印元数据扩展）
--   print_order           打印订单（到店：支付成功即打印，取件码核销完成）
--   print_order_item      订单明细，兼打印任务
--   print_banner          门店轮播图
--   支付与钱包 6 张：
--   pay_channel_config    支付渠道配置（微信小程序 / 微信 H5 / 支付宝）
--   pay_wallet            会员钱包账户（只有会员有余额）
--   pay_wallet_transaction 钱包流水
--   pay_payment           支付单（订单支付 / 钱包充值共用）
--   pay_refund            退款单
--   pay_notify_log        支付渠道回调日志
--   会员 1 张：
--   member_card           会员卡（等级权益卡，等级折扣的载体）
-- =====================================================================


-- ---------------------------------------------------------------------
-- 门店
-- ---------------------------------------------------------------------
CREATE TABLE `print_shop`
(
    `id`             bigint         NOT NULL COMMENT '主键',
    `name`           varchar(64)    NOT NULL COMMENT '门店名称',
    `contact_name`   varchar(32)    NOT NULL DEFAULT '' COMMENT '门店联系人',
    `contact_mobile` varchar(20)    NOT NULL DEFAULT '' COMMENT '门店联系电话',
    `province`       varchar(32)    NOT NULL DEFAULT '' COMMENT '省',
    `city`           varchar(32)    NOT NULL DEFAULT '' COMMENT '市',
    `district`       varchar(32)    NOT NULL DEFAULT '' COMMENT '区/县',
    `address`        varchar(255)   NOT NULL DEFAULT '' COMMENT '详细地址',
    `longitude`      decimal(10, 7) NOT NULL DEFAULT '0.0000000' COMMENT '经度',
    `latitude`       decimal(10, 7) NOT NULL DEFAULT '0.0000000' COMMENT '纬度',
    `avatar`         varchar(512)   NOT NULL DEFAULT '' COMMENT '门头照地址',
    `business_hours` json           NULL COMMENT '营业时间 JSON，如 [{"start":"09:00","end":"21:00"}]；NULL 表示全天营业',
    `status`         tinyint        NOT NULL DEFAULT '1' COMMENT '营业状态：0 停业 / 1 营业',
    `audit_status`   tinyint        NOT NULL DEFAULT '0' COMMENT '审核状态：0 待审核 / 1 审核通过 / 2 审核拒绝',
    `audit_remark`   varchar(255)   NOT NULL DEFAULT '' COMMENT '审核意见',
    `remark`         varchar(255)   NOT NULL DEFAULT '' COMMENT '备注',
    `creator`        varchar(64)    NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`     datetime       NULL COMMENT '创建时间',
    `updater`        varchar(64)    NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`     datetime       NULL COMMENT '更新时间',
    `deleted_at`     datetime       NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_audit_status` (`audit_status`),
    KEY `idx_geo` (`latitude`, `longitude`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='打印门店表';


-- ---------------------------------------------------------------------
-- 门店成员：后台账号与门店的多对多绑定，角色按 (user_id, shop_id) 存
-- 解绑走物理删除（同 member_social_binding），避免唯一键挡住重新绑定
-- ---------------------------------------------------------------------
CREATE TABLE `print_shop_staff`
(
    `id`         bigint      NOT NULL COMMENT '主键',
    `user_id`    bigint      NOT NULL COMMENT '后台账号 ID，对应 system_users.id',
    `shop_id`    bigint      NOT NULL COMMENT '门店 ID，对应 print_shop.id',
    `role`       varchar(16) NOT NULL DEFAULT 'staff' COMMENT '门店角色：owner 店主 / staff 店员',
    `status`     tinyint     NOT NULL DEFAULT '1' COMMENT '状态：0 停用 / 1 启用',
    `creator`    varchar(64) NOT NULL DEFAULT '' COMMENT '创建人（操作人）',
    `created_at` datetime    NULL COMMENT '创建时间',
    `updater`    varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at` datetime    NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_shop` (`user_id`, `shop_id`),
    KEY `idx_shop` (`shop_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='门店成员绑定表';


-- ---------------------------------------------------------------------
-- 打印设备：安装在门店的 Windows 客户端实例
-- 客户端由店主/店员用 system_users 账号登录（后台令牌 + 权限点），无独立设备令牌；
-- 设备归属门店 = 该账号在 print_shop_staff 中的绑定门店，首次登录时登记设备
-- 在线状态由心跳推导，online_status 为冗余字段便于统计在线率
-- ---------------------------------------------------------------------
CREATE TABLE `print_device`
(
    `id`                 bigint       NOT NULL COMMENT '主键',
    `print_shop_id`      bigint       NOT NULL COMMENT '门店 ID',
    `code`               varchar(64)  NOT NULL COMMENT '设备唯一编码（客户端机器指纹），全局唯一',
    `name`               varchar(128) NOT NULL COMMENT '设备名称',
    `client_version`     varchar(32)  NOT NULL DEFAULT '' COMMENT 'Windows 客户端版本号，用于强制升级校验',
    `status`             tinyint      NOT NULL DEFAULT '1' COMMENT '状态：0 禁用 / 1 启用',
    `online_status`      tinyint      NOT NULL DEFAULT '0' COMMENT '在线状态：0 离线 / 1 在线（心跳超过 90 秒由服务端置 0）',
    `ip`                 varchar(50)  NOT NULL DEFAULT '' COMMENT '最近上报 IP',
    `last_heartbeat_at`  datetime     NULL COMMENT '最后心跳时间',
    `last_login_user_id` bigint       NOT NULL DEFAULT '0' COMMENT '最近登录该设备的后台账号，对应 system_users.id',
    `last_error_info`    varchar(512) NOT NULL DEFAULT '' COMMENT '最后错误信息',
    `remark`             varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `creator`            varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`         datetime     NULL COMMENT '创建时间',
    `updater`            varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`         datetime     NULL COMMENT '更新时间',
    `deleted_at`         datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_shop` (`print_shop_id`),
    KEY `idx_online_status` (`online_status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='打印设备表';


-- ---------------------------------------------------------------------
-- 打印机：一条记录 = 某台设备上的一个 Windows 打印队列（队列装在具体机器上，
-- 同型号打印机在两台电脑上安装会产生两条记录，各自队列名可能相同）。
-- 只收录 Windows 能识别的打印机，客户端按 system_name 指定队列静默打印，
-- 不直连 IPP / 厂商 SDK；不设默认打印机，路由只按本表配置的参数能力匹配，
-- 无匹配打印机时下单直接报错，不兜底、不降级打印；
-- 删除打印机用 status=0 停用，避免队列名被删除记录占用。
-- ---------------------------------------------------------------------
CREATE TABLE `print_printer`
(
    `id`             bigint       NOT NULL COMMENT '主键',
    `print_shop_id`  bigint       NOT NULL COMMENT '门店 ID（冗余，便于按店查询）',
    `device_id`      bigint       NOT NULL COMMENT '所属设备 ID，对应 print_device.id；打印队列装在这台 Windows 机器上',
    `name`           varchar(64)  NOT NULL COMMENT '打印机别名，如 A4 黑白 1 号机（后台展示用）',
    `system_name`    varchar(128) NOT NULL COMMENT 'Windows 打印队列名称（spooler 中的 PrinterName），同一设备内唯一，静默打印按此指定',
    `model`          varchar(64)  NOT NULL DEFAULT '' COMMENT '型号',
    `paper_types`    varchar(64)  NOT NULL DEFAULT 'A4' COMMENT '本机支持的纸张，逗号分隔：A4,A3，由门店后台勾选配置',
    `color_support`  tinyint      NOT NULL DEFAULT '0' COMMENT '是否支持彩色：0 否 / 1 是，由门店后台配置，彩色任务遇 0 必须失败',
    `duplex_support` tinyint      NOT NULL DEFAULT '0' COMMENT '是否支持双面：0 否 / 1 是，由门店后台配置',
    `status`         tinyint      NOT NULL DEFAULT '1' COMMENT '状态：0 停用 / 1 启用',
    `remark`         varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `creator`        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`     datetime     NULL COMMENT '创建时间',
    `updater`        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`     datetime     NULL COMMENT '更新时间',
    `deleted_at`     datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_device_system_name` (`device_id`, `system_name`),
    KEY `idx_shop` (`print_shop_id`),
    KEY `idx_device` (`device_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='设备打印机表';


-- ---------------------------------------------------------------------
-- 打印价格配置：门店 × 纸张 × 色彩 × 单双面 唯一确定一个单价（分/页）
-- 价格规则可重建，删除走物理删除
-- ---------------------------------------------------------------------
CREATE TABLE `print_price_config`
(
    `id`            bigint      NOT NULL COMMENT '主键',
    `print_shop_id` bigint      NOT NULL COMMENT '门店 ID',
    `paper_type`    varchar(16) NOT NULL COMMENT '纸张类型：A4 / A3',
    `color_type`    varchar(16) NOT NULL COMMENT '色彩：bw 黑白 / color 彩色',
    `print_mode`    varchar(16) NOT NULL COMMENT '单双面：single 单面 / double 双面',
    `price`         bigint      NOT NULL COMMENT '单价（分/页）',
    `status`        tinyint     NOT NULL DEFAULT '1' COMMENT '状态：0 停用 / 1 启用',
    `creator`       varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`    datetime    NULL COMMENT '创建时间',
    `updater`       varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`    datetime    NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_rule` (`print_shop_id`, `paper_type`, `color_type`, `print_mode`),
    KEY `idx_shop` (`print_shop_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='打印价格配置表';


-- ---------------------------------------------------------------------
-- 打印文件：不重复建存储，源文件与可打印 PDF 都落在 infra_file
-- DOCX 由服务端转 PDF 后写 pdf_file_id；打印参数只在订单明细上，本表不存
-- ---------------------------------------------------------------------
CREATE TABLE `print_file`
(
    `id`             bigint       NOT NULL COMMENT '主键',
    `member_user_id` bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id',
    `infra_file_id`  bigint       NOT NULL COMMENT '源文件，对应 infra_file.id',
    `pdf_file_id`    bigint       NOT NULL DEFAULT '0' COMMENT '可打印 PDF，对应 infra_file.id；源文件本身即 PDF 时与 infra_file_id 相同',
    `file_name`      varchar(255) NOT NULL COMMENT '文件名快照（含扩展名）',
    `file_ext`       varchar(16)  NOT NULL COMMENT '扩展名：pdf / jpg / png / docx',
    `file_size`      bigint       NOT NULL DEFAULT '0' COMMENT '文件大小（字节）',
    `file_sha256`    char(64)     NOT NULL DEFAULT '' COMMENT '可打印文件 SHA-256，客户端下载后校验',
    `page_count`     int          NOT NULL DEFAULT '0' COMMENT '可打印 PDF 总页数，由服务端计算',
    `convert_status` tinyint      NOT NULL DEFAULT '0' COMMENT '转换状态：0 待处理 / 1 转换中 / 2 可打印 / 3 转换失败',
    `convert_error`  varchar(255) NOT NULL DEFAULT '' COMMENT '转换失败原因',
    `expire_time`    datetime     NULL COMMENT '文件清理时间，到期由定时任务删除存储与记录',
    `creator`        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`     datetime     NULL COMMENT '创建时间',
    `updater`        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`     datetime     NULL COMMENT '更新时间',
    `deleted_at`     datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_infra_file` (`infra_file_id`),
    KEY `idx_member` (`member_user_id`),
    KEY `idx_expire` (`expire_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='打印文件表';


-- ---------------------------------------------------------------------
-- 打印订单（到店）：支付成功即自动打印，打印完成转待取件，取件码核销完成
-- 订单状态单线推进：0 → 1 → 2 → 3 → 4，旁支 5/6/7/8
-- pay_status 只表达钱，order_status 表达履约，两者组合由服务层约束
-- ---------------------------------------------------------------------
CREATE TABLE `print_order`
(
    `id`                   bigint       NOT NULL COMMENT '主键',
    `order_no`             varchar(32)  NOT NULL COMMENT '订单号，全局唯一',
    `print_shop_id`        bigint       NOT NULL COMMENT '门店 ID',
    `member_user_id`       bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id',
    `device_id`            bigint       NOT NULL DEFAULT '0' COMMENT '执行设备 ID，0 表示尚未路由（支付成功后由服务端分配）',
    `pickup_code`          char(6)      NOT NULL DEFAULT '' COMMENT '取件码，6 位数字，门店内未完成订单唯一',
    `pickup_expire_time`   datetime     NULL COMMENT '取件保留到期时间，默认支付后 7 天，到期转 8 超期关闭',
    `pickup_verify_at`     datetime     NULL COMMENT '取件核销时间',
    `pickup_verify_user_id` bigint      NOT NULL DEFAULT '0' COMMENT '核销操作人，对应 system_users.id',
    `item_count`           int          NOT NULL DEFAULT '0' COMMENT '明细数量（文件数）',
    `total_pages`          int          NOT NULL DEFAULT '0' COMMENT '总页数（不含份数）',
    `total_sheets`         int          NOT NULL DEFAULT '0' COMMENT '总印张数（页数 × 份数）',
    `total_amount`         bigint       NOT NULL DEFAULT '0' COMMENT '明细金额合计（分）',
    `discount_percent`     int          NOT NULL DEFAULT '100' COMMENT '下单时的会员折扣快照：100 无折扣 / 95 表示 95 折；等级变更不影响历史订单',
    `discount_amount`      bigint       NOT NULL DEFAULT '0' COMMENT '优惠金额（分）= 明细合计 × (100 - discount_percent) / 100，向下取整',
    `pay_amount`           bigint       NOT NULL DEFAULT '0' COMMENT '应付金额（分）= 明细合计 - 优惠',
    `refund_amount`        bigint       NOT NULL DEFAULT '0' COMMENT '已退款金额（分）',
    `pay_status`           tinyint      NOT NULL DEFAULT '0' COMMENT '支付状态：0 未支付 / 1 已支付 / 2 退款中 / 3 已退款',
    `pay_type`             varchar(16)  NOT NULL DEFAULT '' COMMENT '支付方式：wechat 微信 / alipay 支付宝 / wallet 余额',
    `trade_no`             varchar(64)  NOT NULL DEFAULT '' COMMENT '第三方支付交易号',
    `pay_time`             datetime     NULL COMMENT '支付时间',
    `order_status`         tinyint      NOT NULL DEFAULT '0' COMMENT '订单状态：0 待支付 / 1 待打印 / 2 打印中 / 3 待取件 / 4 已完成 / 5 已取消 / 6 退款中 / 7 已退款 / 8 超期关闭',
    `pay_expire_time`      datetime     NULL COMMENT '未支付关单时间，默认创建后 30 分钟',
    `finish_time`          datetime     NULL COMMENT '完成时间',
    `cancel_time`          datetime     NULL COMMENT '取消时间',
    `cancel_reason`        varchar(255) NOT NULL DEFAULT '' COMMENT '取消原因',
    `close_time`           datetime     NULL COMMENT '关闭时间（超时关单 / 超期未取）',
    `remark`               varchar(255) NOT NULL DEFAULT '' COMMENT '用户备注',
    `creator`              varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`           datetime     NULL COMMENT '创建时间',
    `updater`              varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`           datetime     NULL COMMENT '更新时间',
    `deleted_at`           datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_member` (`member_user_id`, `created_at`),
    KEY `idx_shop_status` (`print_shop_id`, `order_status`, `created_at`),
    KEY `idx_shop_pickup` (`print_shop_id`, `pickup_code`),
    KEY `idx_trade_no` (`trade_no`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='打印订单表';


-- ---------------------------------------------------------------------
-- 订单明细，同时承担「打印任务」：一条明细 = 一个文件 + 一套参数 + 一台打印机
-- task_no 全局唯一，客户端据此去重防重复打印
-- 打印参数以本表为唯一来源，print_config 只是全量快照，客户端与门店后台均不得覆盖
-- ---------------------------------------------------------------------
CREATE TABLE `print_order_item`
(
    `id`              bigint       NOT NULL COMMENT '主键',
    `print_shop_id`   bigint       NOT NULL COMMENT '门店 ID（冗余，便于按店统计）',
    `order_id`        bigint       NOT NULL COMMENT '订单 ID，对应 print_order.id',
    `task_no`         varchar(64)  NOT NULL COMMENT '任务号，全局唯一，客户端防重复打印依据',
    `file_id`         bigint       NOT NULL COMMENT '打印文件 ID，对应 print_file.id',
    `file_name`       varchar(255) NOT NULL COMMENT '文件名快照',
    `file_page_count` int          NOT NULL DEFAULT '0' COMMENT '文件总页数快照',
    `paper_type`      varchar(16)  NOT NULL DEFAULT 'A4' COMMENT '纸张类型：A4 / A3',
    `color_type`      varchar(16)  NOT NULL DEFAULT 'bw' COMMENT '色彩：bw 黑白 / color 彩色',
    `print_mode`      varchar(16)  NOT NULL DEFAULT 'single' COMMENT '单双面：single 单面 / double 双面',
    `orientation`     varchar(16)  NOT NULL DEFAULT 'portrait' COMMENT '纸张方向：portrait 纵向 / landscape 横向',
    `page_range`      varchar(128) NOT NULL DEFAULT '' COMMENT '页码范围，如 1-5,8,10-12；空表示全部',
    `page_count`      int          NOT NULL DEFAULT '0' COMMENT '本明细计费页数（按页码范围计算）',
    `copies`          int          NOT NULL DEFAULT '1' COMMENT '打印份数',
    `sheet_count`     int          NOT NULL DEFAULT '0' COMMENT '印张数 = page_count × copies',
    `print_config`    json         NULL COMMENT '打印参数全量快照，服务端下发、客户端只读',
    `amount`          bigint       NOT NULL DEFAULT '0' COMMENT '小计金额（分）',
    `status`          tinyint      NOT NULL DEFAULT '0' COMMENT '任务状态：0 待分配 / 1 已分配 / 2 打印中 / 3 已完成 / 4 失败 / 5 已取消',
    `device_id`       bigint       NOT NULL DEFAULT '0' COMMENT '执行设备 ID，0 表示未分配',
    `printer_id`      bigint       NOT NULL DEFAULT '0' COMMENT '执行打印机 ID，0 表示未分配；创建后锁定，换机需服务端重建任务',
    `retry_count`     int          NOT NULL DEFAULT '0' COMMENT '重试次数，上限 3，超限转人工并触发订单退款',
    `error_code`      varchar(32)  NOT NULL DEFAULT '' COMMENT '失败错误码，如 PAPER_JAM / COLOR_UNSUPPORTED / FILE_BROKEN',
    `fail_reason`     varchar(255) NOT NULL DEFAULT '' COMMENT '失败原因',
    `printed_pages`   int          NOT NULL DEFAULT '0' COMMENT '实际打印页数（客户端回执）',
    `started_at`      datetime     NULL COMMENT '开始打印时间',
    `finished_at`     datetime     NULL COMMENT '打印结束时间',
    `creator`         varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`      datetime     NULL COMMENT '创建时间',
    `updater`         varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`      datetime     NULL COMMENT '更新时间',
    `deleted_at`      datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_task_no` (`task_no`),
    KEY `idx_order` (`order_id`),
    KEY `idx_device_status` (`device_id`, `status`),
    KEY `idx_shop_status` (`print_shop_id`, `status`),
    KEY `idx_file` (`file_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='打印订单明细表';


-- ---------------------------------------------------------------------
-- 门店轮播图：print_shop_id = 0 表示全平台通用
-- ---------------------------------------------------------------------
CREATE TABLE `print_banner`
(
    `id`            bigint       NOT NULL COMMENT '主键',
    `print_shop_id` bigint       NOT NULL DEFAULT '0' COMMENT '门店 ID，0 表示全平台通用',
    `title`         varchar(128) NOT NULL COMMENT '标题',
    `image_url`     varchar(512) NOT NULL COMMENT '图片地址',
    `link_type`     varchar(16)  NOT NULL DEFAULT 'none' COMMENT '跳转类型：none 不跳 / url 外部链接 / page 站内页面',
    `link_url`      varchar(512) NOT NULL DEFAULT '' COMMENT '跳转地址',
    `sort`          int          NOT NULL DEFAULT '0' COMMENT '排序，越小越靠前',
    `status`        tinyint      NOT NULL DEFAULT '1' COMMENT '状态：0 停用 / 1 启用',
    `creator`       varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`    datetime     NULL COMMENT '创建时间',
    `updater`       varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`    datetime     NULL COMMENT '更新时间',
    `deleted_at`    datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_shop_sort` (`print_shop_id`, `sort`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='门店轮播图表';


-- =====================================================================
-- 支付与钱包
-- =====================================================================


-- ---------------------------------------------------------------------
-- 支付渠道配置：后台维护，支付与退款时按 (channel, scene) 取启用配置。
-- 平台统一收款：所有门店订单进同一个商户号，不做门店独立商户号与分账，
-- 因此本表不挂 print_shop_id，唯一键固定为 (channel, scene)。
-- 微信支持三种场景，各配一条记录（同一商户号，appid 不同）：
--   wechat_mini 小程序 JSAPI（openid 属于小程序 appid）
--   wechat_mp   公众号 JSAPI（微信内置浏览器，openid 属于公众号 appid）
--   wechat_h5   手机浏览器 H5（MWEB，仅微信外浏览器，不需要 openid，需要 client_ip）
-- 前端按 UA 分流：小程序内 → mini；微信内置浏览器 → mp；其余手机浏览器 → h5。
-- openid 按 appid 隔离，同一会员的小程序与公众号 openid 不同，靠 unionid 关联；
-- 两套 openid 分别存在 member_social_user / member_social_binding（type 区分渠道）。
-- 公众号支付的 appsecret 复用 member_social_client(wechat_mp)，本表不重复存。
-- 支付宝只是渠道大类不同（alipay + wap/qr），后续接入只加配置、不改表结构。
-- 密钥类字段（private_key / api_key / platform_public_key）明文落库，接口一律不回显，
-- 编辑时留空表示不修改；因此本表必须按机密级别管控：
-- 数据库账号只授予业务库最小权限、备份加密存放、禁止把本表数据带入测试环境。
-- ---------------------------------------------------------------------
CREATE TABLE `pay_channel_config`
(
    `id`                  bigint       NOT NULL COMMENT '主键',
    `code`                varchar(32)  NOT NULL COMMENT '配置编码，全局唯一：wechat_mini / wechat_h5 / alipay_wap',
    `name`                varchar(64)  NOT NULL COMMENT '配置名称，如「微信小程序支付」',
    `channel`             varchar(16)  NOT NULL COMMENT '渠道大类：wechat 微信 / alipay 支付宝',
    `scene`               varchar(16)  NOT NULL COMMENT '支付场景：mini 小程序 / jsapi 公众号（微信内置浏览器） / h5 手机浏览器 / wap 支付宝手机网站 / qr 扫码',
    `app_id`              varchar(64)  NOT NULL COMMENT '应用 ID：微信 appid / 支付宝 appid',
    `mch_id`              varchar(64)  NOT NULL DEFAULT '' COMMENT '商户号：微信 mchid / 支付宝 PID',
    `serial_no`           varchar(64)  NOT NULL DEFAULT '' COMMENT '微信 APIv3 证书序列号',
    `private_key`         text         NULL COMMENT '商户私钥（明文，接口不回显）',
    `api_key`             varchar(256) NOT NULL DEFAULT '' COMMENT '微信 APIv3 密钥（明文，接口不回显）',
    `platform_public_key` text         NULL COMMENT '平台证书 / 支付宝公钥（明文，接口不回显）',
    `gateway`             varchar(128) NOT NULL DEFAULT '' COMMENT '网关地址，支付宝必填、微信可空',
    `notify_url`          varchar(255) NOT NULL DEFAULT '' COMMENT '支付回调地址，空表示走默认路由',
    `return_url`          varchar(255) NOT NULL DEFAULT '' COMMENT 'H5 支付完成后跳回的页面地址',
    `extra`               json         NULL COMMENT '其他渠道参数 JSON（如 H5 授权域名、小程序原始 ID）',
    `status`              tinyint      NOT NULL DEFAULT '1' COMMENT '状态：0 停用 / 1 启用',
    `sort`                int          NOT NULL DEFAULT '0' COMMENT '排序，越小越靠前',
    `remark`              varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `creator`             varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`          datetime     NULL COMMENT '创建时间',
    `updater`             varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`          datetime     NULL COMMENT '更新时间',
    `deleted_at`          datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    UNIQUE KEY `uk_channel_scene` (`channel`, `scene`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='支付渠道配置表';


-- ---------------------------------------------------------------------
-- 钱包账户：只有会员有余额账户，后台账号不参与支付
-- 余额只允许通过「钱包流水」同事务变更，禁止直接 UPDATE balance 且不写流水
-- balance 必须配合 WHERE balance >= amount 条件扣减，防止并发扣成负数
-- ---------------------------------------------------------------------
CREATE TABLE `pay_wallet`
(
    `id`             bigint      NOT NULL COMMENT '主键',
    `member_user_id` bigint      NOT NULL COMMENT '会员 ID，对应 member_user.id，一个会员一个钱包',
    `balance`        bigint      NOT NULL DEFAULT '0' COMMENT '可用余额（分），不允许为负',
    `total_income`   bigint      NOT NULL DEFAULT '0' COMMENT '累计收入（分）',
    `total_expense`  bigint      NOT NULL DEFAULT '0' COMMENT '累计支出（分）',
    `status`         tinyint     NOT NULL DEFAULT '1' COMMENT '状态：0 冻结（禁止支付） / 1 正常',
    `version`        int         NOT NULL DEFAULT '0' COMMENT '乐观锁版本号，每次余额变动 +1',
    `creator`        varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`     datetime    NULL COMMENT '创建时间',
    `updater`        varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`     datetime    NULL COMMENT '更新时间',
    `deleted_at`     datetime    NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_member` (`member_user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='会员钱包表';


-- ---------------------------------------------------------------------
-- 钱包流水：每一次余额变动必须有一条对应流水，与余额更新同事务
-- uk_biz 是幂等键：同一业务单重复入账/重复扣款会直接撞键，业务层据此判定「已处理」
-- ---------------------------------------------------------------------
CREATE TABLE `pay_wallet_transaction`
(
    `id`             bigint       NOT NULL COMMENT '主键',
    `wallet_id`      bigint       NOT NULL COMMENT '钱包 ID，对应 pay_wallet.id',
    `member_user_id` bigint       NOT NULL COMMENT '会员 ID（冗余，便于按人查流水）',
    `biz_type`       varchar(32)  NOT NULL COMMENT '业务类型：recharge 充值 / print_consume 打印消费 / print_refund 打印退款 / manual_adjust 人工调账',
    `biz_id`         bigint       NOT NULL COMMENT '业务单 ID：充值为 pay_payment.id，消费与退款为 print_order.id',
    `type`           tinyint      NOT NULL COMMENT '资金方向：1 收入 / 2 支出',
    `amount`         bigint       NOT NULL COMMENT '变动金额（分，正数，方向由 type 表示）',
    `balance_after`  bigint       NOT NULL COMMENT '变动后可用余额（分）',
    `title`          varchar(64)  NOT NULL DEFAULT '' COMMENT '展示文案，如「打印订单 P2026... 消费」',
    `remark`         varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `creator`        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`     datetime     NULL COMMENT '创建时间',
    `updater`        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`     datetime     NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_biz` (`biz_type`, `biz_id`, `type`),
    KEY `idx_wallet_created` (`wallet_id`, `created_at`),
    KEY `idx_member` (`member_user_id`, `created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='钱包流水表';


-- ---------------------------------------------------------------------
-- 支付单：打印订单支付与钱包充值共用一张表，用 biz_type + biz_id 区分业务
-- payment_no 即微信支付的 out_trade_no，回调与对账都按它定位
-- 同一业务单同时只允许一个进行中（0 待支付 / 1 支付中）的支付单，由服务层保证
-- ---------------------------------------------------------------------
CREATE TABLE `pay_payment`
(
    `id`               bigint       NOT NULL COMMENT '主键',
    `payment_no`       varchar(32)  NOT NULL COMMENT '商户支付单号（雪花），全局唯一，即微信 out_trade_no',
    `biz_type`         varchar(32)  NOT NULL COMMENT '业务类型：print_order 打印订单 / wallet_recharge 钱包充值',
    `biz_id`           bigint       NOT NULL COMMENT '业务单 ID：打印订单为 print_order.id，充值为 pay_wallet.id',
    `member_user_id`   bigint       NOT NULL COMMENT '付款会员 ID，对应 member_user.id',
    `channel`          varchar(16)  NOT NULL COMMENT '渠道大类：wechat 微信 / alipay 支付宝 / balance 钱包余额',
    `channel_code`     varchar(32)  NOT NULL DEFAULT '' COMMENT '渠道配置编码，对应 pay_channel_config.code；balance 为空',
    `scene`            varchar(16)  NOT NULL DEFAULT '' COMMENT '支付场景：mini 小程序 / h5 / wap；balance 为空',
    `client_ip`        varchar(50)  NOT NULL DEFAULT '' COMMENT '下单客户端 IP，微信 H5 支付必填',
    `amount`           bigint       NOT NULL COMMENT '应付金额（分）',
    `paid_amount`      bigint       NOT NULL DEFAULT '0' COMMENT '实付金额（分）',
    `refund_amount`    bigint       NOT NULL DEFAULT '0' COMMENT '已退款金额（分），不得超过 paid_amount',
    `status`           tinyint      NOT NULL DEFAULT '0' COMMENT '状态：0 待支付 / 1 支付中 / 2 已支付 / 3 支付失败 / 4 已关闭 / 5 部分退款 / 6 全额退款',
    `channel_trade_no` varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道交易号（微信 transaction_id）',
    `open_id`          varchar(64)  NOT NULL DEFAULT '' COMMENT '微信 openid，JSAPI 支付必需',
    `expire_time`      datetime     NULL COMMENT '支付超时时间，默认创建后 30 分钟',
    `success_time`     datetime     NULL COMMENT '支付成功时间',
    `close_time`       datetime     NULL COMMENT '关闭时间（超时关单）',
    `remark`           varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `creator`          varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`       datetime     NULL COMMENT '创建时间',
    `updater`          varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`       datetime     NULL COMMENT '更新时间',
    `deleted_at`       datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_payment_no` (`payment_no`),
    KEY `idx_biz` (`biz_type`, `biz_id`),
    KEY `idx_channel_code` (`channel_code`),
    KEY `idx_channel_trade_no` (`channel_trade_no`),
    KEY `idx_member` (`member_user_id`, `created_at`),
    KEY `idx_status_expire` (`status`, `expire_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='会员支付单表';


-- ---------------------------------------------------------------------
-- 退款单：一次支付可以有多笔退款，累计退款金额不得超过实付金额
-- 微信渠道走异步回调，钱包渠道直接入账并立即置成功
-- ---------------------------------------------------------------------
CREATE TABLE `pay_refund`
(
    `id`                bigint       NOT NULL COMMENT '主键',
    `refund_no`         varchar(32)  NOT NULL COMMENT '退款单号（雪花），全局唯一',
    `payment_id`        bigint       NOT NULL COMMENT '支付单 ID，对应 pay_payment.id',
    `payment_no`        varchar(32)  NOT NULL COMMENT '支付单号（冗余，便于对账）',
    `biz_type`          varchar(32)  NOT NULL COMMENT '业务类型，与支付单一致',
    `biz_id`            bigint       NOT NULL COMMENT '业务单 ID，与支付单一致',
    `member_user_id`    bigint       NOT NULL COMMENT '收款会员 ID，对应 member_user.id',
    `channel`           varchar(16)  NOT NULL COMMENT '渠道大类：与支付单一致，wechat / alipay / balance',
    `channel_code`      varchar(32)  NOT NULL DEFAULT '' COMMENT '渠道配置编码，与支付单一致；balance 为空',
    `amount`            bigint       NOT NULL COMMENT '退款金额（分）',
    `reason`            varchar(255) NOT NULL DEFAULT '' COMMENT '退款原因',
    `status`            tinyint      NOT NULL DEFAULT '0' COMMENT '状态：0 待退款 / 1 退款中 / 2 退款成功 / 3 退款失败 / 4 已取消',
    `channel_refund_no` varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道退款单号（微信 refund_id）',
    `success_time`      datetime     NULL COMMENT '退款成功时间',
    `fail_reason`       varchar(255) NOT NULL DEFAULT '' COMMENT '失败原因',
    `operator_id`       bigint       NOT NULL DEFAULT '0' COMMENT '发起退款的后台操作人，对应 system_users.id；0 表示系统自动（如打印失败自动退款）',
    `creator`           varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`        datetime     NULL COMMENT '创建时间',
    `updater`           varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`        datetime     NULL COMMENT '更新时间',
    `deleted_at`        datetime     NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_refund_no` (`refund_no`),
    KEY `idx_payment` (`payment_id`),
    KEY `idx_biz` (`biz_type`, `biz_id`),
    KEY `idx_member` (`member_user_id`, `created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='会员退款单表';


-- ---------------------------------------------------------------------
-- 支付渠道回调日志：微信回调无法重放，必须落原始报文（脱敏）便于排查
-- 日志表只留 created_at，不需要审计四字段
-- ---------------------------------------------------------------------
CREATE TABLE `pay_notify_log`
(
    `id`            bigint       NOT NULL COMMENT '主键',
    `channel`       varchar(16)  NOT NULL COMMENT '渠道：wechat',
    `notify_type`   varchar(16)  NOT NULL COMMENT '回调类型：payment 支付 / refund 退款',
    `out_trade_no`  varchar(64)  NOT NULL DEFAULT '' COMMENT '商户单号（支付回调用 out_trade_no，退款回调用 out_refund_no）',
    `channel_no`    varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道单号（transaction_id / refund_id）',
    `body`          text         NULL COMMENT '原始报文（脱敏后存储）',
    `verified`      tinyint      NOT NULL DEFAULT '0' COMMENT '验签结果：0 失败 / 1 通过',
    `handle_status` tinyint      NOT NULL DEFAULT '0' COMMENT '处理结果：0 待处理 / 1 处理成功 / 2 幂等忽略 / 3 处理失败',
    `handle_msg`    varchar(255) NOT NULL DEFAULT '' COMMENT '处理说明',
    `created_at`    datetime     NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_out_trade_no` (`out_trade_no`),
    KEY `idx_created_at` (`created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='支付渠道回调日志表';


-- =====================================================================
-- 会员卡
-- =====================================================================


-- ---------------------------------------------------------------------
-- 会员卡：独立的折扣载体，卡自带等级与折扣，不与 member_level 等级体系关联
-- 等级折扣唯一来源是本表 discount_percent（100 无折扣 / 95 表示 95 折）；
-- 卡记录不删除，停用或过期改 status，避免 uk_member 被删除记录占住
-- ---------------------------------------------------------------------
CREATE TABLE `member_card`
(
    `id`               bigint       NOT NULL COMMENT '主键',
    `member_user_id`   bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id，一个会员一张卡',
    `card_no`          varchar(32)  NOT NULL COMMENT '卡号，服务端生成（C + 雪花后 10 位），全局唯一',
    `card_level`       varchar(32)  NOT NULL DEFAULT 'normal' COMMENT '卡等级编码：normal 普通 / silver 银卡 / gold 金卡 / black 黑金，可按业务自定义',
    `card_name`        varchar(64)  NOT NULL COMMENT '卡名称，展示用，如「金卡会员」',
    `discount_percent` int          NOT NULL DEFAULT '100' COMMENT '卡自带折扣百分比：100 无折扣 / 95 表示 95 折，支付折扣的唯一来源',
    `card_type`        varchar(16)  NOT NULL DEFAULT 'level' COMMENT '卡类型：level 折扣卡（预留 stored 储值卡）',
    `source`           varchar(16)  NOT NULL DEFAULT 'admin' COMMENT '开卡来源：register 注册赠送 / admin 后台发放 / purchase 购买（预留）',
    `status`           tinyint      NOT NULL DEFAULT '1' COMMENT '状态：0 停用 / 1 正常 / 2 已过期',
    `start_time`       datetime     NULL COMMENT '开卡时间',
    `expire_time`      datetime     NULL COMMENT '到期时间，NULL 表示永久有效；到期由定时任务置 2，折扣随之失效',
    `remark`           varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `creator`          varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
    `created_at`       datetime     NULL COMMENT '创建时间',
    `updater`          varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
    `updated_at`       datetime     NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_card_no` (`card_no`),
    UNIQUE KEY `uk_member` (`member_user_id`),
    KEY `idx_card_level` (`card_level`),
    KEY `idx_status_expire` (`status`, `expire_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='会员卡表';
