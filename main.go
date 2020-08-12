package main

import (
	"fmt"

	"auth/api"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// var db *gorm.DB

func main() {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/auth?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	e, err := api.NewEnforcer(db, "config/auth.json")

	// e.CreateTable()

	// b, err := e.CheckApiAuth(1, "111")
	// fmt.Println(b, err)

	auths, err := e.GetAuths()
	fmt.Println(auths, err)

	// e.CreateOrUpdateUser(1)
	// e.CreateOrUpdateRole(1, "技术人员")
	// e.CreateOrUpdateUserRole(1, []int64{1})
	// auths, err := e.GetUserApiAuths(1)
	// err = e.CreateOrUpdateRoleAuths(2, []string{"1223", "2d34"})
	// fmt.Println(err)

	defer db.Close()
}
