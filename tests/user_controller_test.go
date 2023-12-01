package test

import (
	"CarCrudv2/controllers"
	"CarCrudv2/middleware"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

var user_ctrl = controllers.UserController{}

func TestLoginUser(t *testing.T) {
	t.Run("Login User", func(t *testing.T) {
		endPoint := "/v1/login/"
		ctrl := &middleware.MiddlewareController{}
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com",
			"password": "123456"
		}`)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := beego.NewControllerRegister()
		router.Add(endPoint, ctrl, beego.WithRouterMethods(ctrl, "post:Login"))
		router.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestGetusers(t *testing.T) {
	t.Run("Get All User", func(t *testing.T) {
		endPoint := "/v1/users/users"
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDE0MzgzMzB9.e32LZsFTAehoXxR9v__btd9IiEoCMhMDDLV1tLXNyaU"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		r, err := http.NewRequest("GET", endPoint, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "get:GetAllUser")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestNewUser(t *testing.T) {
	t.Run("new user", func(t *testing.T) {
		endPoint := "/v1/user/register"
		var jsonStr = []byte(`{"first_name":"Dwarkesh","last_name":"patel","email":"dwarkeshpatel@gmail.com","country":"India","role":"Developer","age":30,"phone_number":"1234567890","password":"123456"}`)
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDE0MDk0MDR9.b_Faa-LR74VJubKMoO6NyCh2RPNLKxy_dNT_vs3iFTY"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:PostRegisterNewUser")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Update user", func(t *testing.T) {
		endPoint := "/v1/user/update"
		var jsonStr = []byte(`{"user_id":9,"first_name":"Dwarkesh","last_name":"patel","email":"dwarkeshpatel@gmail.com","country":"India","role":"Developer","age":30,"phone_number":"1234567890"}`)
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)
		r, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "put:UpdateUser")

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

// Please run Indidual t.run in this function to test the send otp and verify otp
func TestVerifyEmail(t *testing.T) {
	t.Run("Send email otp", func(t *testing.T) {
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)
		endPoint := "/v1/user/verify_email"
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com"
		}`)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:VerifyUserEmail")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})

	t.Run("Verify email otp", func(t *testing.T) {
		endPoint := "/v1/user/verify_email_otp"

		// enter otp after getting otp from send otp test
		var jsonStr = []byte(`{
			"username": "1234567890",
			"otp":"9310" 
		}`)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := beego.NewControllerRegister()
		router.Add(endPoint, &user_ctrl, beego.WithRouterMethods(&user_ctrl, "post:VerifyEmailOTP"))
		router.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestForgotPassword(t *testing.T) {
	t.Run("Send  otp", func(t *testing.T) {
		endPoint := "/v1/user/forgot_pass"
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com"
		}`)
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		r.Header.Set("Content-Type", "application/json")
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:SendOtp")
		router := beego.NewControllerRegister()
		router.Add(endPoint, &user_ctrl, beego.WithRouterMethods(&user_ctrl, "post:SendOtp"))
		router.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
	t.Run("Verify email otp", func(t *testing.T) {
		endPoint := "/v1/user/reset_pass_otp"
		// enter otp after getting otp from send otp test
		var jsonStr = []byte(`{
			"email" : "rideshnath.siliconithub@gmail.com",
			"otp":"0703",
		"new_password":"123456"
		}`)
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:VerifyOtpResetpassword")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})

}

func TestSerchUser(t *testing.T) {
	t.Run("Search users", func(t *testing.T) {
		endPoint := "/v1/user/search"
		var jsonStr = []byte(`{
			"search" : "rid"
		}`)

		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:SearchUser")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}
