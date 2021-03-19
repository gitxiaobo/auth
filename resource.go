package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gitxiaobo/auth/models"
	"github.com/wxnacy/wgo/arrays"
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
	r.ResourceKey = key
	r.DealerID = dealerID

	var or models.Resource
	e.DB.Where("resource_key = ? and resource_value = ? and dealer_id = ?", key, r.ResourceValue, dealerID).Find(&or)
	if or.ID == 0 {
		err = e.DB.Create(&r).Error
	} else {
		err = e.DB.Model(&or).Update(map[string]interface{}{"area_id": r.AreaID, "resource_name": r.ResourceName}).Error
	}

	// err = e.DB.FirstOrCreate(&r, models.Resource{ResourceKey: key, ResourceValue: r.ResourceValue, DealerID: dealerID}).Error
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
	e.DB.Where("resource_key = ? and dealer_id = ? and all_area = 0", key, dealerID).Find(&userResources)
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
	var ur []models.UserResource
	err = e.DB.Where("user_id = ? and resource_key = ?", user.ID, key).Find(&ur).Error

	for _, u := range ur {
		fieldName = u.FieldName
		if u.AllArea == 1 {
			var vals []string
			e.DB.Table("auth_resources").Where("dealer_id = ? and resource_key = ? and area_id = ?", user.DealerID, key, u.AreaID).Pluck("resource_value", &vals)
			for _, vl := range vals {
				id, _ := strconv.Atoi(vl)
				value = append(value, int64(id))
			}
		} else {
			var v []int64
			json.Unmarshal([]byte(u.ResourceValue), &v)
			value = append(value, v...)
		}
	}
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

	// json.Unmarshal([]byte(ur.ResourceValue), &value)
	// fieldName = ur.FieldName
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
func (e *Enforcer) CreateOrUpdateUserResouce(userID int64, key string, ids []int64, areaID int, all int) (err error) {
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
	ur.AreaID = areaID
	ur.AllArea = all

	err = e.DB.Where("user_id = ? and resource_key = ? and area_id = ?", user.ID, key, areaID).First(&ur).Error
	if err != nil {
		e.DB.Create(&ur)
		return
	}

	// err = e.DB.Model(&ur).Update("resource_value", string(idsString)).Error
	err = e.DB.Model(&ur).Update(map[string]interface{}{"resource_value": string(idsString), "all_area": all}).Error
	return
}

// 添加自己为自己的资源
func (e *Enforcer) AddSelfToResoure(userID int64, key string, id int64, areaID int) (err error) {
	user, err := e.findUserByUserID(userID)
	if err != nil {
		return
	}

	var ur models.UserResource
	ur.DealerID = user.DealerID
	ur.ResourceKey = key
	ur.AreaID = areaID
	ur.UserID = user.ID
	err = e.DB.Where("user_id = ? and resource_key = ? and area_id = ?", user.ID, key, areaID).First(&ur).Error
	if err != nil {
		ur.FieldName = e.getFieldNameByKey(key)
		idsString, _ := json.Marshal([]int64{id})
		ur.ResourceValue = string(idsString)
		e.DB.Create(&ur)
		return
	}

	var v []int64
	json.Unmarshal([]byte(ur.ResourceValue), &v)

	if arrays.ContainsInt(v, id) == -1 {
		v = append(v, id)
		idsString, _ := json.Marshal(v)
		err = e.DB.Model(&ur).Update(map[string]interface{}{"resource_value": string(idsString)}).Error
	}

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

//	查询拥有某个资源的所有人
func (e *Enforcer) GetFatherUserIDs(key string, sourceValue string) (ids []int64, err error) {
	var resource models.Resource
	err = e.DB.Where("resource_key = ? and resource_value = ?", key, sourceValue).First(&resource).Error
	if err == nil {
		var userIDs []int64
		v, _ := strconv.Atoi(resource.ResourceValue)
		err = e.DB.Table("auth_user_resources").Where("resource_key = ? and (JSON_CONTAINS(resource_value, JSON_ARRAY(?)) or (area_id = ? and all_area = 1))", resource.ResourceKey, v, resource.AreaID).Pluck("user_id", &userIDs).Error
		if err == nil {
			err = e.DB.Table("auth_users").Where("id in (?)", userIDs).Pluck("user_id", &ids).Error
		}
	}
	return
}

//	查询拥有某个资源且拥有某个角色的所有人
func (e *Enforcer) GetUserIDsByResourceAndRole(key string, sourceValue string, roleID int64) (ids []int64, err error) {
	var resource models.Resource
	err = e.DB.Where("resource_key = ? and resource_value = ?", key, sourceValue).First(&resource).Error
	if err == nil {
		var userIDs []int64
		v, _ := strconv.Atoi(resource.ResourceValue)
		err = e.DB.Table("auth_user_resources").Where("resource_key = ? and (JSON_CONTAINS(resource_value, JSON_ARRAY(?)) or (area_id = ? and all_area = 1))", resource.ResourceKey, v, resource.AreaID).Pluck("user_id", &userIDs).Error
		if err == nil {
			err = e.DB.Table("auth_users").Where("id in (?) and id in (select user_id from auth_user_roles where role_id = ?)", userIDs, roleID).Pluck("user_id", &ids).Error
		}
	}
	return
}
