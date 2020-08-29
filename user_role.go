package auth

import (
	"github.com/gitxiaobo/auth/errors"
	"github.com/gitxiaobo/auth/models"
)

// 用户绑定角色
func (e *Enforcer) CreateOrUpdateUserRole(userID int64, roleIDs []int64) (userRoles []models.UserRole, err error) {
	user, err := e.findUserByUserID(userID)
	if err != nil {
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

	e.DB.Model(&user).Update("AuthStatus", 2)
	return
}
