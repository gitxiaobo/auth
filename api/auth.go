package api

import (
	"auth/models"
	"encoding/json"

	"github.com/wxnacy/wgo/arrays"
)

// 创建权限相关表
func (e *Enforcer) CreateTable() {
	e.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}, &models.RoleAuthority{})
}

// func DataInit(db *gorm.DB) {
// 	db.Create(&models.User{Name: "伍小波"})
// 	db.Create(&models.Role{Name: "技术人员"})
// 	db.Create(&models.UserRole{UserID: 1, RoleID: 1})

// 	codes := []string{"123","456"}
// 	valueString, _ := json.Marshal(codes)

// 	db.Create(&models.RoleAuthority{RoleID: 1, ApiAuthCodes: string(valueString)})
// }

// API权限验证
func (e *Enforcer) CheckApiAuth(userID int64, apiAuthCode string) (b bool, err error) {
	userAuthCodes, err := e.GetUserApiAuths(userID)
	if err != nil {
		return
	}

	index := arrays.ContainsString(userAuthCodes, apiAuthCode)
	if index != -1 {
		return true, nil
	}

	return
}

// 获取用户API权限码
func (e *Enforcer) GetUserApiAuths(userID int64) (authCodes []string, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).Preload("Roles").Preload("Roles.Role").Preload("Roles.Role.Auths").First(&user).Error
	if err != nil {
		return
	}

	for _, userRole := range user.Roles {
		for _, roleAuth := range userRole.Role.Auths {
			var apiAuths []string
			json.Unmarshal([]byte(roleAuth.ApiAuthCodes), &apiAuths)
			authCodes = append(authCodes, apiAuths...)
		}
	}
	return
}
