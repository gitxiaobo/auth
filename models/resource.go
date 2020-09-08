package models

import "time"

type Resource struct {
	ID     int64
	Key    string `json:"key"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Remark string `json:"remark"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Resource) TableName() string {
	return "auth_resources"
}
