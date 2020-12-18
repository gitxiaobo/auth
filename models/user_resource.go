package models

import "time"

type UserResource struct {
	ID            int64
	UserID        int64  `gorm:"index:index_user_id"`
	ResourceKey   string `json:"key"`
	ResourceValue string `gorm:"type:text" json:"value"`
	FieldName     string `json:"field_name"`
	DealerID      int64  `gorm:"index:index_dealer_id"`
	AreaID        int    `gorm:"index:index_ared_id;default:0" json:"area_id"` //区域
	AllArea       int    `json:"all_area" gorm:"default:0"`                    //1 - 区域下所有资源， 0 - 选中人员

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserResource) TableName() string {
	return "auth_user_resources"
}
