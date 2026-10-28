package controller

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/support"
)

// operatorName 返回当前后台账号标识，用于 updater 等留痕字段。
// 后台令牌里没有昵称，这里落账号 ID，避免为了一个字段多查一次库。
func operatorName(ctx http.Context) string {
	id := support.UserID(ctx.Context())
	if id <= 0 {
		return "system"
	}

	return strconv.FormatInt(id, 10)
}
