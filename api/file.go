package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AuthModule struct {
	Name     string
	ModuleID int
	Auths    []Auth
}

type Auth struct {
	Name         string
	Desc         string
	FunAuthCode  string
	ApiAuthCodes []string
}

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
