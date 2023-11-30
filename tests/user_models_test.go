package test

import (
	"CarCrudv2/models"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func TestUserModels(t *testing.T) {
	t.Run("Register new user", func(t *testing.T) {
		TruncateTable("users")
		var user = models.NewUserRequest{
			FirstName:   "Devendra",
			LastName:    "pohekar",
			Email:       "rideshnath.siliconithub@gmail.com",
			PhoneNumber: "1234567890",
			Role:        "developer",
			Country:     "India",
			Age:         24,
			Password:    "123456",
		}
		data, err := models.InsertNewUser(user)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Id != 1 {
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
		t.Log(data)
	})
	t.Run("Get UserByEmail", func(t *testing.T) {
		data, err := models.GetUserByEmail("rideshnath.siliconithub@gmail.com")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Get UserId", func(t *testing.T) {
		data, err := models.GetUserDetails(1)
		if err != nil {
			t.Errorf(err.Error())
			return
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
			Country:     "India",
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
		t.Log(data)
	})

}
