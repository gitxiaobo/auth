package api

import (
	"encoding/json"

	"github.com/gitxiaobo/auth/models"

	"github.com/wxnacy/wgo/arrays"
)

// 创建或更新用户
func (e *Enforcer) CreateOrUpdateUser(userID int64) (user models.User, err error) {
	err = e.DB.FirstOrCreate(&user, models.User{UserID: userID}).Error
	return
}

// 删除用户
func (e *Enforcer) DeleteUser(userID int64) (user models.User, err error) {
	err = e.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return
	}

	e.DB.Delete(models.UserRole{}, "user_id = ?", user.ID)
	e.DB.Delete(&user)
	return
}

// API权限验证
func (e *Enforcer) CheckApiAuth(userID int64, url string, method string) (b bool, err error) {
	if e.Disabled {
		return true, nil
	}

	apiAuthCode, err := e.GetApiAuthCode(url, method)
	if len(apiAuthCode) == 0 {
		return true, nil
	}

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

// 获取用户前端权限码
func (e *Enforcer) GetUserFuncAuths(userID int64) (authCodes []string, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).Preload("Roles").Preload("Roles.Role").Preload("Roles.Role.Auths").First(&user).Error
	if err != nil {
		return
	}

	for _, userRole := range user.Roles {
		for _, roleAuth := range userRole.Role.Auths {
			var codes []string
			json.Unmarshal([]byte(roleAuth.FuncAuthCodes), &codes)
			authCodes = append(authCodes, codes...)
		}
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

func (e *Enforcer) GetUserRoles(userID int64) (roles []models.Role, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).Preload("Roles").Preload("Roles.Role").First(&user).Error
	if err != nil {
		return
	}
	for _, ur := range user.Roles {
		roles = append(roles, ur.Role)
	}
	return
}
