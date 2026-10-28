package migrations

import (
	"eqray-eprint/app/facades"
)

// M20261027000016CreatePrintFileTable 建 print_file 打印文件表。
// 不重复建文件存储：源文件与可打印 PDF 都落在 infra_file，
// 本表只保存打印元数据（页数、转换状态、下载校验用的 SHA-256）。
// 打印参数属于订单明细，本表不存，避免两处真相。
type M20261027000016CreatePrintFileTable struct{}

// Signature 返回迁移唯一标识，必须与文件名（去掉 .go）一致。
func (r *M20261027000016CreatePrintFileTable) Signature() string {
	return "20261027000016_create_print_file_table"
}

// Up 建打印文件表；源文件一对一（uk_infra_file），按会员与过期时间建索引。
// 无参数；表已存在时直接返回 nil。
func (r *M20261027000016CreatePrintFileTable) Up() error {
	if facades.Schema().HasTable("print_file") {
		return nil
	}

	return facades.DB().Statement(`
CREATE TABLE print_file (
  id             bigint       NOT NULL COMMENT '主键',
  member_user_id bigint       NOT NULL COMMENT '会员 ID，对应 member_user.id',
  infra_file_id  bigint       NOT NULL COMMENT '源文件，对应 infra_file.id',
  pdf_file_id    bigint       NOT NULL DEFAULT 0 COMMENT '可打印 PDF，对应 infra_file.id；源文件本身是 PDF 时与 infra_file_id 相同',
  file_name      varchar(255) NOT NULL COMMENT '文件名快照（含扩展名）',
  file_ext       varchar(16)  NOT NULL COMMENT '扩展名：pdf / jpg / png / docx',
  file_size      bigint       NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
  file_sha256    char(64)     NOT NULL DEFAULT '' COMMENT '可打印文件 SHA-256，客户端下载后校验',
  page_count     int          NOT NULL DEFAULT 0 COMMENT '可打印 PDF 总页数，由服务端计算',
  convert_status tinyint      NOT NULL DEFAULT 0 COMMENT '转换状态：0 待处理 / 1 转换中 / 2 可打印 / 3 转换失败',
  convert_error  varchar(255) NOT NULL DEFAULT '' COMMENT '转换失败原因',
  expire_time    datetime     NULL COMMENT '文件清理时间，到期由定时任务删除存储与记录',
  creator        varchar(64)  NOT NULL DEFAULT '' COMMENT '创建人',
  created_at     datetime     NULL COMMENT '创建时间',
  updater        varchar(64)  NOT NULL DEFAULT '' COMMENT '更新人',
  updated_at     datetime     NULL COMMENT '更新时间',
  deleted_at     datetime     NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_infra_file (infra_file_id),
  KEY idx_member (member_user_id),
  KEY idx_expire (expire_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='打印文件表'`)
}

// Down 删除打印文件表，无参数，返回删除结果。
func (r *M20261027000016CreatePrintFileTable) Down() error {
	return facades.Schema().DropIfExists("print_file")
}
