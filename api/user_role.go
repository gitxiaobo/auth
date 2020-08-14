package api

import (
	"github.com/gitxiaobo/auth/errors"
	"github.com/gitxiaobo/auth/models"
)

// 用户绑定角色
func (e *Enforcer) CreateOrUpdateUserRole(userID int64, roleIDs []int64) (userRoles []models.UserRole, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).First(&user).Error

	if err != nil {
		err = errors.USER_NOT_FOUND
		return
	}

	for _, roleID := range roleIDs {
		var role models.Role
		err = e.DB.Where("id = ?", roleID).First(&role).Error
		if err != nil {
			err = errors.ROLE_NOT_FOUND
			return
		}
	}

	err = e.DB.Where("user_id = ?", user.ID).Delete(&models.UserRole{}).Error
	if err != nil {
		return
	}

	for _, roleID := range roleIDs {
		var ur models.UserRole
		err = e.DB.FirstOrCreate(&ur, models.UserRole{UserID: user.ID, RoleID: roleID}).Error
	}
	return
}
