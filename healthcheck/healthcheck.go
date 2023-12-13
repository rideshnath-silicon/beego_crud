package healthcheck

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type DatabaseCheck struct {
}

func (dc *DatabaseCheck) isConnected() bool {
	connection, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		return false
	}
	err = orm.RegisterDriver("postgres", orm.DRPostgres)
	if err != nil {
		return false
	}
	err = orm.RegisterDataBase("default1", "postgres", connection)
	if err != nil {
		return false
	}
	return true
}

func (dc *DatabaseCheck) Check() error {
	if dc.isConnected() {
		return nil
	} else {
		return errors.New("can't connect to the database")
	}
}
