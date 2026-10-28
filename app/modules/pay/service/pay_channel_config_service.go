package service

import (
	"context"
	"encoding/json"
	"strings"

	baseErrs "eqray-eprint/app/modules/base/errs"
	"eqray-eprint/app/modules/base/types"
	payEnums "eqray-eprint/app/modules/pay/enums"
	payErrs "eqray-eprint/app/modules/pay/errs"
	"eqray-eprint/app/modules/pay/gateway"
	"eqray-eprint/app/modules/pay/model"
	"eqray-eprint/app/modules/pay/repository"
)

// ChannelConfigInput 是后台新增 / 编辑渠道配置的入参。
// 密钥字段为空表示「不修改」，编辑时前端不回填也可以安全提交。
type ChannelConfigInput struct {
	Code              string
	Name              string
	Channel           string
	Scene             string
	AppID             string
	MchID             string
	SerialNo          string
	PrivateKey        string
	APIKey            string
	PlatformPublicKey string
	Gateway           string
	NotifyURL         string
	ReturnURL         string
	Extra             string
	Status            int8
	Sort              int
	Remark            string
}

// ChannelConfigListQuery 是渠道配置列表的查询条件。
type ChannelConfigListQuery struct {
	Code     string
	Channel  string
	Status   *int8
	PageNo   int
	PageSize int
}

// ChannelConfigDetail 是后台编辑页的回显结构。
// 后台需要看到并修改密钥，因此这里返回明文；列表接口不返回该结构。
type ChannelConfigDetail struct {
	*model.PayChannelConfigModel
	PrivateKey        string `json:"privateKey"`
	APIKey            string `json:"apiKey"`
	PlatformPublicKey string `json:"platformPublicKey"`
}

// ChannelConfigService 负责支付渠道配置的维护与运行期读取。
type ChannelConfigService struct {
	configs *repository.PayChannelConfigRepository
}

// NewChannelConfigService 构造渠道配置服务。
func NewChannelConfigService() *ChannelConfigService {
	return &ChannelConfigService{configs: repository.NewPayChannelConfigRepository()}
}

// Config 按配置编码取运行期配置，供 Gateway 使用。
// 配置不存在或已停用时返回业务错误码，调用方不得降级到其它渠道。
func (s *ChannelConfigService) Config(ctx context.Context, code string) (*gateway.Config, error) {
	m, err := s.configs.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, baseErrs.New(payErrs.PayChannelConfigNotExists, code)
	}
	if m.Status != payEnums.ChannelConfigEnabled {
		return nil, baseErrs.New(payErrs.PayChannelConfigDisabled, code)
	}
	if err := validateChannelParams(m); err != nil {
		return nil, err
	}

	return &gateway.Config{
		Code:              m.Code,
		Channel:           m.Channel,
		Scene:             m.Scene,
		AppID:             m.AppID,
		MchID:             m.MchID,
		SerialNo:          m.SerialNo,
		PrivateKey:        m.PrivateKey,
		APIKey:            m.APIKey,
		PlatformPublicKey: m.PlatformPublicKey,
		GatewayURL:        m.Gateway,
		NotifyURL:         m.NotifyURL,
		ReturnURL:         m.ReturnURL,
	}, nil
}

// AnyEnabled 取该渠道下任一启用配置，供回调验签使用。
// 平台统一收款：同一渠道的多条场景配置共用商户号与平台公钥，
// 因此回调不需要知道具体是哪个 appid 下的单。
func (s *ChannelConfigService) AnyEnabled(ctx context.Context, channel string) (*gateway.Config, error) {
	m, err := s.configs.FindFirstEnabledByChannel(ctx, channel)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, baseErrs.New(payErrs.PayChannelConfigNotExists, channel)
	}

	return &gateway.Config{
		Code:              m.Code,
		Channel:           m.Channel,
		Scene:             m.Scene,
		AppID:             m.AppID,
		MchID:             m.MchID,
		SerialNo:          m.SerialNo,
		PrivateKey:        m.PrivateKey,
		APIKey:            m.APIKey,
		PlatformPublicKey: m.PlatformPublicKey,
		GatewayURL:        m.Gateway,
		NotifyURL:         m.NotifyURL,
		ReturnURL:         m.ReturnURL,
	}, nil
}

