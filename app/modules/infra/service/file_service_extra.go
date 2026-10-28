package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/filesystem"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	infraErrs "eqray-eprint/app/modules/infra/errs"
	"eqray-eprint/app/modules/infra/model"
)

// UploadAs 按调用方给出的扩展名白名单与体积上限上传文件。
// 参数 file 为上传文件，allowed 为小写扩展名集合，maxMB 为上限（0 表示读参数配置）。
// 业务模块（如打印文件）需要 pdf/docx 时用它，避免把图片白名单改宽影响后台。
func (s *FileService) UploadAs(ctx context.Context, file filesystem.File, allowed map[string]bool, maxMB int) (*model.InfraFileModel, error) {
	if file == nil {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "请选择要上传的文件")
	}

	ext := strings.ToLower(strings.TrimPrefix(file.GetClientOriginalExtension(), "."))
	if !allowed[ext] {
		return nil, baseErrs.New(infraErrs.InfraFileTypeNotAllowed)
	}

	size, err := file.Size()
	if err != nil {
		return nil, baseErrs.New(infraErrs.InfraFileUploadFailed)
	}
	if size <= 0 {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "文件内容为空")
	}

	if maxMB <= 0 {
		maxMB = s.configs.GetInt(ctx, uploadMaxSizeKey, defaultMaxFileSizeMB)
	}
	if maxMB <= 0 {
		maxMB = defaultMaxFileSizeMB
	}
	if size > int64(maxMB)<<20 {
		return nil, baseErrs.New(infraErrs.InfraFileTooLarge, maxMB)
	}

	target, err := s.fileConfigs.Active(ctx)
	if err != nil {
		return nil, err
	}

	driver, err := storageDriver(ctx, target)
	if err != nil {
		return nil, err
	}

	dir := s.dayDir(target)
	name := fmt.Sprintf("%d.%s", support.NewID(), ext)

	path, err := driver.PutFileAs(dir, file, name)
	if err != nil {
		return nil, baseErrs.New(infraErrs.InfraFileUploadFailed)
	}

	mimeType, _ := file.MimeType()

	m := &model.InfraFileModel{
		ConfigID: target.ConfigID,
		Name:     file.GetClientOriginalName(),
		Path:     path,
		Url:      s.urlOf(driver, target, path),
		Type:     mimeType,
		Size:     size,
	}
	if err := s.files.Create(ctx, m); err != nil {
		_ = driver.Delete(path)
		return nil, err
	}

	return m, nil
}

// StoreLocal 把服务端生成的文件（如转换后的 PDF）写入存储并登记 infra_file。
// 参数 localPath 为本地临时文件路径，name 为原始文件名，ext 为扩展名（不带点）。
func (s *FileService) StoreLocal(ctx context.Context, localPath, name, ext string) (*model.InfraFileModel, error) {
	content, err := readLocalFile(localPath)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "生成的文件内容为空")
	}

	target, err := s.fileConfigs.Active(ctx)
	if err != nil {
		return nil, err
	}
	driver, err := storageDriver(ctx, target)
	if err != nil {
		return nil, err
	}

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	dir := s.dayDir(target)
	path := fmt.Sprintf("%s/%d.%s", dir, support.NewID(), ext)

	if err := driver.Put(path, string(content)); err != nil {
		return nil, baseErrs.New(infraErrs.InfraFileUploadFailed)
	}

	m := &model.InfraFileModel{
		ConfigID: target.ConfigID,
		Name:     name,
		Path:     path,
		Url:      s.urlOf(driver, target, path),
		Type:     mimeTypeOf(ext),
		Size:     int64(len(content)),
	}
	if err := s.files.Create(ctx, m); err != nil {
		_ = driver.Delete(path)
		return nil, err
	}

	return m, nil
}

// ReadAll 读取文件内容，供业务模块做转换、校验等需要完整字节的场景。
// 参数 m 为 infra_file 记录；记录不存在或存储不可用时返回业务错误。
func (s *FileService) ReadAll(ctx context.Context, m *model.InfraFileModel) ([]byte, error) {
	if m == nil || m.Path == "" {
		return nil, baseErrs.New(infraErrs.InfraFileNotExists)
	}

	var target *FileStorageTarget
	var err error
	if m.ConfigID > 0 {
		target, err = s.fileConfigs.Target(ctx, m.ConfigID)
	} else {
		target, err = s.fileConfigs.Active(ctx)
	}
	if err != nil {
		return nil, err
	}

	driver, err := storageDriver(ctx, target)
	if err != nil {
		return nil, err
	}

	content, err := driver.Get(m.Path)
	if err != nil {
		return nil, baseErrs.New(infraErrs.InfraFileNotExists)
	}

	return []byte(content), nil
}

// FindByID 查询文件记录，业务模块据此拿到存储路径与配置。
func (s *FileService) FindByID(ctx context.Context, id types.ID) (*model.InfraFileModel, error) {
	return s.files.FindByID(ctx, id)
}

// dayDir 返回按「年/月/日」分目录的上传目录。
func (s *FileService) dayDir(target *FileStorageTarget) string {
	dir := strings.Trim(target.Dir, "/")
	if dir == "" {
		dir = uploadDir
	}

	return dir + "/" + time.Now().Format("2006/01/02")
}

// urlOf 生成本地可访问地址：本地磁盘走 /storage 前缀，对象存储用公网地址。
func (s *FileService) urlOf(driver filesystem.Driver, target *FileStorageTarget, path string) string {
	if target.Storage == fileConfigS3Storage {
		return driver.Url(path)
	}

	return "/storage/" + path
}

// mimeTypeOf 按扩展名给出 MIME，仅覆盖常见类型，未知按二进制处理。
func mimeTypeOf(ext string) string {
	switch ext {
	case "pdf":
		return "application/pdf"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

// readLocalFile 读取服务端本地文件（转换产物都写在系统临时目录）。
func readLocalFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "读取本地文件失败："+err.Error())
	}

	return content, nil
}
