package healthcheck

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type DatabaseCheck struct {
}

func (dc *DatabaseCheck) isConnected() bool {
	o := orm.NewOrm()
	_, err := o.Raw("select 1").Exec()
	return err == nil
}

func (dc *DatabaseCheck) Check() error {
	if dc.isConnected() {
		return nil
	} else {
		return errors.New("can't connect to the database")
	}
}
