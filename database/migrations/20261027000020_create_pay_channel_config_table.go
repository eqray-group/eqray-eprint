package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000020CreatePayChannelConfigTable 建 pay_channel_config 支付渠道配置表。
// 平台统一收款：所有门店订单进同一个商户号，不做门店独立商户号与分账，
// 因此不挂 shop_id，唯一键固定为 (channel, scene)。
// 微信三条配置：wechat_mini 小程序 / wechat_mp 公众号（微信内置浏览器）/ wechat_h5 手机浏览器；
// 支付宝后续接入只加配置（alipay + wap），不改表结构。
// 密钥类字段落库前必须加密，接口一律不返回明文，编辑时留空表示不修改。
type M20261027000020CreatePayChannelConfigTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000020CreatePayChannelConfigTable) Signature() string {
	return "20261027000020_create_pay_channel_config_table"
}

// Up 建支付渠道配置表，code 全局唯一，(channel, scene) 保证同一场景只有一份配置。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000020CreatePayChannelConfigTable) Up() error {
	if facades.Schema().HasTable("pay_channel_config") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE pay_channel_config (
  id                   bigint       NOT NULL COMMENT '主键',
  code                 varchar(32)  NOT NULL COMMENT '配置编码，全局唯一：wechat_mini / wechat_mp / wechat_h5 / alipay_wap',
  name                 varchar(64)  NOT NULL COMMENT '配置名称，如「微信小程序支付」',
  channel              varchar(16)  NOT NULL COMMENT '渠道大类：wechat 微信 / alipay 支付宝',
  scene                varchar(16)  NOT NULL COMMENT '支付场景：mini 小程序 / jsapi 公众号（微信内置浏览器）/ h5 手机浏览器 / wap 支付宝手机网站 / qr 扫码',
  app_id               varchar(64)  NOT NULL COMMENT '应用 ID：微信 appid / 支付宝 appid',
  mch_id               varchar(64)  NOT NULL DEFAULT '' COMMENT '商户号：微信 mchid / 支付宝 PID',
  serial_no            varchar(64)  NOT NULL DEFAULT '' COMMENT '微信 APIv3 证书序列号',
  private_key          text         NULL COMMENT '商户私钥（加密存储）',
  api_key              varchar(256) NOT NULL DEFAULT '' COMMENT '微信 APIv3 密钥（加密存储）',
  platform_public_key  text         NULL COMMENT '平台证书 / 支付宝公钥（加密存储）',
  gateway              varchar(128) NOT NULL DEFAULT '' COMMENT '网关地址，支付宝必填、微信可空',
  notify_url           varchar(255) NOT NULL DEFAULT '' COMMENT '支付回调地址，空表示走默认路由',
  return_url           varchar(255) NOT NULL DEFAULT '' COMMENT 'H5 支付完成后跳回的页面地址',
  extra                json         NULL COMMENT '其他渠道参数 JSON，如 H5 授权域名、小程序原始 ID',
  status               tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 停用 / 1 启用',
  sort                 int          NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  remark               varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  creator              varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at           datetime     NULL COMMENT '创建时间',
  updater              varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at           datetime     NULL COMMENT '更新时间',
  deleted_at           datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_code (code),
  UNIQUE KEY uk_channel_scene (channel, scene),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付渠道配置表'`)
}

// Down 删除支付渠道配置表，无参数，返回删除结果。
func (r *M20261027000020CreatePayChannelConfigTable) Down() error {
	return facades.Schema().DropIfExists("pay_channel_config")
}
