package models

import "time"

type Role struct {
	ID       int64  `json:"id"`
	DealerID int64  `gorm:"dealer_id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`

	Status int `json:"status" gorm:"default:1"`

	Auths []RoleAuthority `json:"auths"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "auth_roles"
}
