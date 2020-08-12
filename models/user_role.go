package models

import "time"

type UserRole struct {
	ID     int64
	UserID int64
	RoleID int64

	Role Role

	CreatedAt time.Time
	UpdatedAt time.Time
}
