// Package service 存放 infra 模块的业务逻辑。
package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/s3"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/types"
	infraErrs "eqray-eprint/app/modules/infra/errs"
	"eqray-eprint/app/modules/infra/model"
	"eqray-eprint/app/modules/infra/repository"
)

// storageConfigMu 保护运行时的 S3 磁盘配置写入，viper 的 Set 不是并发安全的。
var storageConfigMu sync.Mutex

const (
	// storageDisk 与 config/filesystems.go 中的 public 磁盘一致，文件落在 storage/app/public 下。
	storageDisk = "public"
	// uploadDir 是文件在磁盘中的根目录，上传时再拼上「年/月/日」子目录。
	uploadDir = "uploads"
	// defaultMaxFileSizeMB 是单个文件的默认上限（MB），可由 infra_config 的
	// infra.upload.max_size 覆盖。
	defaultMaxFileSizeMB = 5
	// uploadMaxSizeKey 指向参数配置中的上传大小上限。
	uploadMaxSizeKey = "infra.upload.max_size"
)

// allowedExtensions 白名单：只放行常见图片，svg 因可内嵌脚本不在此列。
var allowedExtensions = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"gif":  true,
	"webp": true,
	"bmp":  true,
}

// FileService 负责 infra_file 的上传与删除（docs 6.5）。
type FileService struct {
	files       *repository.FileRepository
	configs     *ConfigService
	fileConfigs *FileConfigService
}

func NewFileService() *FileService {
	return &FileService{
		files:       repository.NewFileRepository(),
		configs:     NewConfigService(),
		fileConfigs: NewFileConfigService(),
	}
}

// Upload 上传后台通用图片：走默认图片白名单与参数配置里的体积上限。
func (s *FileService) Upload(ctx context.Context, file filesystem.File) (*model.InfraFileModel, error) {
	return s.UploadAs(ctx, file, allowedExtensions, 0)
}

// storageDriver 按存储目标返回磁盘驱动；S3 类型在运行时按数据库里的配置构建驱动。
func storageDriver(ctx context.Context, target *FileStorageTarget) (filesystem.Driver, error) {
	if target == nil || target.Storage != fileConfigS3Storage || target.S3 == nil {
		return facades.Storage().Disk(storageDisk), nil
	}

	conf := target.S3
	disk := fmt.Sprintf("infra_s3_%d", target.ConfigID)

	// s3 驱动从 filesystems.disks.<disk> 读取参数，这里把库里的配置写入运行时配置
	storageConfigMu.Lock()
	defer storageConfigMu.Unlock()

	facades.Config().Add("filesystems.disks."+disk, map[string]any{
		"driver":         "custom",
		"key":            conf.AccessKey,
		"secret":         conf.SecretKey,
		"region":         conf.Region,
		"bucket":         conf.Bucket,
		"endpoint":       conf.Endpoint,
		"url":            conf.URL,
		"use_path_style": conf.UsePathStyle,
		"via": func() (filesystem.Driver, error) {
			return s3.NewS3(ctx, facades.Config(), disk)
		},
	})

	driver, err := s3.NewS3(ctx, facades.Config(), disk)
	if err != nil {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "S3 存储配置不可用："+err.Error())
	}

	return driver, nil
}

// FileListQuery 是文件列表的查询条件。
type FileListQuery struct {
	Keyword  string
	PageNo   int
	PageSize int
}

// Page 分页查询已上传的文件，供文件管理页使用。
func (s *FileService) Page(ctx context.Context, q FileListQuery) ([]model.InfraFileModel, int64, error) {
	return s.files.Paginate(ctx, repository.FileFilter{Keyword: q.Keyword}, q.PageNo, q.PageSize)
}

// Delete 软删除文件记录，同时删除磁盘上的物理文件；物理文件缺失不阻断删除。
func (s *FileService) Delete(ctx context.Context, id types.ID) error {
	file, err := s.files.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return baseErrs.New(infraErrs.InfraFileNotExists, id)
	}

	if err := s.files.Delete(ctx, file); err != nil {
		return err
	}

	s.deletePhysical(ctx, file)

	return nil
}

// deletePhysical 按文件所属配置选择磁盘删除物理文件，配置缺失或不可用时回落到本地磁盘。
func (s *FileService) deletePhysical(ctx context.Context, file *model.InfraFileModel) {
	target := &FileStorageTarget{Storage: fileConfigLocalStorage, Dir: uploadDir}
	if file.ConfigID > 0 {
		if t, err := s.fileConfigs.Target(ctx, file.ConfigID); err == nil && t != nil {
			target = t
		}
	}

	driver, err := storageDriver(ctx, target)
	if err != nil {
		driver = facades.Storage().Disk(storageDisk)
	}

	_ = driver.Delete(file.Path)
}
