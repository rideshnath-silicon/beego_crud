package main

import (
	_ "CarCrudv2/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func init() {
	connection, _ := beego.AppConfig.String("sqlconn")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", connection)
	orm.RunSyncdb("default", false, true)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	initLogs()
	beego.Run()
}

func initLogs() {
	// from, _ := config.String("EMAIL")
	// password, _ := config.String("PASSWORD")
	// emailConfig := `{
	// 	"username": "` + from + `",
	// 	"password": "` + password + `",
	// 	"host":     "smtp.gmail.com:587",
	// 	"sendTos":  ["rideshnath.siliconithub@gmail.com"]
	// }`

	// // Initialize the logs with the email adapter
	// err := logs.SetLogger(logs.AdapterMail, emailConfig)
	// if err != nil {
	// 	panic("Error setting email logger: " + err.Error())
	// }
	err := logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	if err != nil {
		panic("Error setting logger: " + err.Error())
	}

	logs.Debug("Debug message")
	logs.Info("Info message")
	logs.Warn("Warning message")
	logs.Error("Error message")
	logs.Critical("Critical message")
}
