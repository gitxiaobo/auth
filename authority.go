package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/wxnacy/wgo/arrays"
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
	return
}

// func (e *Enforcer) GetFuncAuthArray() (arr []Auth, err error) {
// 	auths, err := e.GetFuncAuths()
// 	for _, l1 := range auths {
// 		arr = append(arr, Auth{Code: l1.Code, ApiCodes: l1.ApiCodes})
// 		for _, l2 := range l1.Auths {
// 			arr = append(arr, Auth{Code: l2.Code, ApiCodes: l2.ApiCodes})
// 			for _, l3 := range l2.Auths {
// 				arr = append(arr, Auth{Code: l3.Code, ApiCodes: l3.ApiCodes})
// 				for _, l4 := range l3.Auths {
// 					arr = append(arr, Auth{Code: l4.Code, ApiCodes: l4.ApiCodes})
// 					for _, l5 := range l4.Auths {
// 						arr = append(arr, Auth{Code: l5.Code, ApiCodes: l5.ApiCodes})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// 字符串数组去重
func removeDuplicateElement(arr []string) []string {
	result := make([]string, 0, len(arr))
	temp := map[string]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// 获取父级功能权限码和api权限码
func (e *Enforcer) GetCodesByFuncAuthCodes(funcAuthCodes []string) (funcCodes []string, apiCodes []string, err error) {
	auths, err := e.GetFuncAuths()
	for _, l1 := range auths {
		index := arrays.ContainsString(funcAuthCodes, l1.Code)
		if index != -1 {
			funcCodes = append(funcCodes, l1.Code)
			apiCodes = append(apiCodes, l1.ApiCodes...)

			funcAuthCodes = append(funcAuthCodes[:index], funcAuthCodes[index+1:]...)
			if len(funcAuthCodes) == 0 {
				goto Finish
			}
		}

		for _, l2 := range l1.Auths {
			index := arrays.ContainsString(funcAuthCodes, l2.Code)
			if index != -1 {
				funcCodes = append(funcCodes, l1.Code)
				apiCodes = append(apiCodes, l1.ApiCodes...)
				funcCodes = append(funcCodes, l2.Code)
				apiCodes = append(apiCodes, l2.ApiCodes...)

				funcAuthCodes = append(funcAuthCodes[:index], funcAuthCodes[index+1:]...)
				if len(funcAuthCodes) == 0 {
					goto Finish
				}
			}

			for _, l3 := range l2.Auths {
				index := arrays.ContainsString(funcAuthCodes, l3.Code)
				if index != -1 {
					funcCodes = append(funcCodes, l1.Code)
					apiCodes = append(apiCodes, l1.ApiCodes...)
					funcCodes = append(funcCodes, l2.Code)
					apiCodes = append(apiCodes, l2.ApiCodes...)
					funcCodes = append(funcCodes, l3.Code)
					apiCodes = append(apiCodes, l3.ApiCodes...)

					funcAuthCodes = append(funcAuthCodes[:index], funcAuthCodes[index+1:]...)
					if len(funcAuthCodes) == 0 {
						goto Finish
					}
				}

				for _, l4 := range l3.Auths {
					index := arrays.ContainsString(funcAuthCodes, l3.Code)
					if index != -1 {
						funcCodes = append(funcCodes, l1.Code)
						apiCodes = append(apiCodes, l1.ApiCodes...)
						funcCodes = append(funcCodes, l2.Code)
						apiCodes = append(apiCodes, l2.ApiCodes...)
						funcCodes = append(funcCodes, l3.Code)
						apiCodes = append(apiCodes, l3.ApiCodes...)
						funcCodes = append(funcCodes, l4.Code)
						apiCodes = append(apiCodes, l4.ApiCodes...)

						funcAuthCodes = append(funcAuthCodes[:index], funcAuthCodes[index+1:]...)
						if len(funcAuthCodes) == 0 {
							goto Finish
						}
					}

					for _, l5 := range l4.Auths {
						index := arrays.ContainsString(funcAuthCodes, l3.Code)
						if index != -1 {
							funcCodes = append(funcCodes, l1.Code)
							apiCodes = append(apiCodes, l1.ApiCodes...)
							funcCodes = append(funcCodes, l2.Code)
							apiCodes = append(apiCodes, l2.ApiCodes...)
							funcCodes = append(funcCodes, l3.Code)
							apiCodes = append(apiCodes, l3.ApiCodes...)
							funcCodes = append(funcCodes, l4.Code)
							apiCodes = append(apiCodes, l4.ApiCodes...)
							funcCodes = append(funcCodes, l5.Code)
							apiCodes = append(apiCodes, l5.ApiCodes...)

							funcAuthCodes = append(funcAuthCodes[:index], funcAuthCodes[index+1:]...)
							if len(funcAuthCodes) == 0 {
								goto Finish
							}
						}
					}
				}
			}
		}
	}

Finish:
	funcCodes = removeDuplicateElement(funcCodes)
	apiCodes = removeDuplicateElement(apiCodes)
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
