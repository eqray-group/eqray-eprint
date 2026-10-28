package repository

import (
	"context"

	"github.com/goravel/framework/contracts/database/db"

	"eqray-eprint/app/facades"
	"eqray-eprint/app/modules/base/types"
	printEnums "eqray-eprint/app/modules/print/enums"
	"eqray-eprint/app/modules/print/model"
)

const orderItemTable = model.PrintOrderItemTable

// PrintOrderItemRepository 是 print_order_item 表的仓储（明细兼打印任务）。
type PrintOrderItemRepository struct{}

// NewPrintOrderItemRepository 构造订单明细仓储。
func NewPrintOrderItemRepository() *PrintOrderItemRepository {
	return &PrintOrderItemRepository{}
}

// BatchCreate 在事务内批量插入明细。
func (r *PrintOrderItemRepository) BatchCreate(ctx context.Context, tx db.Tx, items []model.PrintOrderItemModel) error {
	if len(items) == 0 {
		return nil
	}

	rows := make([]map[string]any, 0, len(items))
	for i := range items {
		m := &items[i]
		m.FillCreate(ctx)
		rows = append(rows, map[string]any{
			"id":              m.ID,
			"print_shop_id":   m.PrintShopID,
			"order_id":        m.OrderID,
			"task_no":         m.TaskNo,
			"file_id":         m.FileID,
			"file_name":       m.FileName,
			"file_page_count": m.FilePageCount,
			"paper_type":      m.PaperType,
			"color_type":      m.ColorType,
			"print_mode":      m.PrintMode,
			"orientation":     m.Orientation,
			"page_range":      m.PageRange,
			"page_count":      m.PageCount,
			"copies":          m.Copies,
			"sheet_count":     m.SheetCount,
			"print_config":    m.PrintConfig,
			"amount":          m.Amount,
			"status":          m.Status,
			"device_id":       m.DeviceID,
			"printer_id":      m.PrinterID,
			"creator":         m.Creator,
			"created_at":      m.CreatedAt,
			"updater":         m.Updater,
			"updated_at":      m.UpdatedAt,
		})
	}

	_, err := tx.Table(orderItemTable).Insert(rows)

	return err
}

// ListByOrder 查询订单明细。
func (r *PrintOrderItemRepository) ListByOrder(ctx context.Context, orderID types.ID) ([]model.PrintOrderItemModel, error) {
	list := make([]model.PrintOrderItemModel, 0)

	if err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("order_id", orderID).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// ListPendingByDevice 查询设备待打印任务：已分配且未开始的明细。
// 客户端拉起任务后需回执「打印中」，避免重复拉取时重复打印。
func (r *PrintOrderItemRepository) ListPendingByDevice(ctx context.Context, deviceID types.ID, limit int) ([]model.PrintOrderItemModel, error) {
	list := make([]model.PrintOrderItemModel, 0)

	if err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("device_id", deviceID).
		Where("status", printEnums.TaskStatusAssigned).
		WhereNull("deleted_at").
		OrderBy("id", "asc").
		Limit(uint64(limit)).
		Get(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// FindByTaskNo 按任务号查询，回执时定位明细。
func (r *PrintOrderItemRepository) FindByTaskNo(ctx context.Context, taskNo string) (*model.PrintOrderItemModel, error) {
	var m model.PrintOrderItemModel

	if err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("task_no", taskNo).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// MarkPrinting 把任务置为打印中并记录开始时间。
func (r *PrintOrderItemRepository) MarkPrinting(ctx context.Context, id types.ID) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("id", id).
		Where("status", printEnums.TaskStatusAssigned).
		Update(map[string]any{
			"status":     printEnums.TaskStatusPrinting,
			"started_at": now,
			"updated_at": now,
		})

	return err
}

// MarkDone 把任务置为完成并记录实际打印页数。
func (r *PrintOrderItemRepository) MarkDone(ctx context.Context, id types.ID, printedPages int) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("id", id).
		Update(map[string]any{
			"status":        printEnums.TaskStatusDone,
			"printed_pages": printedPages,
			"finished_at":   now,
			"updated_at":    now,
		})

	return err
}

// MarkFailed 记录失败信息并累加重试次数，返回累加后的重试次数需要调用方查询。
func (r *PrintOrderItemRepository) MarkFailed(ctx context.Context, id types.ID, errorCode, reason string, retryCount int) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("id", id).
		Update(map[string]any{
			"status":      printEnums.TaskStatusFailed,
			"error_code":  errorCode,
			"fail_reason": reason,
			"retry_count": retryCount,
			"finished_at": now,
			"updated_at":  now,
		})

	return err
}

// Retry 把失败任务重置为已分配，重试次数由调用方控制。
func (r *PrintOrderItemRepository) Retry(ctx context.Context, id types.ID, retryCount int) error {
	_, err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("id", id).
		Update(map[string]any{
			"status":      printEnums.TaskStatusAssigned,
			"retry_count": retryCount,
			"error_code":  "",
			"fail_reason": "",
			"updated_at":  types.Now(),
		})

	return err
}

// CancelByOrder 把订单下未开始的任务置为已取消，订单超时关闭时调用。
// 已打印或已完成的任务保持原状，避免吞掉真实打印记录。
func (r *PrintOrderItemRepository) CancelByOrder(ctx context.Context, orderID types.ID) error {
	_, err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("order_id", orderID).
		Where("status", printEnums.TaskStatusAssigned).
		Update(map[string]any{
			"status":     printEnums.TaskStatusCanceled,
			"updated_at": types.Now(),
		})

	return err
}

// CountByStatus 统计订单下某状态的任务数量。
func (r *PrintOrderItemRepository) CountByStatus(ctx context.Context, orderID types.ID, status int8) (int64, error) {
	total, err := facades.DB().WithContext(ctx).Table(orderItemTable).
		Where("order_id", orderID).
		Where("status", status).
		WhereNull("deleted_at").
		Count()
	if err != nil {
		return 0, err
	}

	return total, nil
}
