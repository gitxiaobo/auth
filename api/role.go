package api

import (
	"auth/models"
	"encoding/json"
)

func (e *Enforcer) CreateOrUpdateRole(roleID int64, name string) (role models.Role, err error) {
	role.ID = roleID
	role.Name = name

	err = e.DB.Where("id = ?", role.ID).First(&role).Error
	if err != nil {
		e.DB.Create(&role)
	}

	e.DB.Model(&role).Update("Name", name)
	return
}

func (e *Enforcer) Deleterole(roleID int64) (role models.Role, err error) {
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

	e.DB.Delete(&role)
	return
}

func (e *Enforcer) GetRoles() (roles []models.Role, err error) {
	err = e.DB.Find(&roles).Error
	return
}

func (e *Enforcer) CreateOrUpdateRoleAuths(roleID int64, authCodes []string) (err error) {
	var role models.Role
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

	authCodesString, _ := json.Marshal(authCodes)

	var ra models.RoleAuthority
	ra.RoleID = roleID
	ra.FuncAuthCodes = string(authCodesString)

	err = e.DB.Where("role_id = ?", ra.RoleID).First(&ra).Error
	if err != nil {
		err = e.DB.Create(&ra).Error
		return
	}

	err = e.DB.Model(&ra).Updates(models.RoleAuthority{FuncAuthCodes: string(authCodesString)}).Error
	return
}
