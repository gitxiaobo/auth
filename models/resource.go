package models

import "time"

type Resource struct {
	ID            int64
	ResourceKey   string `json:"key"`
	ResourceName  string `json:"name"`
	ResourceValue string `json:"value"`
	Remark        string `json:"remark"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Resource) TableName() string {
	return "auth_resources"
}
