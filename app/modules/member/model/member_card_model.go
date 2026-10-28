package model

import (
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// MemberCardTable 是会员卡表名常量，Model 与 Repository 共用同一来源。
const MemberCardTable = "member_card"

// 会员卡状态：正常 / 停用 / 已过期，过期由定时任务写入。
const (
	MemberCardStatusDisabled int8 = 0
	MemberCardStatusNormal   int8 = 1
	MemberCardStatusExpired  int8 = 2
)

// 开卡来源：注册赠送 / 后台发放 / 购买（预留）。
const (
	MemberCardSourceRegister = "register"
	MemberCardSourceAdmin    = "admin"
	MemberCardSourcePurchase = "purchase"
)

// MemberCardModel 对应 member_card 表。
// 会员卡是独立的权益载体：卡自带等级与折扣，不引用会员等级表，
// 支付折扣只读本表的 DiscountPercent（100 无折扣 / 95 表示 95 折）。
type MemberCardModel struct {
	baseModel.BaseModel
	MemberUserID types.ID `db:"member_user_id" json:"memberUserId"`
	// CardNo 卡号，服务端生成，对外展示
	CardNo string `db:"card_no" json:"cardNo"`
	// CardLevel 卡等级编码，如 normal / silver / gold / black
	CardLevel string `db:"card_level" json:"cardLevel"`
	// CardName 卡名称，展示用
	CardName string `db:"card_name" json:"cardName"`
	// DiscountPercent 折扣百分比：100 无折扣 / 95 表示 95 折
	DiscountPercent int `db:"discount_percent" json:"discountPercent"`
	// CardType 卡类型：level 折扣卡（预留 stored 储值卡）
	CardType string `db:"card_type" json:"cardType"`
	// Source 开卡来源：register / admin / purchase
	Source string `db:"source" json:"source"`
	// Status 0 停用 / 1 正常 / 2 已过期
	Status    int8             `db:"status" json:"status"`
	StartTime *types.DateTime  `db:"start_time" json:"startTime"`
	ExpireTime *types.DateTime `db:"expire_time" json:"expireTime"`
	Remark    string           `db:"remark" json:"remark"`
	// Enabled 是派生字段，不落库：仅用于前端展示卡是否可用
	Enabled bool `db:"-" json:"enabled"`
}

// TableName 返回表名，供仓储层构造查询时引用。
func (m *MemberCardModel) TableName() string {
	return MemberCardTable
}

// IsUsable 判断卡当前是否可用于折扣计算：状态正常且未过期。
// 无参数；返回 true 表示可享卡折扣，false 表示按无折扣处理。
func (m *MemberCardModel) IsUsable() bool {
	if m == nil || m.Status != MemberCardStatusNormal {
		return false
	}
	if m.ExpireTime == nil {
		return true
	}

	return m.ExpireTime.Time().After(types.Now().Time())
}
