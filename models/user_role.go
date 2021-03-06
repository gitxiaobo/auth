package models

import "time"

type UserRole struct {
	ID     int64
	User   User  `json:"-"`
	UserID int64 `gorm:"index:index_user_id"`
	RoleID int64 `gorm:"index:index_role_id"`

	Role Role

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserRole) TableName() string {
	return "auth_user_roles"
}
