package test

import (
	"CarCrudv2/controllers"
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

var car_ctrl = controllers.CarController{}

func TestGetCars(t *testing.T) {
	t.Run("Get All Cars", func(t *testing.T) {
		endPoint := "/v1/car/cars"
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		r, err := http.NewRequest("GET", endPoint, nil)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "get:GetAllCars")

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})

	t.Run("Get Single Car", func(t *testing.T) {
		endPoint := "/v1/car/"
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		var jsonStr = []byte(`{"car_id":1}`)
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:GetSingleCar")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestInsertNewcar(t *testing.T) {
	endPoint := "v1/car/create"
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
	tokan := fmt.Sprintf("Bearer %s", validToken)

	var jsonStr = []byte(`{
				"car_name" : "Thar",
				"model": "4*4",
				"modified_by": "mahindara",
				"type":"sedan"
			}`)

	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:AddNewCar")

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}

func TestUpdatecar(t *testing.T) {
	endPoint := "v1/car/update"
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
	tokan := fmt.Sprintf("Bearer %s", validToken)

	var jsonStr = []byte(`{
				"car_id":1,
				"car_name" : "Thar",
				"model": "4*4",
				"modified_by": "mahindara",
				"type":"sedan"
			}`)

	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:UpdateCar")

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}

func TestSearchCar(t *testing.T) {
	t.Run("Search users", func(t *testing.T) {
		endPoint := "/v1/car/search"
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
		tokan := fmt.Sprintf("Bearer %s", validToken)

		var jsonStr = []byte(`{
			"search" : "se"
		}`)

		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:GetCarUsingSearch")

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestRemoveCar(t *testing.T) {
	endPoint := "/v1/car/delete"
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDEzNDY4ODl9.t-cNDRqPHygAu1yGHjOtpJWvhj1qaBk0WpTGHxM9Vm4"
	tokan := fmt.Sprintf("Bearer %s", validToken)

	var jsonStr = []byte(`{"car_id":1}`)
	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:DeleteCar")
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}
