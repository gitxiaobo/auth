package models

import "time"

type User struct {
	ID         int64
	UserID     int64  `gorm:"index:index_user_id"`
	DealerID   int64  `gorm:"index:index_dealer_id"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	Status     int    `gorm:"default:1"`
	AuthStatus int    `gorm:"default:1"` //权限是否发生变化 1 - 否， 2 - 是

	Roles []UserRole

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "auth_users"
}
