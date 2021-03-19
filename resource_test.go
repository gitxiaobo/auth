package auth

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestResource(t *testing.T) {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/auth?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	e, _ := NewEnforcer(db, "config/auth.json", "config/api_auth.json", "config/resource.json")
	ids, _ := e.GetUserIDsByResourceAndRole("user", "30", 16)
	fmt.Println(ids)
	// e.CreateTable()
	// e.AddSelfToResoure(1, "user", 8, 2)

	// e.SetSingleResource("user", 2, models.Resource{ResourceName: "222", ResourceValue: "1111", AreaID: 233})
	// e.getFieldNameByKey("user")
	// e.CreateOrUpdateRole(models.Role{Name: "53333", ID: 8}, []string{"121101", "151104", "251104"})
	// e.SetResource("user", 1, []models.Resource{models.Resource{ResourceName: "xxxx", ResourceValue: "23333", AreaID: 2}})
	// e.SetResource("user", 3, []models.Resource{models.Resource{ResourceName: "xxxx", ResourceValue: "23333"}, models.Resource{ResourceName: "xxxx", ResourceValue: "23", AreaID: 1}})
	// r, _ := e.GetResources()
	// fmt.Println(r)

	// r, s, _ := e.GetUserResourcesByKey(1, "user")
	// fmt.Println(r, s)
	// e.CreateOrUpdateUser(1, 1)
	// b := e.isSuperAdmin(1)
	// fmt.Println(b)

	// e.CreateOrUpdateUserResouce(1, "user", []int64{22, 3}, 2, 1)
	defer db.Close()
}
