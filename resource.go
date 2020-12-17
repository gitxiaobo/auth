package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

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
func (e *Enforcer) GetResources(dealerID int64) (resoures []Resource, err error) {
	jsonFile, err := os.Open(e.ResourceConfigPath)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &resoures)
	for index, r := range resoures {
		err = e.DB.Where("resource_key = ? and dealer_id =?", r.Key, dealerID).Find(&r.Items).Error
		if err == nil {
			resoures[index].Items = r.Items
		}
	}
	return
}

// 资源池设置
func (e *Enforcer) SetResource(key string, dealerID int64, resources []models.Resource) (err error) {
	e.DB.Where("resource_key = ? and dealer_id = ?", key, dealerID).Delete(&models.Resource{})
	var value []int64
	for _, r := range resources {
		r.ResourceKey = key
		r.DealerID = dealerID
		// err = e.DB.FirstOrCreate(&r, models.Resource{ResourceKey: r.ResourceKey, ResourceValue: r.ResourceValue, DealerID: r.DealerID}).Error
		var re models.Resource
		err = e.DB.Where("resource_key = ? and resource_value = ? and dealer_id = ?", r.ResourceKey, r.ResourceValue, r.DealerID).Find(&re).Error
		if err != nil {
			err = e.DB.Create(&r).Error
		}
		if err == nil {
			v, _ := strconv.Atoi(r.ResourceValue)
			value = append(value, int64(v))
		}
	}

	e.checkUserResource(key, value, dealerID)
	return
}

// 资源池增加单个资源
func (e *Enforcer) SetSingleResource(key string, dealerID int64, r models.Resource) (err error) {
	err = e.DB.FirstOrCreate(&r, models.Resource{ResourceKey: key, ResourceValue: r.ResourceValue, DealerID: dealerID}).Error
	return
}

//求交集
func intersect(slice1, slice2 []int64) []int64 {
	m := make(map[int64]int)
	nn := make([]int64, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 检查用户资源池
func (e *Enforcer) checkUserResource(key string, nv []int64, dealerID int64) (err error) {
	var userResources []models.UserResource
	e.DB.Where("resource_key = ? and dealer_id = ?", key, dealerID).Find(&userResources)
	for _, ur := range userResources {
		var value []int64
		json.Unmarshal([]byte(ur.ResourceValue), &value)
		v := intersect(nv, value)

		idsString, _ := json.Marshal(v)
		err = e.DB.Model(&ur).Update("resource_value", string(idsString)).Error
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
		resources, _ := e.GetResources(user.DealerID)
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
	ur.DealerID = user.DealerID

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
	resources, err := e.GetResources(0)
	if err == nil {
		for _, r := range resources {
			if r.Key == key {
				return r.FieldName
			}
		}
	}
	return ""
}
