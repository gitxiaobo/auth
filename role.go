package auth

import (
	"encoding/json"

	"github.com/gitxiaobo/auth/errors"

	"github.com/gitxiaobo/auth/models"
)

// 创建或更新角色
func (e *Enforcer) CreateOrUpdateRole(role models.Role, codes []string) (rl models.Role, err error) {
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
		err = e.CreateOrUpdateRoleAuths(role.ID, codes)
	}

	err = e.DB.Where("id = ?", role.ID).Preload("Auths").First(&rl).Error
	return
}

// 删除角色
func (e *Enforcer) DeleteRole(roleID int64) (role models.Role, err error) {
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

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
func (e *Enforcer) CreateOrUpdateRoleAuths(roleID int64, authCodes []string) (err error) {
	var role models.Role
	err = e.DB.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return
	}

	c1, c2, err := e.GetCodesByFuncAuthCodes(authCodes)
	chosedCodesString, _ := json.Marshal(authCodes)

	// 获取父级权限
	authCodesString, _ := json.Marshal(c1)
	apiAuthCodesString, _ := json.Marshal(c2)

	var ra models.RoleAuthority
	ra.RoleID = roleID
	ra.FuncAuthCodes = string(authCodesString)
	ra.ApiAuthCodes = string(apiAuthCodesString)
	ra.ChosedCodes = string(chosedCodesString)

	err = e.DB.Where("role_id = ?", ra.RoleID).First(&ra).Error
	if err != nil {
		err = e.DB.Create(&ra).Error
		return
	}

	err = e.DB.Model(&ra).Updates(models.RoleAuthority{FuncAuthCodes: string(authCodesString), ApiAuthCodes: string(apiAuthCodesString), ChosedCodes: string(chosedCodesString)}).Error
	return
}
