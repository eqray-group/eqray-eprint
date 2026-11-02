package repository

import (
	"context"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
)

// systemUserTable 是后台账号表名。print 模块只读取账号展示字段，
// 不引用 system 模块的 Service，避免跨模块 Go 依赖（docs 3.3）。
const systemUserTable = "system_users"

// SystemUserBrief 是绑定门店成员时需要展示的账号摘要。
type SystemUserBrief struct {
	ID       types.ID `db:"id"`
	Username string   `db:"username"`
	Nickname string   `db:"nickname"`
	Mobile   *string  `db:"mobile"`
}

// SystemUserReader 提供后台账号的只读查询，仅供门店成员绑定使用。
type SystemUserReader struct{}

// NewSystemUserReader 构造后台账号只读仓储。
func NewSystemUserReader() *SystemUserReader {
	return &SystemUserReader{}
}

// SearchActive 按用户名 / 昵称 / 手机号模糊搜索启用中的账号。
// 参数 keyword 为空时返回空列表；limit 控制条数，避免下拉框一次塞入过多选项。
func (r *SystemUserReader) SearchActive(ctx context.Context, keyword string, limit int) ([]SystemUserBrief, error) {
	list := make([]SystemUserBrief, 0)
	if keyword == "" {
		return list, nil
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	like := "%" + keyword + "%"
	if err := facades.DB().WithContext(ctx).Table(systemUserTable).
		Select("id", "username", "nickname", "mobile").
		WhereNull("deleted_at").
		Where("status", 1).
		WhereRaw("(username LIKE ? OR nickname LIKE ? OR mobile LIKE ?)", []any{like, like, like}).
		OrderByDesc("id").
		Limit(uint64(limit)).
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// ListByIDs 批量查询账号摘要，用于门店成员列表回填用户名与昵称。
// 参数 ids 为空时直接返回空列表，不发起查询。
func (r *SystemUserReader) ListByIDs(ctx context.Context, ids []types.ID) ([]SystemUserBrief, error) {
	list := make([]SystemUserBrief, 0)
	if len(ids) == 0 {
		return list, nil
	}

	values := make([]any, 0, len(ids))
	for _, id := range ids {
		values = append(values, id)
	}

	if err := facades.DB().WithContext(ctx).Table(systemUserTable).
		Select("id", "username", "nickname", "mobile").
		WhereIn("id", values).
		WhereNull("deleted_at").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}