// Page 分页查询渠道配置，列表不返回密钥。
func (s *ChannelConfigService) Page(ctx context.Context, q ChannelConfigListQuery) ([]model.PayChannelConfigModel, int64, error) {
	return s.configs.Paginate(ctx, repository.ChannelConfigFilter{
		Code:    strings.TrimSpace(q.Code),
		Channel: strings.TrimSpace(q.Channel),
		Status:  q.Status,
	}, q.PageNo, q.PageSize)
}

// Detail 查询单条配置并回显密钥，仅用于后台编辑页。
func (s *ChannelConfigService) Detail(ctx context.Context, id types.ID) (*ChannelConfigDetail, error) {
	m, err := s.configs.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, baseErrs.New(payErrs.PayChannelConfigNotExists, id)
	}

	return &ChannelConfigDetail{
		PayChannelConfigModel: m,
		PrivateKey:            m.PrivateKey,
		APIKey:                m.APIKey,
		PlatformPublicKey:     m.PlatformPublicKey,
	}, nil
}

// Create 新增渠道配置，要求 code 与 (channel, scene) 均唯一。
func (s *ChannelConfigService) Create(ctx context.Context, in ChannelConfigInput) (types.ID, error) {
	normalizeChannelInput(&in)
	if err := validateChannelInput(in); err != nil {
		return 0, err
	}

	exists, err := s.configs.ExistsByCode(ctx, in.Code, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, baseErrs.New(payErrs.PayChannelConfigCodeExists, in.Code)
	}

	exists, err = s.configs.ExistsByScene(ctx, in.Channel, in.Scene, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, baseErrs.New(payErrs.PayChannelSceneExists, in.Scene)
	}

	m := &model.PayChannelConfigModel{
		Code:              in.Code,
		Name:              in.Name,
		Channel:           in.Channel,
		Scene:             in.Scene,
		AppID:             in.AppID,
		MchID:             in.MchID,
		SerialNo:          in.SerialNo,
		PrivateKey:        in.PrivateKey,
		APIKey:            in.APIKey,
		PlatformPublicKey: in.PlatformPublicKey,
		Gateway:           in.Gateway,
		NotifyURL:         in.NotifyURL,
		ReturnURL:         in.ReturnURL,
		Extra:             in.Extra,
		Status:            in.Status,
		Sort:              in.Sort,
		Remark:            in.Remark,
	}
	if err := s.configs.Create(ctx, m); err != nil {
		return 0, err
	}

	return m.ID, nil
}

// Update 编辑渠道配置；密钥字段为空时保留原值，避免前端回显不全导致密钥被清空。
func (s *ChannelConfigService) Update(ctx context.Context, id types.ID, in ChannelConfigInput) error {
	normalizeChannelInput(&in)
	if err := validateChannelInput(in); err != nil {
		return err
	}

	m, err := s.configs.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return baseErrs.New(payErrs.PayChannelConfigNotExists, id)
	}

	if in.Code != m.Code {
		exists, err := s.configs.ExistsByCode(ctx, in.Code, id)
		if err != nil {
			return err
		}
		if exists {
			return baseErrs.New(payErrs.PayChannelConfigCodeExists, in.Code)
		}
	}
	if in.Channel != m.Channel || in.Scene != m.Scene {
		exists, err := s.configs.ExistsByScene(ctx, in.Channel, in.Scene, id)
		if err != nil {
			return err
		}
		if exists {
			return baseErrs.New(payErrs.PayChannelSceneExists, in.Scene)
		}
	}

	m.Code = in.Code
	m.Name = in.Name
	m.Channel = in.Channel
	m.Scene = in.Scene
	m.AppID = in.AppID
	m.MchID = in.MchID
	m.SerialNo = in.SerialNo
	m.Gateway = in.Gateway
	m.NotifyURL = in.NotifyURL
	m.ReturnURL = in.ReturnURL
	m.Extra = in.Extra
	m.Status = in.Status
	m.Sort = in.Sort
	m.Remark = in.Remark
	if in.PrivateKey != "" {
		m.PrivateKey = in.PrivateKey
	}
	if in.APIKey != "" {
		m.APIKey = in.APIKey
	}
	if in.PlatformPublicKey != "" {
		m.PlatformPublicKey = in.PlatformPublicKey
	}

	return s.configs.Update(ctx, m)
}

// Delete 软删除渠道配置；已产生的支付单仍保留 channel_code 便于追溯。
func (s *ChannelConfigService) Delete(ctx context.Context, id types.ID) error {
	m, err := s.configs.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return baseErrs.New(payErrs.PayChannelConfigNotExists, id)
	}

	return s.configs.Delete(ctx, m)
}

