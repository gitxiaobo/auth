package models

import "time"

type UserResource struct {
	ID            int64
	UserID        int64  `gorm:"index:index_user_id"`
	ResourceKey   string `json:"key"`
	ResourceValue string `json:"value"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserResource) TableName() string {
	return "auth_user_resources"
}
