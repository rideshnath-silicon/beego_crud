package test

import (
	"CarCrudv2/middleware"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=root dbname=mydb sslmode=disable")
	// orm.RunSyncdb("default", false, true)
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func RunControllerRoute(endPoint string, r *http.Request, ctrl beego.ControllerInterface, tokan string, methodFuction string) *httptest.ResponseRecorder {
	r.Header.Set("Authorization", tokan)
	w := httptest.NewRecorder()
	router := beego.NewControllerRegister()
	router.InsertFilter(endPoint, beego.BeforeRouter, middleware.JWTMiddleware, beego.WithCaseSensitive(false))
	router.Add(endPoint, ctrl, beego.WithRouterMethods(ctrl, methodFuction))
	router.ServeHTTP(w, r)
	return w
}

func TruncateTable(tableName string) {
	o := orm.NewOrm()
	_, err := o.Raw("TRUNCATE TABLE " + tableName).Exec()

	if err != nil {
		fmt.Println("Failed to truncate table:", err)
		return
	}
	orm.NewOrm().Raw(`SELECT setval('"` + tableName + `_id_seq"', 1, false)`).Exec()
}
