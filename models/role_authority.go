package models

import "time"

var CategoryMap = map[int]string{
	1: "网页端",
}

type RoleAuthority struct {
	ID            int64
	RoleID        int64
	Category      int `gorm:"default:1"`
	FuncAuthCodes string
	ApiAuthCodes  string

	CreatedAt time.Time
	UpdatedAt time.Time
}
