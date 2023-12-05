package controllers

import (
	"CarCrudv2/helpers"
	"CarCrudv2/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/httplib"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
)

type UserController struct {
	beego.Controller
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<checking modules>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
var filepath = &cache.FileCache{
	CachePath:  "cache",
	FileSuffix: ".cache",
}

// Checkmoduls
// @Title check modules
// @Description sesseions and all module to be checked in this
// @Success 200 {object} object
// @Failure 403
// @router /modulecheck [get]
func (c *UserController) Checkmodul() {
	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<Session and cache>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	c.SetSession("userData", "rideshNathji")
	// cache.Cache.Put("myKey", "Hello, World!", 3600)
	c.Ctx.Input.CruSession.Set(context.Background(), "myKey", "Hello, World!")

	cache.Cache.Put(filepath, context.Background(), "key", "ridesh", 60*time.Second)

	userData := c.GetSession("userData")
	if userData == nil {
		helpers.ApiFailure(c.Ctx, "Session is destroyed please login again", http.StatusBadRequest, 1001)
		return
	}
	c.Ctx.WriteString(fmt.Sprintf("sesssion Value: %v", userData))
	value, err := cache.Cache.Get(filepath, context.Background(), "key")
	if value != nil {
		// Value found in the cache
		c.Ctx.WriteString(fmt.Sprintf("\nCached Value: %v", value))
	}
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	session := c.StartSession()
	// Retrieve the session ID
	sessionID := session.SessionID(context.Background())
	// Use the session ID as neede

	c.Ctx.WriteString("\nSession ID: " + sessionID)
	// c.DestroySession()
	session.Delete(context.Background(), "userData")
	userData = c.GetSession("userData")

	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<Session and cache>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<httplib>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	req := httplib.Get("http://localhost:8080/v1/user/secure/users")
	req.Header("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDE3NzIwMDR9.5FHeDnH39qG_kUfM-9arLEZP_F9FMJAqIXi_8akQRFI")
	str, err := req.String()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	c.Ctx.WriteString(fmt.Sprintf("\nValue: %v,sesssion value: %v\n", str, userData))
	req = httplib.Post("http://localhost:8080/v1/user/secure/search")
	req.Header("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDE3NzIwMDR9.5FHeDnH39qG_kUfM-9arLEZP_F9FMJAqIXi_8akQRFI")
	req.Header("Content-Type", "application/json")
	req.Body(`{"search":"de"}`)
	str, err = req.String()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	c.Ctx.WriteString(fmt.Sprintf("\nValue: %v,sesssion value: %v\n", str, userData))

	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<httplib>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< checking end >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

var key, _ = beego.AppConfig.String("JWT_SEC_KEY")
var jwtKey = []byte(key)

// Login ...
// @Title login User
// @Desciption login
// @Param body body models.UserLoginRequest true "login User"
// @Success 201 {object} object
// @Failure 403
// @router /login [post]
func (c *UserController) Login() {
	var user models.UserLoginRequest
	err := helpers.RequestBody(c.Ctx, &user)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	HashPassWord, err := models.GetUserByEmail(user.Email)

	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	if HashPassWord.Password == "" {
		helpers.ApiFailure(c.Ctx, "please enter valid Username Or Passwordsss ", http.StatusBadRequest, 1001)
		return
	}

	err = helpers.VerifyHashedData(HashPassWord.Password, user.Password)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	userData, _ := models.LoginUser(user.Email, HashPassWord.Password)
	if userData.Email == "" && userData.FirstName == "" {
		helpers.ApiFailure(c.Ctx, "Unauthorized User", http.StatusBadRequest, 5001)
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.JwtClaim{Email: userData.Email, ID: int(userData.Id), StandardClaims: jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 5001)
		return
	}
	data := map[string]interface{}{"User_Data": token.Claims, "Tokan": tokenString}
	helpers.ApiSuccess(c.Ctx, data, http.StatusOK, 5000)
}

// Logout ...
// @Title logout user
// @Description logout
// @Sucess 200 {string} string
// @Failure 403
// @router /logout [get]
func (c *UserController) Logout() {
	c.DestroySession()
	helpers.ApiSuccess(c.Ctx, "logout successfully ", http.StatusOK, 0000)
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.Users
// @Failure 403
// @router /secure/users [get]
func (c *UserController) GetAllUser() {
	user, err := models.GetAllUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, user, http.StatusOK, 1000)
}

// PostRegisterNewUser ...
// @Title Insert New User
// @Desciption new users
// @Param body body models.NewUserRequest true "Insert New User"
// @Success 201 {object} models.Users
// @Failure 403
// @router /register [post]
func (c *UserController) PostRegisterNewUser() {
	var bodyData models.NewUserRequest
	if err := c.ParseForm(&bodyData); err != nil {
		c.Ctx.WriteString("Error while parsing form data: " + err.Error())
		return
	}
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, _ := models.GetUserByEmail(bodyData.Email)

	if data.Email == bodyData.Email {
		helpers.ApiFailure(c.Ctx, "Email already used by another account please try with new email", http.StatusBadRequest, 10001)
		return
	}
	output, err := models.InsertNewUser(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1002)
}

// UpdateUser ...
// @Title update User
// @Desciption update users
// @Param body body models.UpdateUserRequest true "update New User"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/update [put]
func (c *UserController) UpdateUser() {
	var bodyData models.UpdateUserRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}

	data, err := models.GetUserDetails(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.Email != data.Email {
		res, _ := models.GetUserByEmail(bodyData.Email)
		if res.Email == bodyData.Email {
			helpers.ApiFailure(c.Ctx, "Email already used by another account please try with new email", http.StatusBadRequest, 10001)
			return
		}
	}
	output, err := models.UpdateUser(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1003)
}

// ResetPassword ...
// @Title Reset password
// @Desciption Reset password
// @Param body body models.ResetUserPassword true "reset password"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/reset_pass [post]
func (c *UserController) ResetPassword() {
	claims := helpers.GetUserDataFromTokan(c.Ctx)
	id := claims["User_id"].(float64)
	output, err := models.GetUserDetails(id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	var bodyData models.ResetUserPassword
	err = helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	err = helpers.VerifyHashedData(output.Password, bodyData.CurrentPass)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.ConfirmPass != bodyData.NewPass {
		helpers.ApiFailure(c.Ctx, "Please match new password and confirm password", http.StatusBadRequest, 1001)
		return
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, uppass, http.StatusOK, 1003)
}

// SendOtp ...
// @Title forgot password
// @Desciption forgot password
// @Param body body models.SendOtpData true "forgot password this is send otp on mobile and email"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/forgot_pass [post]
func (c *UserController) SendOtp() {
	var bodyData models.SendOtpData
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Username)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	_, err = helpers.TwilioSendOTP(output.PhoneNumber)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	otp, err := helpers.SendMailOTp(output.Email, output.FirstName)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	Response, err := models.UpadteOtpForEmail(output.Id, otp)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, Response, http.StatusOK, 1000)
	go func() {
		newOtp := helpers.GenerateOtp()
		models.UpdateColumnOTP(output.Id, newOtp)
	}()
}

// VerifyOtpResetpassword ...
// @Title verify otp
// @Desciption otp verification for forgot password
// @Param body body models.ResetUserPasswordOtp true "otp verification for forgot password"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/reset_pass_otp [post]
func (c *UserController) VerifyOtpResetpassword() {
	var bodyData models.ResetUserPasswordOtp
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Email)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	err = helpers.TwilioVerifyOTP(output.PhoneNumber, bodyData.Otp)
	if err != nil {
		data, err := models.VerifyEmailOTP(bodyData.Email, bodyData.Otp)
		if err != nil {
			helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
			return
		}
		if data.Otp != bodyData.Otp {
			helpers.ApiFailure(c.Ctx, "Please Eenter Valid otp", http.StatusBadRequest, 5001)
		}
		err = models.UpdateVerified(data.Id)
		if err != nil {
			helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
			return
		}
		uppass, err := models.ResetPassword(bodyData.NewPass, float64(output.Id))
		if err != nil {
			helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
			return
		}
		helpers.ApiSuccess(c.Ctx, uppass, http.StatusOK, 1003)
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, float64(output.Id))
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, uppass, http.StatusOK, 1003)
}

