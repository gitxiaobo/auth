package auth

import (
	"encoding/json"

	"github.com/gitxiaobo/auth/errors"

	"github.com/gitxiaobo/auth/models"
)

// 创建或更新角色
func (e *Enforcer) CreateOrUpdateRole(role models.Role, category int, codes []string) (rl models.Role, err error) {
	var r models.Role
	err = e.DB.Where("name = ? and dealer_id = ?", role.Name, role.DealerID).First(&r).Error

	if (role.ID == 0 && err == nil) || (role.ID > 0 && role.ID != r.ID && err == nil) {
		err = errors.ROLE_NAME_REPEAT
		return
	}

	err = e.DB.Where("id = ?", role.ID).First(&r).Error
	if role.ID > 0 && err != nil {
		err = errors.DB_ERROR
		return
	}

	if err != nil {
		err = e.DB.Create(&role).Error
	} else {
		err = e.DB.Omit("status").Save(&role).Error
	}

	if err == nil {
		err = e.createOrUpdateRoleAuths(role.ID, category, codes)
	}

	err = e.DB.Where("id = ?", role.ID).Preload("Auths").First(&rl).Error
	return
}

// 角色状态切换
func (e *Enforcer) SwitchRoleStatus(id int64) (err error) {
	var r models.Role
	err = e.DB.Where("id = ?", id).First(&r).Error
	if err != nil {
		return
	}

	status := 1
	if r.Status == 1 {
		status = 2
	}

	err = e.DB.Model(&r).Update("Status", status).Error
	e.changeUserAuthStatus(r.ID)
	return
}

// 删除角色
func (e *Enforcer) DeleteRole(roleID int64) (role models.Role, err error) {
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

	e.changeUserAuthStatus(roleID)
	e.DB.Delete(models.RoleAuthority{}, "role_id = ?", role.ID)
	e.DB.Delete(models.UserRole{}, "role_id = ?", role.ID)
	e.DB.Delete(&role)
	return
}

// 获取角色列表
func (e *Enforcer) GetRoles(args map[string]interface{}) (roles []models.Role, err error) {
	err = e.DB.Where(args).Preload("Auths").Find(&roles).Error
	return
}

// 更新角色权限
func (e *Enforcer) createOrUpdateRoleAuths(roleID int64, category int, authCodes []string) (err error) {
	chosedCodesString, _ := json.Marshal(authCodes)
	c1, c2, err := e.getCodesByFuncAuthCodes(authCodes)

	// 获取父级权限
	authCodesString, _ := json.Marshal(c1)
	apiAuthCodesString, _ := json.Marshal(c2)

	var ra models.RoleAuthority
	ra.RoleID = roleID
	ra.Category = category
	ra.FuncAuthCodes = string(authCodesString)
	ra.ApiAuthCodes = string(apiAuthCodesString)
	ra.ChosedCodes = string(chosedCodesString)

	err = e.DB.Where("role_id = ? and category = ?", ra.RoleID, category).First(&ra).Error
	if err != nil {
		err = e.DB.Create(&ra).Error
		return
	}

	err = e.DB.Model(&ra).Updates(models.RoleAuthority{FuncAuthCodes: string(authCodesString), ApiAuthCodes: string(apiAuthCodesString), ChosedCodes: string(chosedCodesString)}).Error

	if err == nil {
		e.changeUserAuthStatus(roleID)
	}
	return
}

func (e *Enforcer) changeUserAuthStatus(roleID int64) {
	var userRoles []models.UserRole
	err := e.DB.Where("role_id = ?", roleID).Preload("User").Find(&userRoles).Error
	if err == nil {
		for _, ur := range userRoles {
			if ur.User.ID > 0 {
				e.DB.Model(&ur.User).Update("AuthStatus", 2)
			}
		}
	}
	return
}
