package auth

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
func (e *Enforcer) DeleteUser(userID int64) (err error) {
	user, err := e.findUserByUserID(userID)
	if err != nil {
		return
	}

	e.DB.Delete(models.UserRole{}, "user_id = ?", user.ID)
	e.DB.Delete(&user)
	return
}

// 获取用户角色
func (e *Enforcer) GetUserRoles(userID int64) (roles []models.Role, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).Preload("Roles").Preload("Roles.Role").First(&user).Error
	if err != nil {
		return
	}
	for _, ur := range user.Roles {
		if ur.Role.Status != 1 {
			continue
		}
		roles = append(roles, ur.Role)
	}
	return
}

// API权限验证
func (e *Enforcer) CheckApiAuth(userID int64, url string, method string) (b bool, err error) {
	if e.Disabled {
		return true, nil
	}

	user, err := e.findUserByUserID(userID)
	if err != nil {
		return
	}

	if b := e.isSuperAdmin(user.ID); b {
		return true, nil
	}

	apiAuthCode, err := e.getApiAuthCode(url, method)
	if len(apiAuthCode) == 0 {
		return true, nil
	}

	userAuthCodes, err := e.getUserApiAuths(user.ID)
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
		if userRole.Role.Status != 1 {
			continue
		}
		for _, roleAuth := range userRole.Role.Auths {
			var codes []string
			json.Unmarshal([]byte(roleAuth.FuncAuthCodes), &codes)
			authCodes = append(authCodes, codes...)
		}
	}
	authCodes = removeDuplicateElement(authCodes)
	return
}

// 获取用户API权限码
func (e *Enforcer) getUserApiAuths(uid int64) (authCodes []string, err error) {
	var userRoles []models.UserRole

	err = e.DB.Where("user_id = ?", uid).Preload("Role").Preload("Role.Auths").Find(&userRoles).Error
	if err != nil {
		return
	}

	for _, userRole := range userRoles {
		if userRole.Role.Status != 1 {
			continue
		}
		for _, roleAuth := range userRole.Role.Auths {
			var apiAuths []string
			json.Unmarshal([]byte(roleAuth.ApiAuthCodes), &apiAuths)
			authCodes = append(authCodes, apiAuths...)
		}
	}
	authCodes = removeDuplicateElement(authCodes)
	return
}

// 超级管理员判断
func (e *Enforcer) isSuperAdmin(uid int64) (b bool) {
	var ur models.UserRole
	err := e.DB.Table("auth_roles").Joins("RIGHT JOIN auth_user_roles on auth_user_roles.role_id = auth_roles.id and auth_user_roles.user_id = ?", uid).Where("auth_roles.name = ?", "超级管理员").First(&ur).Error
	return err == nil
}

func (e *Enforcer) findUserByUserID(userID int64) (user models.User, err error) {
	err = e.DB.Where("user_id = ?", userID).First(&user).Error
	return
}
