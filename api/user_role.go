package api

import (
	"auth/models"
)

func (e *Enforcer) CreateOrUpdateUserRole(userID int64, roleIDs []int64) (userRoles []models.UserRole, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).First(&user).Error

	if err != nil {
		return
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
