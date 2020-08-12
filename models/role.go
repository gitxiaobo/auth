package models

import "time"

type Role struct {
	ID     int64
	Name   string
	Status int `gorm:"default:1"`

	Auths []RoleAuthority

	CreatedAt time.Time
	UpdatedAt time.Time
}
