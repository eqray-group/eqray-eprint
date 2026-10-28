package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// PrintFileTable 是打印文件表名常量。
const PrintFileTable = "print_file"

// PrintFileModel 对应 print_file，只做 infra_file 的打印元数据扩展。
// 源文件与可打印 PDF 都落在 infra_file，本表不存 URL，避免两套存储。
type PrintFileModel struct {
	baseModel.BaseModel
	MemberUserID  types.ID        `db:"member_user_id" json:"memberUserId"`
	InfraFileID   types.ID        `db:"infra_file_id" json:"infraFileId"`
	PDFFileID     types.ID        `db:"pdf_file_id" json:"pdfFileId"`
	FileName      string          `db:"file_name" json:"fileName"`
	FileExt       string          `db:"file_ext" json:"fileExt"`
	FileSize      int64           `db:"file_size" json:"fileSize"`
	FileSha256    string          `db:"file_sha256" json:"fileSha256"`
	PageCount     int             `db:"page_count" json:"pageCount"`
	ConvertStatus int8            `db:"convert_status" json:"convertStatus"`
	ConvertError  string          `db:"convert_error" json:"convertError"`
	ExpireTime    *types.DateTime `db:"expire_time" json:"expireTime"`
}

// TableName 返回表名。
func (m *PrintFileModel) TableName() string {
	return PrintFileTable
}
