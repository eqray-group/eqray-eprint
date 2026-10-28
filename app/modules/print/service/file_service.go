package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/filesystem"
	"rsc.io/pdf"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	infraService "eqray-eprint/app/modules/infra/service"
	printEnums "eqray-eprint/app/modules/print/enums"
	printErrs "eqray-eprint/app/modules/print/errs"
	"eqray-eprint/app/modules/print/model"
	"eqray-eprint/app/modules/print/repository"
)

const (
	// printMaxFileMB 是打印文件体积上限：PDF/DOCX 通常比图片大，单独放宽到 20MB。
	printMaxFileMB = 20
	// fileKeep 是文件保留期，到期由定时任务清理，与订单取件保留期一致。
	fileKeep = 7 * 24 * time.Hour
	// convertTimeout 是 DOCX 转 PDF 的超时，LibreOffice 冷启动较慢。
	convertTimeout = 90 * time.Second
	// sofficeConfigKey 指向 LibreOffice 可执行文件路径，默认取 PATH 里的 soffice。
	sofficeConfigKey = "print.convert.soffice"
)

// printAllowedExtensions 是打印文件白名单；后台图片上传不受影响。
var printAllowedExtensions = map[string]bool{
	"pdf":  true,
	"docx": true,
	"jpg":  true,
	"jpeg": true,
	"png":  true,
}

// FileService 负责打印文件的上传登记与可打印化（页数统计、DOCX 转 PDF）。
type FileService struct {
	files *repository.PrintFileRepository
	infra *infraService.FileService
}

// NewFileService 构造打印文件服务。
func NewFileService() *FileService {
	return &FileService{
		files: repository.NewPrintFileRepository(),
		infra: infraService.NewFileService(),
	}
}

// Upload 上传打印文件：先落 infra_file，再登记 print_file 并立即尝试可打印化。
// 参数 memberID 为上传人，file 为 multipart 文件。
// 转换失败不算上传失败：记录会保留并标记失败原因，前端可提示重试。
func (s *FileService) Upload(ctx context.Context, memberID types.ID, file filesystem.File) (*model.PrintFileModel, error) {
	if file == nil {
		return nil, baseErrs.WithMsg(baseErrs.BadRequest, "请选择要上传的文件")
	}

	ext := strings.ToLower(strings.TrimPrefix(file.GetClientOriginalExtension(), "."))
	if !printAllowedExtensions[ext] {
		return nil, baseErrs.New(printErrs.FileUnsupported, ext)
	}

	infraFile, err := s.infra.UploadAs(ctx, file, printAllowedExtensions, printMaxFileMB)
	if err != nil {
		return nil, err
	}

	expire := types.DateTime(types.Now().Time().Add(fileKeep))
	record := &model.PrintFileModel{
		MemberUserID:  memberID,
		InfraFileID:   infraFile.ID,
		PDFFileID:     0,
		FileName:      infraFile.Name,
		FileExt:       ext,
		FileSize:      infraFile.Size,
		ConvertStatus: printEnums.ConvertStatusPending,
		ExpireTime:    &expire,
	}
	record.ID = support.NewID()

	if err := s.files.Create(ctx, record); err != nil {
		return nil, err
	}

	// PDF / 图片转换在毫秒级，同步做完；DOCX 要拉起 LibreOffice，放后台避免请求超时
	if ext == "docx" {
		go s.convertAsync(record.ID)
	} else {
		_ = s.Convert(ctx, record.ID)
	}

	return s.Detail(ctx, memberID, record.ID)
}

