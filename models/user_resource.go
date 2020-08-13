package models

import "time"

type UserResource struct {
	ID          int64
	UserID      int64
	ResourceKey string
	Resoures    string

	CreatedAt time.Time
	UpdatedAt time.Time
}
