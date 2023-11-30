package test

import (
	"CarCrudv2/controllers"
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

var Home_ctrl = controllers.HomeSettingController{}

func TestGetHomeSettings(t *testing.T) {
	endPoint := "/v1/home/"
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
	tokan := fmt.Sprintf("Bearer %s", validToken)

	var jsonStr = []byte(`{"home_seting_id":1}`)
	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	w := RunControllerRoute(endPoint, r, &Home_ctrl, tokan, "post:GetHomeSetting")
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}

func TestInsertHomeSetting(t *testing.T) {
	endPoint := "/v1/home/create"
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
	tokan := fmt.Sprintf("Bearer %s", validToken)

	var jsonStr = []byte(`{
				"section" : "right",
				"type": "title",
				"value":"Heading page"
			}`)

	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := RunControllerRoute(endPoint, r, &Home_ctrl, tokan, "post:InsertNewHomeSetting")

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}

func TestUpdateHomeSetting(t *testing.T) {
	endPoint := "/v1/home/update"
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
	tokan := fmt.Sprintf("Bearer %s", validToken)

	var jsonStr = []byte(`{
				"home_seting_id":3,
				"section" : "right",
				"type": "titdddd",
				"value":"Heading page"
			}`)

	r, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := RunControllerRoute(endPoint, r, &Home_ctrl, tokan, "put:UpdateHomeSeting")

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}
