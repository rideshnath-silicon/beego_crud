package main_test

import (
	"CarCrudv2/models"
	"log"
	"testing"

	// "github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

func TestInit(t *testing.T) {
	t.Run("demo", func(t *testing.T) {
		data, err := models.GetAllUser()
		if err != nil {
			t.Errorf("Error ;- %s", err.Error())
		}
		log.Print(data)
	})
}
