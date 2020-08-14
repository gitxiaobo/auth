package api

import (
	"encoding/json"

	"github.com/gitxiaobo/auth/models"
)

// 创建或更新角色
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

// 删除角色
func (e *Enforcer) Deleterole(roleID int64) (role models.Role, err error) {
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

	e.DB.Delete(&role)
	return
}

// 获取角色列表
func (e *Enforcer) GetRoles() (roles []models.Role, err error) {
	err = e.DB.Find(&roles).Error
	return
}

// 更新角色权限
func (e *Enforcer) CreateOrUpdateRoleAuths(roleID int64, authCodes []string) (err error) {
	var role models.Role
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

	apiAuthCodes := []string{}

	auths, _ := e.GetAuths()
	for _, funcAuthCode := range authCodes {
		for _, moduleAuth := range auths {
			for _, funcAuth := range moduleAuth.Auths {
				if funcAuthCode != funcAuth.FuncAuthCode {
					continue
				}
				for _, apiAuth := range funcAuth.ApiAuths {
					apiAuthCodes = append(apiAuthCodes, apiAuth.AuthCode)
				}
			}
		}
	}

	authCodesString, _ := json.Marshal(authCodes)
	apiAuthCodesString, _ := json.Marshal(apiAuthCodes)

	var ra models.RoleAuthority
	ra.RoleID = roleID
	ra.FuncAuthCodes = string(authCodesString)
	ra.ApiAuthCodes = string(apiAuthCodesString)

	err = e.DB.Where("role_id = ?", ra.RoleID).First(&ra).Error
	if err != nil {
		err = e.DB.Create(&ra).Error
		return
	}

	err = e.DB.Model(&ra).Updates(models.RoleAuthority{FuncAuthCodes: string(authCodesString), ApiAuthCodes: string(apiAuthCodesString)}).Error
	return
}

// 获取角色前端权限码
func (e *Enforcer) GetRoleFuncAuths(roleID int64) (authCodes []string, err error) {
	var role models.Role
	err = e.DB.Where("id = ?", roleID).Preload("Auths").First(&role).Error
	if err != nil {
		return
	}

	for _, roleAuth := range role.Auths {
		var codes []string
		json.Unmarshal([]byte(roleAuth.FuncAuthCodes), &codes)
		authCodes = append(authCodes, codes...)
	}
	return
}
