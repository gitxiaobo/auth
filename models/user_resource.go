package models

import "time"

type UserResource struct {
	ID          int64
	UserID      int64 `gorm:"index:index_user_id"`
	ResourceKey string
	Resoures    string

	CreatedAt time.Time
	UpdatedAt time.Time
}