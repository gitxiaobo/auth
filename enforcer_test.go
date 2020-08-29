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

	// err = e.SwitchRoleStatus(1)
	// fmt.Println(err)

	// roles, _ := e.GetRoles(map[string]interface{}{})
	// fmt.Println(roles)
	// e.getCodesByFuncAuthCodes([]string{"1100", "1200"})
	// e.CreateOrUpdateUser(2)
	// e.CreateOrUpdateUserRole(1, []int64{1, 4})
	// roles, _ := e.GetUserRoles(1)
	// fmt.Println(roles)
	// e.CreateOrUpdateUserResouce(2, "region", []int64{1})

	// codes, _ := e.GetRoleFuncAuths(1)
	// fmt.Println(codes)
	// r, _ := e.GetUserResources(1, "region")
	// e.CreateOrUpdateRole(models.Role{Name: "1687", ID: 1}, 1, []string{"123", "34f"})
	// e.CreateOrUpdateUser(1)
	// e.CreateOrUpdateRole(models.Role{Name: "16844447"}, 1, []string{"123", "34f"})
	// e.DeleteRole(2)
	// e.NomarlUserAuthStatus(2)
	b, err := e.CheckApiAuth(1, "/454353453", "ddd")
	// e.CreateOrUpdateUserRole(1, []int64{2})

	// e.CreateOrUpdateUserRole(2, []int64{1})
	// e.createOrUpdateRoleAuths(1, 1, []string{"11301"})
	// s, err := e.GetUserFuncAuths(1)
	fmt.Println(b, err)
	// e.DeleteRole(1)

	// e.createOrUpdateRoleAuths(4, 1, []string{"1100", "200"})
	// e.CreateOrUpdateUserRole(1, []int64{1, 2})

	// codes, err := e.getUserApiAuths(1)
	// fmt.Println(codes, err)

	// e.CreateOrUpdateUser(1)
	// e.CreateOrUpdateRole(1, "技术人员")
	// e.CreateOrUpdateUserRole(1, []int64{1})

	// e.CreateOrUpdateRoleAuths(1, []string{"1", "2"})

	// roles, err := e.GetUserFuncAuths(1)
	// fmt.Println(roles, err)
	// e.CreateOrUpdateRoleAuths(1, []string{"1100", "200"})

	// b, err := e.CheckApiAuth(1, "/api/customers/4", "get")
	// fmt.Println("========")
	// fmt.Println(b, err)

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
