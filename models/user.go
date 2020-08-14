package models

import "time"

type User struct {
	ID     int64
	UserID int64 `gorm:"index:index_user_id"`
	Name   string
	Mobile string
	Status int `gorm:"default:1"`

	Roles []UserRole

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "auth_users"
}
