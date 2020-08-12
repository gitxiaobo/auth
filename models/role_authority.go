package models

import "time"

type RoleAuthority struct {
	ID            int64
	RoleID        int64
	FuncAuthCodes string
	ApiAuthCodes  string

	CreatedAt time.Time
	UpdatedAt time.Time
}
