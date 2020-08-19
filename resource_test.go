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
	e.CreateTable()
	// r, _ := e.GetResources()
	// fmt.Println(r)

	r, s, _ := e.GetUserResourcesByKey(1, "user")
	fmt.Println(r, s)
	// e.CreateOrUpdateUser(1)
	// b := e.isSuperAdmin(1)
	// fmt.Println(b)

	// e.CreateOrUpdateUserResouce(1, "user", []int64{1, 2, 3})
	defer db.Close()
}
