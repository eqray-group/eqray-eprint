package controller

import (
	"github.com/goravel/framework/contracts/http"

	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/types"
	"eqray-eprint/app/modules/print/service"
)

// FileController 提供打印文件的后台只读查询与统计概览。
type FileController struct {
	files      *service.FileService
	statistics *service.StatisticsService
}

// NewFileController 构造打印文件后台控制器。
func NewFileController() *FileController {
	return &FileController{
		files:      service.NewFileService(),
		statistics: service.NewStatisticsService(),
	}
}

// Index 分页查询打印文件，可按会员、文件名、转换状态筛选。
// 只读：文件内容属于会员，后台不做下载与删除，清理交给保留期任务。
func (c *FileController) Index(ctx http.Context) http.Response {
	page, err := param.FromQuery(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	q := service.AdminFileQuery{
		FileName: ctx.Request().Query("fileName"),
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	}
	if raw := ctx.Request().QueryInt("memberId", 0); raw > 0 {
		q.MemberID = types.ID(raw)
	}
	if status := ctx.Request().QueryInt("convertStatus", -1); status >= 0 {
		v := int8(status)
		q.ConvertStatus = &v
	}

	list, total, err := c.files.AdminPage(ctx.Context(), q)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Page(ctx, list, total)
}

// Statistics 返回打印业务概览，供控制台卡片使用。
func (c *FileController) Statistics(ctx http.Context) http.Response {
	overview, err := c.statistics.Overview(ctx.Context())
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, overview)
}
