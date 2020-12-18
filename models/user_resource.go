package models

import "time"

type UserResource struct {
	ID            int64
	UserID        int64  `gorm:"index:index_user_id"`
	ResourceKey   string `json:"key"`
	ResourceValue string `gorm:"type:text" json:"value"`
	FieldName     string `json:"field_name"`
	DealerID      int64  `gorm:"index:index_dealer_id"`
	AreaID        int    `gorm:"index:index_ared_id"` //区域
	All           int    `json:"all"`                 //1 - 区域下所有资源， 0 - 选中人员

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserResource) TableName() string {
	return "auth_user_resources"
}
