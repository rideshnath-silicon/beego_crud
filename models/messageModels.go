package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type NewMessageCode struct {
	Code    string
	English string
	Hindi   string
}

func InsertNewMessage(data NewMessageCode) (string, error) {
	o := orm.NewOrm()
	// orm.Debug = true
	sql := "INSERT INTO language_code (language_code, created_date) VALUES (?, ?) RETURNING langcodeid"
	params := []interface{}{data.Code, time.Now()} // Replace with your actual values
	var langId int64
	err := o.Raw(sql, params...).QueryRow(&langId)
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	sql1 := "INSERT INTO engilsh_lang_message (langcodeid, value) VALUES (?, ?)"
	paramsql1 := []interface{}{langId, data.English}
	_, err = o.Raw(sql1, paramsql1...).Exec()
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	sql2 := "INSERT INTO hindi_lang_message (langcodeid, value) VALUES (?, ?)"
	paramsql2 := []interface{}{langId, data.Hindi}
	_, err = o.Raw(sql2, paramsql2...).Exec()
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	var langCode string
	sql = "select language_code from language_code where langcodeid = ?"
	err = o.Raw(sql, langId).QueryRow(&langCode)
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	return "Generated Language code is :- " + langCode, nil	
}
