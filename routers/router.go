// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"CarCrudv2/controllers"
	"CarCrudv2/middleware"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/login",
			beego.NSRouter("/", &middleware.MiddlewareController{}, "post:Login"),
		),
		beego.NSNamespace("/user",
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSRouter("/forgot_pass", &controllers.UserController{}, "post:SendOtp"),
			beego.NSRouter("/reset_pass_otp", &controllers.UserController{}, "post:VerifyOtpResetpassword"),
			beego.NSRouter("/users", &controllers.UserController{}, "Get:GetAllUser"),
			beego.NSRouter("/register", &controllers.UserController{}, "post:PostRegisterNewUser"),
			beego.NSRouter("/verify_email", &controllers.UserController{}, "post:VerifyUserEmail"),
			beego.NSRouter("/verify_email_otp", &controllers.UserController{}, "post:VerifyEmailOTP"),
			beego.NSRouter("/update", &controllers.UserController{}, "put:UpdateUser"),
			beego.NSRouter("/reset_pass", &controllers.UserController{}, "post:ResetPassword"),
			beego.NSRouter("/contries", &controllers.UserController{}, "get:GetCountryWiseCountUser"),
			beego.NSRouter("/verified_user", &controllers.UserController{}, "get:GetVerifiedUsers"),
			beego.NSRouter("/search", &controllers.UserController{}, "post:SearchUser"),
		),
		beego.NSNamespace("/car",
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSRouter("/", &controllers.CarController{}, "post:GetSingleCar"),
			beego.NSRouter("/cars", &controllers.CarController{}, "get:GetAllCars"),
			beego.NSRouter("/search", &controllers.CarController{}, "post:GetCarUsingSearch"),
			beego.NSRouter("/create", &controllers.CarController{}, "post:AddNewCar"),
			beego.NSRouter("/update", &controllers.CarController{}, "put:UpdateCar"),
			beego.NSRouter("/delete", &controllers.CarController{}, "delete:DeleteCar"),
		),
		beego.NSNamespace("/home",
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSRouter("/", &controllers.HomeSettingController{}, "post:GetHomeSetting"),
			beego.NSRouter("/userwise", &controllers.HomeSettingController{}, "post:GetUserWiseHome"),
			beego.NSRouter("/create", &controllers.HomeSettingController{}, "post:InsertNewHomeSetting"),
			beego.NSRouter("/update", &controllers.HomeSettingController{}, "put:UpdateHomeSeting"),
		),
	)
	beego.AddNamespace(ns)
}
