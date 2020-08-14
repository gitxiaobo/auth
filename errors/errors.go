package errors

import "errors"

// Global errors for rbac defined here
var (
	USER_NOT_FOUND   = errors.New("用户未找到")
	ROLE_NOT_FOUND   = errors.New("角色未找到")
	ROLE_NAME_REPEAT = errors.New("角色名重复")
	PARAMS_ERROR     = errors.New("参数错误")
	DB_ERROR         = errors.New("数据库错误")
)
