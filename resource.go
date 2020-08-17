package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gitxiaobo/auth/models"
)

type Resource struct {
	Name      string                 `json:"name"`
	FieldName string                 `json:"field_name"`
	Key       string                 `json:"key"`
	Table     string                 `json:"table"`
	Items     map[string]interface{} `json:"items"`
}

func (e *Enforcer) GetResources() (resoures []Resource, err error) {
	jsonFile, err := os.Open(e.ResourceConfigPath)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &resoures)
	return
}

func (e *Enforcer) GetUserResourcesByKey(userID int64, key string) (value []int64, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return
	}
	var ur models.UserResource
	err = e.DB.Where("user_id = ? and resource_key = ?", user.ID, key).First(&ur).Error
	if err != nil {
		return
	}

	json.Unmarshal([]byte(ur.ResourceValue), &value)
	return
}

func (e *Enforcer) GetUserResources(userID int64, key string) (resoures []models.UserResource, err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return
	}

	db := e.DB.Where("user_id = ?", user.ID)
	if len(key) > 0 {
		db = db.Where("resource_key = ?", key)
	}

	err = db.Find(&resoures).Error
	return
}

func (e *Enforcer) CreateOrUpdateUserResouce(userID int64, key string, ids []int64) (err error) {
	var user models.User
	err = e.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return
	}

	idsString, _ := json.Marshal(ids)

	var ur models.UserResource
	ur.UserID = user.ID
	ur.ResourceKey = key
	ur.ResourceValue = string(idsString)

	err = e.DB.Where("user_id = ? and resource_key = ?", user.ID, key).First(&ur).Error
	if err != nil {
		e.DB.Create(&ur)
		return
	}

	err = e.DB.Model(&ur).Update("resource_value", string(idsString)).Error
	return
}
