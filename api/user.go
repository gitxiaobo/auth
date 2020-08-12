package api

import "auth/models"

func (e *Enforcer) CreateOrUpdateUser(userID int64) (user models.User, err error) {
	err = e.DB.FirstOrCreate(&user, models.User{UserID: userID}).Error
	return
}

func (e *Enforcer) DeleteUser(userID int64) (user models.User, err error) {
	err = e.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return
	}

	e.DB.Delete(&user)
	return
}