// validateChannelInput 校验必填项与枚举取值，避免存进无法使用的配置。
func validateChannelInput(in ChannelConfigInput) error {
	if in.Code == "" || in.Name == "" {
		return baseErrs.WithMsg(baseErrs.BadRequest, "配置编码与名称不能为空")
	}
	if in.Channel != payEnums.ChannelWechat && in.Channel != payEnums.ChannelAlipay {
		return baseErrs.WithMsg(baseErrs.BadRequest, "渠道只能是 wechat 或 alipay")
	}
	if !validScene(in.Channel, in.Scene) {
		return baseErrs.WithMsg(baseErrs.BadRequest, "支付场景取值不合法")
	}
	if in.Status != payEnums.ChannelConfigEnabled && in.Status != payEnums.ChannelConfigDisabled {
		return baseErrs.WithMsg(baseErrs.BadRequest, "状态取值不合法")
	}
	// extra 落 JSON 列，非法 JSON 会让 MySQL 直接报错，这里提前拦成业务错误
	if in.Extra != "" && !json.Valid([]byte(in.Extra)) {
		return baseErrs.WithMsg(baseErrs.BadRequest, "扩展参数必须是合法的 JSON")
	}

	return nil
}

// validateChannelParams 在运行期校验渠道参数完整性，缺一项都不允许发起支付。
func validateChannelParams(m *model.PayChannelConfigModel) error {
	missing := make([]string, 0)

	if m.Channel == payEnums.ChannelWechat {
		if m.AppID == "" {
			missing = append(missing, "app_id")
		}
		if m.MchID == "" {
			missing = append(missing, "mch_id")
		}
		if m.SerialNo == "" {
			missing = append(missing, "serial_no")
		}
		if m.PrivateKey == "" {
			missing = append(missing, "private_key")
		}
		if m.APIKey == "" {
			missing = append(missing, "api_key")
		}
		if m.PlatformPublicKey == "" {
			missing = append(missing, "platform_public_key")
		}
	} else {
		if m.AppID == "" {
			missing = append(missing, "app_id")
		}
		if m.PrivateKey == "" {
			missing = append(missing, "private_key")
		}
		if m.PlatformPublicKey == "" {
			missing = append(missing, "platform_public_key")
		}
	}
	if m.NotifyURL == "" {
		missing = append(missing, "notify_url")
	}

	if len(missing) > 0 {
		return baseErrs.New(payErrs.PayChannelConfigInvalid, strings.Join(missing, "、"))
	}

	return nil
}

// validScene 判断场景是否属于对应渠道，防止出现 wechat + wap 这类无效组合。
// 支付宝当前只实现手机网站支付（wap），扫码（qr）待接入后再放开。
func validScene(channel, scene string) bool {
	if channel == payEnums.ChannelWechat {
		return scene == payEnums.SceneMini || scene == payEnums.SceneJsapi || scene == payEnums.SceneH5
	}

	return scene == payEnums.SceneWap
}

// normalizeChannelInput 统一去空格并补齐默认值。
func normalizeChannelInput(in *ChannelConfigInput) {
	in.Code = strings.TrimSpace(in.Code)
	in.Name = strings.TrimSpace(in.Name)
	in.Channel = strings.TrimSpace(in.Channel)
	in.Scene = strings.TrimSpace(in.Scene)
	in.AppID = strings.TrimSpace(in.AppID)
	in.MchID = strings.TrimSpace(in.MchID)
	in.SerialNo = strings.TrimSpace(in.SerialNo)
	in.PrivateKey = strings.TrimSpace(in.PrivateKey)
	in.APIKey = strings.TrimSpace(in.APIKey)
	in.PlatformPublicKey = strings.TrimSpace(in.PlatformPublicKey)
	in.Gateway = strings.TrimSpace(in.Gateway)
	in.NotifyURL = strings.TrimSpace(in.NotifyURL)
	in.ReturnURL = strings.TrimSpace(in.ReturnURL)
	in.Extra = strings.TrimSpace(in.Extra)
	in.Remark = strings.TrimSpace(in.Remark)
	if in.Status != payEnums.ChannelConfigEnabled && in.Status != payEnums.ChannelConfigDisabled {
		in.Status = payEnums.ChannelConfigEnabled
	}
}
