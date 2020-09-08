package auth

import (
	"github.com/gitxiaobo/auth/models"

	"github.com/jinzhu/gorm"
)

type Enforcer struct {
	DB                 *gorm.DB
	Disabled           bool
	FuncAuthConfigPath string
	ApiAuthConfigPath  string
	ResourceConfigPath string
}

// 获取实例
func NewEnforcer(db *gorm.DB, funcAuthConfigPath string, apiAuthConfigPath string, resourceConfigPath string) (*Enforcer, error) {
	e := &Enforcer{}

	e.DB = db
	e.FuncAuthConfigPath = funcAuthConfigPath
	e.ApiAuthConfigPath = apiAuthConfigPath
	e.ResourceConfigPath = resourceConfigPath

	if !db.HasTable(&models.User{}) {
		err := e.CreateTable()
		if err != nil {
			return e, err
		}
	}

	return e, nil
}

// 创建权限相关表
func (e *Enforcer) CreateTable() (err error) {
	err = e.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}, &models.RoleAuthority{}, &models.UserResource{}, &models.Resource{}).Error
	return
}
