package models

import (
	"CarCrudv2/helpers"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func GetUserByEmail(username string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	// orm.Debug = true
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).All(&user)
	if err != nil {
		return user,  errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("LOGIN_ERROR")
	}
	return user, nil
}

func LoginUser(username string, pass string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).Filter("password", pass).All(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func GetUserDetails(id interface{}) (Users, error) {
	o := orm.NewOrm()
	// orm.Debug = true
	var user Users
	num, err := o.QueryTable(new(Users)).Filter("id", id).All(&user, "first_name", "last_name", "email", "password", "phone_number")
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func GetAllUser() ([]orm.Params, error) {
	o := orm.NewOrm()
	// orm.Debug = true
	var user []orm.Params
	sqlQuery := `
		SELECT u.id as user_id , u.first_name as name, u.last_name as last_name , u.email  as user_email , u.age as user_age, c.name as country_name
		FROM users as u
		JOIN mod_country_master as c ON c.country_id = u.country	
		
	`
	_, err := o.Raw(sqlQuery).Values(&user)
	if err != nil {
		return nil,  errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func InsertNewUser(Data NewUserRequest) (Users, error) {
	o := orm.NewOrm()
	var user Users
	pass, err := helpers.HashData(Data.Password)
	if err != nil {
		return user, err
	}
	user = Users{
		FirstName:   Data.FirstName,
		LastName:    Data.LastName,
		Country:     Data.Country,
		Email:       Data.Email,
		PhoneNumber: Data.PhoneNumber,
		Age:         Data.Age,
		Password:    pass,
		Role:        Data.Role,
		CreatedAt:   time.Now(),
	}
	num, err := o.Insert(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func ImportNewUser(Data NewUserRequest) (Users, error) {
	o := orm.NewOrm()
	var user = Users{
		FirstName:   Data.FirstName,
		LastName:    Data.LastName,
		Country:     Data.Country,
		Email:       Data.Email,
		PhoneNumber: Data.PhoneNumber,
		Age:         Data.Age,
		Password:    Data.Password,
		Role:        Data.Role,
		CreatedAt:   time.Now(),
	}
	num, err := o.Insert(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}
func UpdateUser(Data UpdateUserRequest) (interface{}, error) {
	var user = Users{
		Id:          Data.Id,
		FirstName:   Data.FirstName,
		LastName:    Data.LastName,
		Country:     Data.Country,
		Email:       Data.Email,
		Age:         Data.Age,
		Role:        Data.Role,
		UpdatedAt:   time.Now(),
		PhoneNumber: Data.PhoneNumber,
	}
	o := orm.NewOrm()
	num, err := o.Update(&user, "id", "first_name", "last_name", "country", "email", "age", "role", "updated_at", "phone_number")
	if err != nil {
		return nil,errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return "DATA_UPDATED", nil
}

func ResetPassword(Password string, id float64) (interface{}, error) {
	pass, err := helpers.HashData(Password)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	var user = Users{Id: uint(id), Password: pass}
	num, err := o.Update(&user, "password")
	if err != nil {
		return num, errors.New("DATABASE_ERROR")
	}
	return "PASSWORD_RESET", nil
}

func UpadteOtpForEmail(id uint, otp string) (string, error) {
	o := orm.NewOrm()
	var user = Users{Id: id, Otp: otp, Verified: "no"}
	num, err := o.Update(&user, "otp", "verified")
	if err != nil {
		return "num", errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return "user", errors.New("DATABASE_ERROR")
	}
	return "OTP_SENT", nil
}

func VerifyEmailOTP(username string, otp string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).Filter("otp", otp).All(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user,errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func UpdateVerified(id uint) error {
	o := orm.NewOrm()
	var user = Users{Id: id, Verified: "yes"}
	num, err := o.Update(&user, "verified")
	if err != nil {
		return errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return errors.New("DATABASE_ERROR")
	}
	return nil
}

func UpdateColumnOTP(id uint, otp string) {
	<-time.After(30 * time.Second)
	o := orm.NewOrm()
	var user = Users{Id: id, Otp: otp}
	_, err := o.Update(&user, "otp")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<OTP IS Expired>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.")
	if err != nil {
		log.Fatalf("Error updating column: %s", err)
	}
}

func GetCountryWiseCountUser() ([]GetCountryWiseUserRequest, error) {
	o := orm.NewOrm()
	var countries []GetCountryWiseUserRequest

	_, err := o.Raw("SELECT country,COUNT(country) AS count FROM users GROUP BY country").QueryRows(&countries)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func GetVerifiedUsers() ([]Users, error) {
	o := orm.NewOrm()
	var user []Users
	num, err := o.QueryTable(new(Users)).Filter("verified", "yes").OrderBy("id").All(&user)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return nil, errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func SearchUser(search string) ([]Users, error) {
	o := orm.NewOrm()
	var user []Users
	// orm.Debug = true
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("first_name__icontains", search).Or("email__icontains", search).Or("last_name__icontains", search).Or("role__icontains", search)).All(&user)
	if err != nil {
		return nil, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return nil, errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func DeleteUser(id uint) (string, error) {
	o := orm.NewOrm()
	var user = Users{Id: id}
	num, err := o.Delete(&user)
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return "", errors.New("DATABASE_ERROR")
	}
	return "Data_Delete", nil
}