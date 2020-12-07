package models

import "time"

type Resource struct {
	ID            int64  `json:id`
	ResourceKey   string `json:"key"`
	ResourceName  string `json:"name"`
	ResourceValue string `json:"value"`
	Remark        string `json:"remark"`
	DealerID      int64  `gorm:"index:index_dealer_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Resource) TableName() string {
	return "auth_resources"
}
