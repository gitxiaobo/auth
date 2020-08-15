package auth

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Test(t *testing.T) {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/auth?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	e, err := NewEnforcer(db, "config/auth.json", "config/api_auth.json", "config/resource.json")

	// auths, err := e.GetFuncAuths()
	// fmt.Println(auths, err)

	// roles, _ := e.GetRoles(map[string]interface{}{})
	// fmt.Println(roles)
	e.CreateOrUpdateUser(1)
	e.CreateOrUpdateUserRole(1, []int64{1})
	roles, _ := e.GetUserRoles(1)
	fmt.Println(roles)
	// e.CreateOrUpdateUserResouce(2, "region", []int64{1})

	// codes, _ := e.GetRoleFuncAuths(1)
	// fmt.Println(codes)
	// r, _ := e.GetUserResources(1, "region")
	fmt.Println(e, err)

	// e.CreateOrUpdateUser(1)
	// e.CreateOrUpdateRole(1, "技术人员")
	// e.CreateOrUpdateUserRole(1, []int64{1})

	// e.CreateOrUpdateRoleAuths(1, []string{"1", "2"})

	// roles, err := e.GetUserFuncAuths(1)
	// fmt.Println(roles, err)
	e.CreateOrUpdateRoleAuths(1, []string{"1100", "200"})

	b, err := e.CheckApiAuth(1, "/api/customers/4", "get")
	fmt.Println("========")
	fmt.Println(b, err)

	// auths, err := e.GetAuths()
	// fmt.Println(auths, err)

	// e.CreateOrUpdateUser(1)
	// e.CreateOrUpdateRole(1, "技术人员")
	// e.CreateOrUpdateUserRole(1, []int64{1})
	// auths, err := e.GetUserApiAuths(1)
	// err = e.CreateOrUpdateRoleAuths(2, []string{"1223", "2d34"})
	// fmt.Println(err)

	defer db.Close()
}
