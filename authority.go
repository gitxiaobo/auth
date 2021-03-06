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
	Codes    []string `json:"codes"` //不同客户端的功能权限码
	ApiCodes []string `json:"api_codes"`
	Platform int      `json:"platform"` // 0 - 共用, 1 - 平台方, 2 - 中间商
	Auths    []*Auth  `json:"children"`
}

// 获取权限配置文件数据
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
func (e *Enforcer) getCodesByFuncAuthCodes(funcAuthCodes []string) (funcCodes []string, apiCodes []string, err error) {
	auths, err := e.GetFuncAuths()
	for _, l1 := range auths {
		for index, code := range l1.Codes {
			if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
				funcCodes = append(funcCodes, l1.Codes[index])
				apiCodes = append(apiCodes, l1.ApiCodes...)
				funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
				if len(funcAuthCodes) == 0 {
					goto Finish
				}
			}
		}

		for _, l2 := range l1.Auths {
			for index, code := range l2.Codes {
				if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
					funcCodes = append(funcCodes, l1.Codes[index])
					apiCodes = append(apiCodes, l1.ApiCodes...)
					funcCodes = append(funcCodes, l2.Codes[index])
					apiCodes = append(apiCodes, l2.ApiCodes...)

					funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
					if len(funcAuthCodes) == 0 {
						goto Finish
					}
				}
			}

			for _, l3 := range l2.Auths {
				for index, code := range l3.Codes {
					if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
						funcCodes = append(funcCodes, l1.Codes[index])
						apiCodes = append(apiCodes, l1.ApiCodes...)
						funcCodes = append(funcCodes, l2.Codes[index])
						apiCodes = append(apiCodes, l2.ApiCodes...)
						funcCodes = append(funcCodes, l3.Codes[index])
						apiCodes = append(apiCodes, l3.ApiCodes...)
						funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
						if len(funcAuthCodes) == 0 {
							goto Finish
						}
					}
				}

				for _, l4 := range l3.Auths {
					for index, code := range l4.Codes {
						if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
							funcCodes = append(funcCodes, l1.Codes[index])
							apiCodes = append(apiCodes, l1.ApiCodes...)
							funcCodes = append(funcCodes, l2.Codes[index])
							apiCodes = append(apiCodes, l2.ApiCodes...)
							funcCodes = append(funcCodes, l3.Codes[index])
							apiCodes = append(apiCodes, l3.ApiCodes...)
							funcCodes = append(funcCodes, l4.Codes[index])
							apiCodes = append(apiCodes, l4.ApiCodes...)
							funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
							if len(funcAuthCodes) == 0 {
								goto Finish
							}
						}
					}

					for _, l5 := range l4.Auths {
						for index, code := range l5.Codes {
							if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
								funcCodes = append(funcCodes, l1.Codes[index])
								apiCodes = append(apiCodes, l1.ApiCodes...)
								funcCodes = append(funcCodes, l2.Codes[index])
								apiCodes = append(apiCodes, l2.ApiCodes...)
								funcCodes = append(funcCodes, l3.Codes[index])
								apiCodes = append(apiCodes, l3.ApiCodes...)
								funcCodes = append(funcCodes, l4.Codes[index])
								apiCodes = append(apiCodes, l4.ApiCodes...)
								funcCodes = append(funcCodes, l5.Codes[index])
								apiCodes = append(apiCodes, l5.ApiCodes...)
								funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
								if len(funcAuthCodes) == 0 {
									goto Finish
								}
							}
						}

						for _, l6 := range l5.Auths {
							for index, code := range l6.Codes {
								if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
									funcCodes = append(funcCodes, l1.Codes[index])
									apiCodes = append(apiCodes, l1.ApiCodes...)
									funcCodes = append(funcCodes, l2.Codes[index])
									apiCodes = append(apiCodes, l2.ApiCodes...)
									funcCodes = append(funcCodes, l3.Codes[index])
									apiCodes = append(apiCodes, l3.ApiCodes...)
									funcCodes = append(funcCodes, l4.Codes[index])
									apiCodes = append(apiCodes, l4.ApiCodes...)
									funcCodes = append(funcCodes, l5.Codes[index])
									apiCodes = append(apiCodes, l5.ApiCodes...)
									funcCodes = append(funcCodes, l6.Codes[index])
									apiCodes = append(apiCodes, l6.ApiCodes...)
									funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
									if len(funcAuthCodes) == 0 {
										goto Finish
									}
								}
							}

							for _, l7 := range l6.Auths {
								for index, code := range l7.Codes {
									if i := arrays.ContainsString(funcAuthCodes, code); i != -1 {
										funcCodes = append(funcCodes, l1.Codes[index])
										apiCodes = append(apiCodes, l1.ApiCodes...)
										funcCodes = append(funcCodes, l2.Codes[index])
										apiCodes = append(apiCodes, l2.ApiCodes...)
										funcCodes = append(funcCodes, l3.Codes[index])
										apiCodes = append(apiCodes, l3.ApiCodes...)
										funcCodes = append(funcCodes, l4.Codes[index])
										apiCodes = append(apiCodes, l4.ApiCodes...)
										funcCodes = append(funcCodes, l5.Codes[index])
										apiCodes = append(apiCodes, l5.ApiCodes...)
										funcCodes = append(funcCodes, l6.Codes[index])
										apiCodes = append(apiCodes, l6.ApiCodes...)
										funcCodes = append(funcCodes, l7.Codes[index])
										apiCodes = append(apiCodes, l7.ApiCodes...)
										funcAuthCodes = append(funcAuthCodes[:i], funcAuthCodes[i+1:]...)
										if len(funcAuthCodes) == 0 {
											goto Finish
										}
									}
								}
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

// 获取API配置文件
func (e *Enforcer) getAPIAuths() (auths []APIAuth, err error) {
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
func (e *Enforcer) getApiAuthCode(url string, method string) (code string, err error) {
	auths, _ := e.getAPIAuths()
	for _, auth := range auths {
		if auth.Method == method {
			if ok, _ := regexp.MatchString("^"+auth.URL+"$", url); ok {
				code = auth.Code
				return
			}
		}
	}
	return
}
