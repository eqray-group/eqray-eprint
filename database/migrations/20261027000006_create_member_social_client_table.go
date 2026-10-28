package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000006CreateMemberSocialClientTable 建 member_social_client 第三方登录渠道配置表。
// 渠道凭据只在后台配置（会员管理 → 社交客户端），登录时按 social_type 读取，
// 不再依赖环境变量或配置文件。
type M20261027000006CreateMemberSocialClientTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000006CreateMemberSocialClientTable) Signature() string {
	return "20261027000006_create_member_social_client_table"
}

// Up 建 member_social_client 表，social_type 唯一保证一个渠道只有一份配置。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000006CreateMemberSocialClientTable) Up() error {
	if facades.Schema().HasTable("member_social_client") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE member_social_client (
  id            bigint       NOT NULL COMMENT '主键',
  name          varchar(64)  NOT NULL COMMENT '应用名',
  social_type   varchar(32)  NOT NULL COMMENT '渠道：wechat_mini / wechat_mp / wechat_open',
  client_id     varchar(128) NOT NULL COMMENT '客户端 ID（微信 appid）',
  client_secret varchar(256) NOT NULL COMMENT '客户端密钥（微信 appsecret）',
  agent_id      varchar(64)  NOT NULL DEFAULT '' COMMENT '代理编号（钉钉/企微等渠道）',
  status        tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0 禁用 / 1 启用',
  creator       varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at    datetime     NULL COMMENT '创建时间',
  updater       varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at    datetime     NULL COMMENT '更新时间',
  deleted_at    datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_social_type (social_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员第三方登录渠道配置表'`)
}

// Down 删除第三方登录渠道配置表，无参数，返回删除结果。
func (r *M20261027000006CreateMemberSocialClientTable) Down() error {
	return facades.Schema().DropIfExists("member_social_client")
}
