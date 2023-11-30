package test

import (
	"CarCrudv2/models"
	"testing"
)

func TestCarModels(t *testing.T) {
	t.Run("Get All cars", func(t *testing.T) {
		data, err := models.GetAllCars()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Get car", func(t *testing.T) {
		data, err := models.GetSingleCar(1)
		if err != nil {
			t.Errorf( err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Register New car", func(t *testing.T) {
		car := models.GetNewCarRequest{
			CarName:    "swift",
			ModifiedBy: "suzuki",
			Model:      "swift dzire",
			Type:       "sedan",
			CarImage:   "swiftImage",
		}
		data, err := models.InsertNewCar(car)
		if err != nil {
			t.Errorf( err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Register New car", func(t *testing.T) {
		car := models.UpdateCarRequest{
			Id:         2,
			CarName:    "swift",
			ModifiedBy: "suzuki",
			Model:      "swift dzire",
			Type:       "sedan",
			CarImage:   "swiftImage",
		}
		data, err := models.UpdateCar(car)
		if err != nil {
			t.Errorf( err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("serch car", func(t *testing.T) {
		search := "th"
		data, err := models.GetCarUsingSearch(search)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Delete car", func(t *testing.T) {
		data, err := models.DeleteCar(22)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
}
