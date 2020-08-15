package models

import "time"

var CategoryMap = map[int]string{
	1: "网页端",
}

type RoleAuthority struct {
	ID            int64
	RoleID        int64  `gorm:"index:index_role_id"`
	Category      int    `gorm:"default:1"`
	FuncAuthCodes string `gorm:"func_auth_codes"`
	ApiAuthCodes  string

	CreatedAt time.Time
	UpdatedAt time.Time
}
