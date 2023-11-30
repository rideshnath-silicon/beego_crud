package test

import (
	"CarCrudv2/models"
	"testing"
)

func TestHomeSettingModels(t *testing.T) {
	t.Run("Get Home Settings", func(t *testing.T) {
		data, err := models.GetHomeSetting(1)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Insert New home settings", func(t *testing.T) {
		homeSetting := models.InserNewHomeSettingRequest{
			Section: "Left",
			Type:    "title",
			Value:   "HomeSetting",
		}
		data, err := models.InsertNewHomeSetting(homeSetting)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})

	t.Run("Insert New home settings", func(t *testing.T) {
		homeSetting := models.UpdateHomeSetingRequest{
			Id:      1,
			Section: "Left",
			Type:    "title",
			Value:   "HomeSetting",
		}
		data, err := models.UpdateHomeSeting(homeSetting)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("Userwise Home Settings", func(t *testing.T) {
		data, err := models.UserWiseHomeseting(1)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
}
