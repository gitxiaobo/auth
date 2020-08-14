package models

import "time"

type Role struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status int    `gorm:"default:1"`

	Auths []RoleAuthority

	CreatedAt time.Time
	UpdatedAt time.Time
}
