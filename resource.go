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
	Platform  int                    `json:"platform"` // 0 - 共用, 1 - 平台方, 2 - 中间商
	Items     []models.Resource      `json:"items"`
	Options   map[string]interface{} `json:"options"`
}

// 获取数据资源配置文件数据
func (e *Enforcer) GetResources() (resoures []Resource, err error) {
	jsonFile, err := os.Open(e.ResourceConfigPath)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &resoures)
	for index, r := range resoures {
		err = e.DB.Where("resource_key = ?", r.Key).Find(&r.Items).Error
		if err == nil {
			resoures[index].Items = r.Items
		}
	}
	return
}

// 资源池设置
func (e *Enforcer) SetResource(key string, resources []models.Resource) (err error) {
	e.DB.Where("resource_key = ?", key).Delete(&models.Resource{})
	for _, r := range resources {
		r.ResourceKey = key
		err = e.DB.FirstOrCreate(&r, models.Resource{ResourceKey: r.ResourceKey, ResourceValue: r.ResourceValue}).Error
	}
	return
}

// 获取用户某个资源的资源值
func (e *Enforcer) GetUserResourcesByKey(userID int64, key string) (value []int64, fieldName string, err error) {
	user, err := e.findUserByUserID(userID)
	if err != nil {
		return
	}
	var ur models.UserResource
	err = e.DB.Where("user_id = ? and resource_key = ?", user.ID, key).First(&ur).Error
	if err != nil {
		resources, _ := e.GetResources()
		for _, r := range resources {
			if r.Key == key {
				fieldName = r.FieldName
				err = nil
				break
			}
		}
		return
	}

	json.Unmarshal([]byte(ur.ResourceValue), &value)
	fieldName = ur.FieldName
	return
}

// 获取用户资源列表
func (e *Enforcer) GetUserResources(userID int64, key string) (resoures []models.UserResource, err error) {
	user, err := e.findUserByUserID(userID)
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

// 用户资源管理
func (e *Enforcer) CreateOrUpdateUserResouce(userID int64, key string, ids []int64) (err error) {
	user, err := e.findUserByUserID(userID)
	if err != nil {
		return
	}

	idsString, _ := json.Marshal(ids)

	fieldName := e.getFieldNameByKey(key)

	var ur models.UserResource
	ur.UserID = user.ID
	ur.ResourceKey = key
	ur.ResourceValue = string(idsString)
	ur.FieldName = fieldName

	err = e.DB.Where("user_id = ? and resource_key = ?", user.ID, key).First(&ur).Error
	if err != nil {
		e.DB.Create(&ur)
		return
	}

	err = e.DB.Model(&ur).Update("resource_value", string(idsString)).Error
	return
}

// 通过key得到查询的字段
func (e *Enforcer) getFieldNameByKey(key string) string {
	resources, err := e.GetResources()
	if err == nil {
		for _, r := range resources {
			if r.Key == key {
				return r.FieldName
			}
		}
	}
	return ""
}
