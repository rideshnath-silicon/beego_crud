package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func InsertNewHomeSetting(data InserNewHomeSettingRequest) (interface{}, error) {
	o := orm.NewOrm()
	var hsetting = HomeSetting{
		Section:     data.Section,
		Type:        data.Type,
		Key:         data.Key,
		Value:       data.Value,
		CreatedDate: time.Now(),
	}
	num, err := o.Insert(&hsetting)
	if err != nil {
		return nil,  errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return nil, errors.New("DATABASE_ERROR")
	}
	return hsetting, nil
}

func UpdateHomeSeting(data UpdateHomeSetingRequest) (interface{}, error) {
	o := orm.NewOrm()
	var hsetting = HomeSetting{
		Id:         data.Id,
		Section:    data.Section,
		Type:       data.Type,
		Key:        data.Key,
		Value:      data.Value,
		UpdateDate: time.Now(),
	}
	num, err := o.Update(&hsetting, "id", "section", "type", "key", "value", "update_at")
	if err != nil {
		return nil,  errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return nil, errors.New("DATABASE_ERROR")
	}
	return hsetting, nil
}

func GetHomeSetting(id uint) (HomeSetting, error) {
	o := orm.NewOrm()
	var hsetting = HomeSetting{Id: id}
	num, err := o.QueryTable(new(HomeSetting)).Filter("id", id).All(&hsetting)
	if err != nil {
		return hsetting,  errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return hsetting, errors.New("LOGIN_ERROR")
	}
	return hsetting, nil
}

func UserWiseHomeseting(id uint) (interface{}, error) {
	o := orm.NewOrm()
	var userhome UserWiseHomeRequest
	sqlQuery := `
		SELECT hs.section, hs.type, hs.value, u.first_name, u.last_name
		FROM home_setting as hs
		JOIN users as u ON u.id = hs.id
		WHERE hs.id = ?
	`
	err := o.Raw(sqlQuery, id).QueryRow(&userhome)
	if err != nil {
		return nil,  errors.New("DATABASE_ERROR")
	}

	return userhome, nil
}

func DeleteHomeSetting(id int) (interface{}, error) {
	o := orm.NewOrm()
	var hsetting = HomeSetting{Id: uint(id)}
	num, err := o.Delete(&hsetting)
	if err != nil {
		return "",  errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return "",  errors.New("DATABASE_ERROR")
	}
	return "DATA_DELETE", nil
}
