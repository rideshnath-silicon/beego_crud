// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact rideshnath.siliconithub@gmail.com
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
		beego.NSNamespace("/user",
			beego.NSInclude(&controllers.UserController{}),
			beego.NSRouter("/register", &controllers.UserController{}, "Post:PostRegisterNewUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/logout", &controllers.UserController{}, "get:Logout"),
			beego.NSRouter("/modulecheck", &controllers.UserController{}, "get:Checkmodul"),
			beego.NSNamespace("/secure",
				beego.NSBefore(middleware.JWTMiddleware),
				beego.NSRouter("/forgot_pass", &controllers.UserController{}, "post:SendOtp"),
				beego.NSRouter("/reset_pass_otp", &controllers.UserController{}, "post:VerifyOtpResetpassword"),
				beego.NSRouter("/users", &controllers.UserController{}, "Get:GetAllUser"),
				beego.NSRouter("/verify_email", &controllers.UserController{}, "post:VerifyUserEmail"),
				beego.NSRouter("/verify_email_otp", &controllers.UserController{}, "post:VerifyEmailOTP"),
				beego.NSRouter("/update", &controllers.UserController{}, "put:UpdateUser"),
				beego.NSRouter("/reset_pass", &controllers.UserController{}, "post:ResetPassword"),
				beego.NSRouter("/contries", &controllers.UserController{}, "get:GetCountryWiseCountUser"),
				beego.NSRouter("/verified_user", &controllers.UserController{}, "get:GetVerifiedUsers"),
				beego.NSRouter("/search", &controllers.UserController{}, "post:SearchUser"),
			),
		),
		beego.NSNamespace("/car",
			beego.NSInclude(&controllers.CarController{}),
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSRouter("/", &controllers.CarController{}, "post:GetSingleCar"),
			beego.NSRouter("/cars", &controllers.CarController{}, "get:GetAllCars"),
			beego.NSRouter("/search", &controllers.CarController{}, "post:GetCarUsingSearch"),
			beego.NSRouter("/create", &controllers.CarController{}, "post:AddNewCar"),
			beego.NSRouter("/update", &controllers.CarController{}, "put:UpdateCar"),
			beego.NSRouter("/delete", &controllers.CarController{}, "delete:DeleteCar"),
		),
		beego.NSNamespace("/home",
			beego.NSInclude(&controllers.HomeSettingController{}),
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSRouter("/", &controllers.HomeSettingController{}, "post:GetHomeSetting"),
			beego.NSRouter("/userwise", &controllers.HomeSettingController{}, "post:GetUserWiseHome"),
			beego.NSRouter("/create", &controllers.HomeSettingController{}, "post:InsertNewHomeSetting"),
			beego.NSRouter("/update", &controllers.HomeSettingController{}, "put:UpdateHomeSeting"),
		),
	)

	beego.AddNamespace(ns)

	ns1 := beego.NewNamespace("/v2",
		beego.NSNamespace("/user",
			beego.NSInclude(&controllers.UserController{}),
			beego.NSRouter("/sign_up", &controllers.UserController{}, "Post:PostRegisterNewUser"),
			beego.NSRouter("/sign_in", &controllers.UserController{}, "post:Login"),
		),
	)
	beego.AddNamespace(ns1)
}

// func init() {
// 	ns := beego.NewNamespace("/v1",
// 		beego.NSNamespace("/user",
// 			beego.NSInclude(&controllers.UserController{}),
// 			beego.NSAutoRouter(&controllers.UserController{}),
// 		),
// 	)
// 	beego.AddNamespace(ns)
// }
