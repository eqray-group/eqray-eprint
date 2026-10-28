package controller

import (
	"github.com/goravel/framework/contracts/http"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/param"
	"eqray-eprint/app/modules/base/result"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/print/service"
)

// 上传字段名，前端 multipart 用 file。
const uploadField = "file"

// FileController 提供会员端打印文件的上传、查询与重试转换。
type FileController struct {
	service *service.FileService
}

// NewFileController 构造文件控制器。
func NewFileController() *FileController {
	return &FileController{service: service.NewFileService()}
}

// Upload 上传打印文件并同步完成可打印化。
// 转换失败不影响上传结果，失败原因在返回的 convertError 里。
func (c *FileController) Upload(ctx http.Context) http.Response {
	file, err := ctx.Request().File(uploadField)
	if err != nil || file == nil {
		return result.Fail(ctx, baseErrs.WithMsg(baseErrs.BadRequest, "请选择要上传的文件"))
	}

	record, err := c.service.Upload(ctx.Context(), support.MemberID(ctx.Context()), file)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, record)
}

// Show 查询文件状态，前端轮询转换结果用。
func (c *FileController) Show(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	record, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, record)
}

// Convert 手动重试转换：适用于服务器随后安装了 LibreOffice 或源文件被替换的场景。
func (c *FileController) Convert(ctx http.Context) http.Response {
	id, err := param.PathID(ctx)
	if err != nil {
		return result.Fail(ctx, err)
	}

	if _, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()), id); err != nil {
		return result.Fail(ctx, err)
	}

	if err := c.service.Convert(ctx.Context(), id); err != nil {
		return result.Fail(ctx, err)
	}

	record, err := c.service.Detail(ctx.Context(), support.MemberID(ctx.Context()), id)
	if err != nil {
		return result.Fail(ctx, err)
	}

	return result.Success(ctx, record)
}
