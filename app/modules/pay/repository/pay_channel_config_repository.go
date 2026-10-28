package repository

import (
	"context"

	"eqray-eprint/app/facades"
	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/support"
	"eqray-eprint/app/modules/base/types"
	payErrs "eqray-eprint/app/modules/pay/errs"
	"eqray-eprint/app/modules/pay/model"
)

const channelConfigTable = model.PayChannelConfigTable

// ChannelConfigFilter 是渠道配置列表的筛选条件，零值表示不限制。
type ChannelConfigFilter struct {
	Code    string
	Channel string
	Status  *int8
}

// PayChannelConfigRepository 是 pay_channel_config 表的仓储。
type PayChannelConfigRepository struct{}

// NewPayChannelConfigRepository 构造渠道配置仓储。
func NewPayChannelConfigRepository() *PayChannelConfigRepository {
	return &PayChannelConfigRepository{}
}

// Create 新增渠道配置；code 或 (channel, scene) 冲突转换为业务错误码。
func (r *PayChannelConfigRepository) Create(ctx context.Context, m *model.PayChannelConfigModel) error {
	m.FillCreate(ctx)

	_, err := facades.DB().WithContext(ctx).Table(channelConfigTable).Insert(r.values(m))
	if err != nil && support.IsDuplicate(err) {
		return baseErrs.New(payErrs.PayChannelConfigCodeExists, m.Code)
	}

	return err
}

// Update 全量更新渠道配置（含密钥字段），调用方负责保留未修改的密钥。
func (r *PayChannelConfigRepository) Update(ctx context.Context, m *model.PayChannelConfigModel) error {
	m.FillUpdate(ctx)

	values := r.values(m)
	delete(values, "id")
	delete(values, "creator")
	delete(values, "created_at")

	_, err := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(values)
	if err != nil && support.IsDuplicate(err) {
		return baseErrs.New(payErrs.PayChannelConfigCodeExists, m.Code)
	}

	return err
}

// UpdateStatus 只更新启用状态，列表切换开关时不动密钥。
func (r *PayChannelConfigRepository) UpdateStatus(ctx context.Context, id types.ID, status int8, updater string) error {
	now := types.Now()

	_, err := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("id", id).
		WhereNull("deleted_at").
		Update(map[string]any{
			"status":     status,
			"updater":    updater,
			"updated_at": now,
		})

	return err
}

// Delete 软删除渠道配置；历史支付单仍按 channel_code 关联，不做物理删除。
func (r *PayChannelConfigRepository) Delete(ctx context.Context, m *model.PayChannelConfigModel) error {
	m.FillUpdate(ctx)

	res, err := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("id", m.ID).
		WhereNull("deleted_at").
		Update(map[string]any{
			"deleted_at": types.Now(),
			"updater":    m.Updater,
			"updated_at": m.UpdatedAt,
		})
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return baseErrs.New(payErrs.PayChannelConfigNotExists, m.ID)
	}

	return nil
}

// FindByID 按主键查询未删除的配置；未命中返回 (nil, nil)。
func (r *PayChannelConfigRepository) FindByID(ctx context.Context, id types.ID) (*model.PayChannelConfigModel, error) {
	var m model.PayChannelConfigModel

	if err := facades.DB().WithContext(ctx).Table(channelConfigTable).
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

// FindByCode 按配置编码查询未删除的记录；未命中返回 (nil, nil)。
func (r *PayChannelConfigRepository) FindByCode(ctx context.Context, code string) (*model.PayChannelConfigModel, error) {
	var m model.PayChannelConfigModel

	if err := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("code", code).
		WhereNull("deleted_at").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// FindFirstEnabledByChannel 取该渠道下任一启用配置。
// 平台统一收款，同一渠道的多条场景配置共用商户号与平台公钥，
// 回调验签只需要其中一条即可，避免按 appid 反推配置。
func (r *PayChannelConfigRepository) FindFirstEnabledByChannel(ctx context.Context, channel string) (*model.PayChannelConfigModel, error) {
	var m model.PayChannelConfigModel

	if err := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("channel", channel).
		Where("status", 1).
		WhereNull("deleted_at").
		OrderBy("sort", "asc").
		First(&m); err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}

	return &m, nil
}

// ExistsByCode 判断编码是否已被占用，excludeID 用于编辑时排除自身。
func (r *PayChannelConfigRepository) ExistsByCode(ctx context.Context, code string, excludeID types.ID) (bool, error) {
	query := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("code", code).
		WhereNull("deleted_at")
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	total, err := query.Count()
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// ExistsByScene 判断 (channel, scene) 是否已配置，保证同一场景只有一份商户配置。
func (r *PayChannelConfigRepository) ExistsByScene(ctx context.Context, channel, scene string, excludeID types.ID) (bool, error) {
	query := facades.DB().WithContext(ctx).Table(channelConfigTable).
		Where("channel", channel).
		Where("scene", scene).
		WhereNull("deleted_at")
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	total, err := query.Count()
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// Paginate 按条件分页查询渠道配置。
func (r *PayChannelConfigRepository) Paginate(ctx context.Context, filter ChannelConfigFilter, pageNo, pageSize int) ([]model.PayChannelConfigModel, int64, error) {
	list := make([]model.PayChannelConfigModel, 0)
	var total int64

	query := facades.DB().WithContext(ctx).Table(channelConfigTable).WhereNull("deleted_at")
	if filter.Code != "" {
		query = query.WhereLike("code", "%"+filter.Code+"%")
	}
	if filter.Channel != "" {
		query = query.Where("channel", filter.Channel)
	}
	if filter.Status != nil {
		query = query.Where("status", *filter.Status)
	}

	if err := query.OrderBy("sort", "asc").OrderByDesc("id").Paginate(pageNo, pageSize, &list, &total); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// values 把模型转成写入用的字段集合；extra 为空时写 NULL，避免 JSON 列拒绝空串。
func (r *PayChannelConfigRepository) values(m *model.PayChannelConfigModel) map[string]any {
	var extra any
	if m.Extra != "" {
		extra = m.Extra
	}

	return map[string]any{
		"id":                  m.ID,
		"code":                m.Code,
		"name":                m.Name,
		"channel":             m.Channel,
		"scene":               m.Scene,
		"app_id":              m.AppID,
		"mch_id":              m.MchID,
		"serial_no":           m.SerialNo,
		"private_key":         m.PrivateKey,
		"api_key":             m.APIKey,
		"platform_public_key": m.PlatformPublicKey,
		"gateway":             m.Gateway,
		"notify_url":          m.NotifyURL,
		"return_url":          m.ReturnURL,
		"extra":               extra,
		"status":              m.Status,
		"sort":                m.Sort,
		"remark":              m.Remark,
		"creator":             m.Creator,
		"created_at":          m.CreatedAt,
		"updater":             m.Updater,
		"updated_at":          m.UpdatedAt,
	}
}
