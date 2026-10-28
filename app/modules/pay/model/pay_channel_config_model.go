package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
)

// PayChannelConfigTable 是支付渠道配置表名常量。
const PayChannelConfigTable = "pay_channel_config"

// PayChannelConfigModel 对应 pay_channel_config，一行代表一个支付场景的商户配置。
// 平台统一收款，同一 (Channel, Scene) 只有一条启用配置；
// PrivateKey / APIKey / PlatformPublicKey 在库中应为密文，读取后需解密再交给 Gateway。
type PayChannelConfigModel struct {
	baseModel.BaseModel
	Code              string `db:"code" json:"code"`
	Name              string `db:"name" json:"name"`
	Channel           string `db:"channel" json:"channel"`
	Scene             string `db:"scene" json:"scene"`
	AppID             string `db:"app_id" json:"appId"`
	MchID             string `db:"mch_id" json:"mchId"`
	SerialNo          string `db:"serial_no" json:"serialNo"`
	PrivateKey        string `db:"private_key" json:"-"`
	APIKey            string `db:"api_key" json:"-"`
	PlatformPublicKey string `db:"platform_public_key" json:"-"`
	Gateway           string `db:"gateway" json:"gateway"`
	NotifyURL         string `db:"notify_url" json:"notifyUrl"`
	ReturnURL         string `db:"return_url" json:"returnUrl"`
	Extra             string `db:"extra" json:"extra"`
	Status            int8   `db:"status" json:"status"`
	Sort              int    `db:"sort" json:"sort"`
	Remark            string `db:"remark" json:"remark"`
}

// TableName 返回表名。
func (m *PayChannelConfigModel) TableName() string {
	return PayChannelConfigTable
}
