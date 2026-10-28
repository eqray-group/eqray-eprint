package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/model"
)

const fileTable = model.PrintFileTable

// PrintFileRepository 是 print_file 表的仓储（infra_file 的打印元数据扩展）。
type PrintFileRepository struct{}

// NewPrintFileRepository 构造打印文件仓储。
func NewPrintFileRepository() *PrintFileRepository {
	return &PrintFileRepository{}
}

// Create 新增打印文件记录。
func (r *PrintFileRepository) Create(ctx context.Context, m *model.PrintFileModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(fileTable).Insert(map[string]any{
		"id":             m.ID,
		"member_user_id": m.MemberUserID,
		"infra_file_id":  m.InfraFileID,
		"pdf_file_id":    m.PDFFileID,
		"file_name":      m.FileName,
		"file_ext":       m.FileExt,
		"file_size":      m.FileSize,
		"file_sha256":    m.FileSha256,
		"page_count":     m.PageCount,
		"convert_status": m.ConvertStatus,
		"convert_error":  m.ConvertError,
		"expire_time":    m.ExpireTime,
		"creator":        m.Creator,
		"created_at":     m.CreatedAt,
		"updater":        m.Updater,
		"updated_at":     m.UpdatedAt,
	})

	return err
}

// Update 更新转换结果与页数。
func (r *PrintFileRepository) Update(ctx context.Context, m *model.PrintFileModel) error {
	m.FillUpdate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(fileTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"pdf_file_id":    m.PDFFileID,
			"file_sha256":    m.FileSha256,
			"page_count":     m.PageCount,
			"convert_status": m.ConvertStatus,
			"convert_error":  m.ConvertError,
			"updater":        m.Updater,
			"updated_at":     m.UpdatedAt,
		})

	return err
}

// FindByID 按主键查询；未命中返回 (nil, nil)。
func (r *PrintFileRepository) FindByID(ctx context.Context, id types.ID) (*model.PrintFileModel, error) {
	var m model.PrintFileModel

	if err := facades.DB().WithContext(ctx).Table(fileTable).
		Where("id", id).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindByIDs 批量查询，返回 ID → 文件映射，避免下单时逐个查询。
func (r *PrintFileRepository) FindByIDs(ctx context.Context, ids []types.ID) (map[types.ID]model.PrintFileModel, error) {
	result := make(map[types.ID]model.PrintFileModel, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	list := make([]model.PrintFileModel, 0, len(ids))
	if err := facades.DB().WithContext(ctx).Table(fileTable).
		WhereIn("id", toAnySlice(ids)).
		WhereNull("deleted_at").
		Get(&list); err != nil {
		return nil, err
	}

	for _, item := range list {
		result[item.ID] = item
	}

	return result, nil
}

// ListExpired 查询已过保留期的文件，供清理任务使用。
// 参数 limit 为单批上限，避免一次删除过多拖慢数据库与存储。
func (r *PrintFileRepository) ListExpired(ctx context.Context, limit int) ([]model.PrintFileModel, error) {
	list := make([]model.PrintFileModel, 0)

	if err := facades.DB().WithContext(ctx).Table(fileTable).
		Where("expire_time", "<", types.Now()).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Limit(uint64(limit)).
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// SoftDelete 软删除文件记录；物理文件由 infra 模块负责删除。
func (r *PrintFileRepository) SoftDelete(ctx context.Context, id types.ID) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(fileTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"deleted_at": now,
			"updated_at": now,
		})

	return err
}

// FileFilter 是后台文件列表的筛选条件，零值表示不限制。
type FileFilter struct {
	MemberID      types.ID
	FileName      string
	ConvertStatus *int8
}

// Paginate 后台分页查询打印文件，按上传时间倒序。
func (r *PrintFileRepository) Paginate(ctx context.Context, filter FileFilter, pageNo, pageSize int) ([]model.PrintFileModel, int64, error) {
	list := make([]model.PrintFileModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(fileTable).WhereNull("deleted_at")
	if filter.MemberID > 0 {
		query = query.Where("member_user_id", filter.MemberID)
	}
	if filter.FileName != "" {
		query = query.WhereLike("file_name", "%"+filter.FileName+"%")
	}
	if filter.ConvertStatus != nil {
		query = query.Where("convert_status", *filter.ConvertStatus)
	}

	if err := query.OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// toAnySlice 把 ID 切片转成查询需要的 []any。
func toAnySlice(ids []types.ID) []any {
	values := make([]any, 0, len(ids))
	for _, id := range ids {
		values = append(values, id)
	}

	return values
}
