package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Auth struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Code     string   `json:"code"`
	ApiCodes []string `json:"api_codes"`
	Auths    []*Auth  `json:"children"`
}

func (e *Enforcer) GetFuncAuths() (auths []Auth, err error) {
	jsonFile, err := os.Open(e.FuncAuthConfigPath)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &auths)

	// Jsondata, _ := json.Marshal(auths)
	// fmt.Println(string(Jsondata))
	return
}

func (e *Enforcer) GetFuncAuthArray() (arr []Auth, err error) {
	auths, err := e.GetFuncAuths()
	for _, l1 := range auths {
		arr = append(arr, Auth{Code: l1.Code, ApiCodes: l1.ApiCodes})
		for _, l2 := range l1.Auths {
			arr = append(arr, Auth{Code: l2.Code, ApiCodes: l2.ApiCodes})
			for _, l3 := range l2.Auths {
				arr = append(arr, Auth{Code: l3.Code, ApiCodes: l3.ApiCodes})
				for _, l4 := range l3.Auths {
					arr = append(arr, Auth{Code: l4.Code, ApiCodes: l4.ApiCodes})
					for _, l5 := range l4.Auths {
						arr = append(arr, Auth{Code: l5.Code, ApiCodes: l5.ApiCodes})
					}
				}
			}
		}
	}
	return
}

type APIAuth struct {
	Code   string `json:"code"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

func (e *Enforcer) GetAPIAuths() (auths []APIAuth, err error) {
	jsonFile, err := os.Open(e.ApiAuthConfigPath)

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
	auths, _ := e.GetAPIAuths()
	for _, auth := range auths {
		if auth.Method == method {
			if ok, _ := regexp.MatchString(auth.URL, url); ok {
				code = auth.Code
				return
			}
		}
	}
	return
}
