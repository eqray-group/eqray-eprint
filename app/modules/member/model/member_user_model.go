package model

import (
	baseEnums "eqray-eprint/app/modules/base/enums"
	baseModel "eqray-eprint/app/modules/base/model"
	"eqray-eprint/app/modules/base/types"
)

// MemberUserTable 是表名常量，Model 与 Repository 共用同一来源。
const MemberUserTable = "member_user"

// MemberUserModel 对应 member_user 表，只描述 C 端会员身份与基础资料。
// 经验值、积分、标签、分组、会员等级体系已下线：折扣由 member_card 承担，
// 等级仅用于展示，因此本表不再持有等级字段。
type MemberUserModel struct {
	baseModel.BaseModel
	Mobile string `db:"mobile" json:"mobile"`
	// Password 预留：后续开放 C 端密码登录时存 bcrypt 哈希
	Password string `db:"password" json:"-"`
	Nickname string `db:"nickname" json:"nickname"`
	Avatar   string `db:"avatar" json:"avatar"`
	Name     string `db:"name" json:"name"`
	// Sex 0 未知 / 1 男 / 2 女
	Sex      int8            `db:"sex" json:"sex"`
	Birthday *types.DateTime `db:"birthday" json:"birthday"`
	AreaID   int             `db:"area_id" json:"areaId"`
	// Mark 后台标记，如「高价值客户」
	Mark       string                     `db:"mark" json:"mark"`
	RegisterIP string                     `db:"register_ip" json:"registerIp"`
	LoginIP    string                     `db:"login_ip" json:"loginIp"`
	LoginDate  *types.DateTime            `db:"login_date" json:"loginDate"`
	Status     baseEnums.CommonStatusEnum `db:"status" json:"status"`
}

// TableName 返回表名，供仓储层构造查询时引用。
func (m *MemberUserModel) TableName() string {
	return MemberUserTable
}
