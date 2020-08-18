package models

import "time"

var CategoryMap = map[int]string{
	1: "网页端",
}

type RoleAuthority struct {
	ID            int64
	RoleID        int64  `gorm:"index:index_role_id"`
	Category      int    `gorm:"default:1"`
	FuncAuthCodes string `json:"func_auth_codes"`
	ChosedCodes   string `json:"chosed_codes"` //前端配置角色权限时显示的权限码，和返回的用户权限码都差别
	ApiAuthCodes  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RoleAuthority) TableName() string {
	return "auth_role_authorities"
}
