package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type AuthModule struct {
	Name     string     `json:"name"`
	ModuleID int        `json:"module_id"`
	Auths    []FuncAuth `json:"auths"`
}

type FuncAuth struct {
	Name         string    `json:"name"`
	Desc         string    `json:"desc"`
	FuncAuthCode string    `json:"func_auth_code"`
	ApiAuths     []ApiAuth `json:"api_auths"`
}

type ApiAuth struct {
	AuthCode string `json:"auth_code"`
	URL      string `json:"url"`
	Method   string `json:"method"`
}

// 获取所有配置权限，分模块返回
func (e *Enforcer) GetAuths() (auths []AuthModule, err error) {
	jsonFile, err := os.Open(e.AuthConfigPath)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &auths)
	return
}

// 获取api权限码
func (e *Enforcer) GetApiAuthCode(url string, method string) (code string, err error) {
	auths, _ := e.GetAuths()
	for _, moduleAuth := range auths {
		for _, funcAuth := range moduleAuth.Auths {
			for _, apiAuth := range funcAuth.ApiAuths {
				if apiAuth.Method == method {
					if ok, _ := regexp.MatchString(apiAuth.URL, url); ok {
						code = apiAuth.AuthCode
						return
					}
				}
			}
		}
	}
	return
}