// VerifyUserEmail ...
// @Title verify email
// @Desciption Verify email
// @Param body body models.SendOtpData true "verify email"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verify_email [post]
func (c *UserController) VerifyUserEmail() {
	var bodyData models.SendOtpData
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Username)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	otp, err := helpers.SendMailOTp(output.Email, output.FirstName)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	res, err := models.UpadteOtpForEmail(output.Id, otp)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, res, http.StatusOK, 1000)
	go func() {
		newOtp := helpers.GenerateOtp()
		models.UpdateColumnOTP(output.Id, newOtp)
	}()
}

// VerifyEmailOTP ...
// @Title verify otp for email
// @Desciption otp verification for eamil
// @Param body body models.VerifyEmailOTPRequest true "otp verification for email"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verify_email_otp [post]
func (c *UserController) VerifyEmailOTP() {
	var bodyData models.VerifyEmailOTPRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.VerifyEmailOTP(bodyData.Username, bodyData.Otp)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	if data.Otp != bodyData.Otp {
		helpers.ApiFailure(c.Ctx, "Please Eenter Valid otp", http.StatusBadRequest, 5001)
	}
	err = models.UpdateVerified(data.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, "Your Account is Successfully Verified", http.StatusOK, 5000)
}

func (c *UserController) GetCountryWiseCountUser() {
	res, err := models.GetCountryWiseCountUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, res, http.StatusOK, 1000)
}

// GetVerifiedUsers ...
// @Title verifid users
// @Desciption Get all verified user
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verified_user [get]
func (c *UserController) GetVerifiedUsers() {
	user, err := models.GetVerifiedUsers()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}

	var output []models.UserDetailsRequest
	for _, user := range user {
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Age: user.Age}
		output = append(output, userDetails)
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
}

// SearchUser ...
// @Title Search User
// @Desciption SearchUser
// @Param body body models.SearchRequest true "otp verification for email"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/search [post]
func (c *UserController) SearchUser() {
	var bodyData models.SearchRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	user, err := models.SearchUser(bodyData.Search)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	var output []models.UserDetailsRequest
	for _, user := range user {
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Age: user.Age}
		output = append(output, userDetails)
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
}