// Convert 把文件转成可打印状态：PDF 统计页数，图片按单页，DOCX 走 LibreOffice。
// 参数 id 为 print_file 主键；重复调用是幂等的，已可打印直接返回。
func (s *FileService) Convert(ctx context.Context, id types.ID) error {
	record, err := s.files.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if record == nil {
		return baseErrs.New(printErrs.FileNotExists)
	}
	if record.ConvertStatus == printEnums.ConvertStatusReady {
		return nil
	}

	infraFile, err := s.infra.FindByID(ctx, record.InfraFileID)
	if err != nil {
		return err
	}
	if infraFile == nil {
		return s.markFailed(ctx, record, "源文件不存在或已被清理")
	}

	content, err := s.infra.ReadAll(ctx, infraFile)
	if err != nil {
		return s.markFailed(ctx, record, "读取源文件失败")
	}

	switch record.FileExt {
	case "pdf":
		pages, err := countPDFPages(content)
		if err != nil {
			return s.markFailed(ctx, record, "PDF 解析失败："+err.Error())
		}

		return s.markReady(ctx, record, record.InfraFileID, pages, "")
	case "jpg", "jpeg", "png":
		// 图片包成单页 A4 PDF，保证客户端只需处理 PDF 一种格式
		pdfContent, err := imageToPDF(content, record.FileExt)
		if err != nil {
			return s.markFailed(ctx, record, "图片转 PDF 失败："+err.Error())
		}

		pdfPath, err := writeTempPDF(pdfContent)
		if err != nil {
			return s.markFailed(ctx, record, err.Error())
		}
		defer os.Remove(pdfPath)

		pdfFile, err := s.infra.StoreLocal(ctx, pdfPath, pdfNameOf(record.FileName), "pdf")
		if err != nil {
			return s.markFailed(ctx, record, "转换后的 PDF 入库失败")
		}

		return s.markReady(ctx, record, pdfFile.ID, 1, pdfFile.Url)
	case "docx":
		pdfPath, err := convertDocxToPDF(content, record.FileName)
		if err != nil {
			return s.markFailed(ctx, record, err.Error())
		}
		defer os.Remove(pdfPath)

		pdfFile, err := s.infra.StoreLocal(ctx, pdfPath, pdfNameOf(record.FileName), "pdf")
		if err != nil {
			return s.markFailed(ctx, record, "转换后的 PDF 入库失败")
		}

		pdfContent, err := os.ReadFile(pdfPath)
		if err != nil {
			return s.markFailed(ctx, record, "读取转换结果失败")
		}
		pages, err := countPDFPages(pdfContent)
		if err != nil {
			return s.markFailed(ctx, record, "转换结果解析失败："+err.Error())
		}

		return s.markReady(ctx, record, pdfFile.ID, pages, pdfFile.Url)
	default:
		return s.markFailed(ctx, record, "不支持的文件类型")
	}
}

// convertAsync 在后台执行 DOCX 转换，使用独立上下文，不受请求结束影响。
// 转换结果写回文件记录，前端轮询 GET /files/{id} 获取状态。
func (s *FileService) convertAsync(id types.ID) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if err := s.Convert(ctx, id); err != nil {
		facades.Log().Errorf("打印文件转换失败 id=%d err=%v", id, err)
	}
}

// AdminFileQuery 是后台文件列表的查询条件。
type AdminFileQuery struct {
	MemberID      types.ID
	FileName      string
	ConvertStatus *int8
	PageNo        int
	PageSize      int
}

// AdminPage 后台分页查询打印文件，用于排查转换失败与占用量。
func (s *FileService) AdminPage(ctx context.Context, q AdminFileQuery) ([]model.PrintFileModel, int64, error) {
	return s.files.Paginate(ctx, repository.FileFilter{
		MemberID:      q.MemberID,
		FileName:      q.FileName,
		ConvertStatus: q.ConvertStatus,
	}, q.PageNo, q.PageSize)
}

// CleanupExpired 清理过期打印文件：先删 infra 物理文件，再软删 print_file 记录。
// 参数 limit 为单批上限；返回清理成功的条数，单条失败只记日志不中断。
func (s *FileService) CleanupExpired(ctx context.Context, limit int) (int, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}

	list, err := s.files.ListExpired(ctx, limit)
	if err != nil {
		return 0, err
	}

	cleaned := 0
	for i := range list {
		record := &list[i]

		// 源文件与转换后的 PDF 是两条 infra_file 记录，都要删；相同则只删一次
		if err := s.infra.Delete(ctx, record.InfraFileID); err != nil {
			facades.Log().Errorf("清理打印文件源文件失败 id=%d infraFileID=%d err=%v", record.ID, record.InfraFileID, err)
		}
		if record.PDFFileID > 0 && record.PDFFileID != record.InfraFileID {
			if err := s.infra.Delete(ctx, record.PDFFileID); err != nil {
				facades.Log().Errorf("清理打印文件 PDF 失败 id=%d pdfFileID=%d err=%v", record.ID, record.PDFFileID, err)
			}
		}

		if err := s.files.SoftDelete(ctx, record.ID); err != nil {
			facades.Log().Errorf("软删除打印文件失败 id=%d err=%v", record.ID, err)
			continue
		}
		cleaned++
	}

	return cleaned, nil
}

