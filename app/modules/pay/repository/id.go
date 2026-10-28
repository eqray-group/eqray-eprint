package repository

import (
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
)

// newID 生成雪花主键；仓储层写非 BaseModel 的表（如回调日志）时使用。
func newID() types.ID {
	return support.NewID()
}
