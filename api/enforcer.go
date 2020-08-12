package api

import (
	"github.com/jinzhu/gorm"
)

type Enforcer struct {
	DB *gorm.DB
	AuthConfigPath string
}

func NewEnforcer(db *gorm.DB, authConfigPath string) (*Enforcer, error) {
	e := &Enforcer{}

	e.DB = db
	e.AuthConfigPath = authConfigPath
	return e, nil
}