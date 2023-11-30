package main_test

import (
	"CarCrudv2/models"
	"log"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

func Init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=root dbname=mydb sslmode=disable")
	log.Print("<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>")
	// orm.RegisterModel(models.Users{})
	orm.RunSyncdb("default", false, true)
}

func TestInit(t *testing.T) {
	Init()
	t.Run("demo", func(t *testing.T) {
		data, err := models.GetAllUser()
		if err != nil {
			t.Errorf("Error ;- %s", err.Error())
		}
		log.Print(data)
	})
}
