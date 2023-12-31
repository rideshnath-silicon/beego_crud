package test

import (
	"CarCrudv2/models"
	"fmt"
	"reflect"
	"testing"

	"github.com/beego/beego/v2/client/httplib"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func TestUserModels(t *testing.T) {
	t.Run("Register new user", func(t *testing.T) {
		TruncateTable("users")
		var user = models.NewUserRequest{
			FirstName:   "Ridesh",
			LastName:    "Nath",
			Email:       "rideshnath.siliconithub@gmail.com",
			PhoneNumber: "1234567890",
			Role:        "developer",
			Country:     1,
			Age:         24,
			Password:    "123456",
		}
		data, err := models.InsertNewUser(user)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Id == 0 {
			t.Errorf("error in register")
			return
		}
		t.Log(data)
	})
	t.Run("Get All User", func(t *testing.T) {
		data, err := models.GetAllUser()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if len(data) == 0 {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Get UserByEmail", func(t *testing.T) {
		data, err := models.GetUserByEmail("rideshnath.siliconithub@gmail.com")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Email == "" {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Get UserId", func(t *testing.T) {
		data, err := models.GetUserDetails(1)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Email == "" {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Login user", func(t *testing.T) {
		username := "1234567890"
		password := "123456"
		data, err := models.LoginUser(username, password)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Email == "" {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Update user", func(t *testing.T) {
		var user = models.UpdateUserRequest{
			Id:          1,
			FirstName:   "Devendra",
			LastName:    "pohekar",
			Email:       "rideshnath.siliconithub@gmail.com",
			PhoneNumber: "1234567890",
			Role:        "deve",
			Country:     5,
			Age:         23,
		}
		data, err := models.UpdateUser(user)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("serch user", func(t *testing.T) {
		search := "rid"
		data, err := models.SearchUser(search)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if len(data) == 0 {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})

}

func TestGetcolumn(t *testing.T) {
	yourStructType := reflect.TypeOf(models.HomeSetting{})
	err := models.SynchronizeModelWithDB("home_setting", yourStructType)
	// datas, _ := json.Marshal(data)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDemo(t *testing.T) {
	req := httplib.Get("https://dummy.restapiexample.com/api/v1/employees")
	str, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}