// Detail 查询文件，仅允许访问本人上传的文件。
func (s *FileService) Detail(ctx context.Context, memberID, id types.ID) (*model.PrintFileModel, error) {
	record, err := s.files.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil || record.MemberUserID != memberID {
		return nil, baseErrs.New(printErrs.FileNotExists)
	}

	return record, nil
}

// markReady 把文件置为可打印并记录页数与 PDF 指针。
func (s *FileService) markReady(ctx context.Context, record *model.PrintFileModel, pdfFileID types.ID, pages int, _ string) error {
	record.PDFFileID = pdfFileID
	record.PageCount = pages
	record.ConvertStatus = printEnums.ConvertStatusReady
	record.ConvertError = ""

	return s.files.Update(ctx, record)
}

// markFailed 把文件置为转换失败并记录原因，返回 nil 表示状态已落库。
func (s *FileService) markFailed(ctx context.Context, record *model.PrintFileModel, reason string) error {
	record.ConvertStatus = printEnums.ConvertStatusFailed
	record.ConvertError = truncateText(reason, 255)

	// 失败原因已写进记录，不再向调用方抛错，避免上传接口因转换失败而整体报错
	_ = s.files.Update(ctx, record)

	return nil
}

// countPDFPages 解析 PDF 页数，解析失败返回错误（页数错会导致计费错）。
func countPDFPages(content []byte) (int, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return 0, err
	}

	pages := reader.NumPage()
	if pages <= 0 {
		return 0, fmt.Errorf("PDF 未解析到任何页面")
	}

	return pages, nil
}

// convertDocxToPDF 调用 LibreOffice 把 DOCX 转成 PDF，返回临时文件路径。
// 服务器未安装 LibreOffice 时返回可读的错误，由前端提示联系管理员。
func convertDocxToPDF(content []byte, fileName string) (string, error) {
	soffice := facades.Config().GetString(sofficeConfigKey, "soffice")
	if _, err := exec.LookPath(soffice); err != nil {
		return "", fmt.Errorf("服务器未安装 LibreOffice，暂不支持 DOCX，请先转成 PDF 再上传")
	}

	workDir, err := os.MkdirTemp("", "print-convert-*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败")
	}

	srcPath := filepath.Join(workDir, "source.docx")
	if err := os.WriteFile(srcPath, content, 0o600); err != nil {
		os.RemoveAll(workDir)

		return "", fmt.Errorf("写入临时文件失败")
	}

	ctx, cancel := context.WithTimeout(context.Background(), convertTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, soffice,
		"--headless",
		"--norestore",
		"--convert-to", "pdf",
		"--outdir", workDir,
		srcPath,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(workDir)

		return "", fmt.Errorf("DOCX 转换失败：%s", truncateText(string(output), 120))
	}

	pdfPath := filepath.Join(workDir, "source.pdf")
	if _, err := os.Stat(pdfPath); err != nil {
		os.RemoveAll(workDir)

		return "", fmt.Errorf("DOCX 转换后未生成 PDF")
	}

	_ = fileName

	return pdfPath, nil
}

// writeTempPDF 把 PDF 字节写到系统临时目录，返回路径，调用方负责删除。
func writeTempPDF(content []byte) (string, error) {
	file, err := os.CreateTemp("", "print-*.pdf")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败")
	}

	if _, err := file.Write(content); err != nil {
		file.Close()
		os.Remove(file.Name())

		return "", fmt.Errorf("写入临时文件失败")
	}
	file.Close()

	return file.Name(), nil
}

// pdfNameOf 把原文件名换成 .pdf 扩展名，保留原名便于门店按文件名对单。
func pdfNameOf(fileName string) string {
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	if base == "" {
		base = "print"
	}

	return base + ".pdf"
}

// truncateText 按字符截断文本，避免超出字段长度。
func truncateText(value string, limit int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= limit {
		return string(runes)
	}

	return string(runes[:limit])
}
