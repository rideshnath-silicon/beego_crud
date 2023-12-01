package controllers

import (
	"CarCrudv2/helpers"
	"CarCrudv2/models"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.Users
// @Failure 403
// @router /users [get]
func (c *UserController) GetAllUser() {
	if c.Ctx.Request.Method != "GET" {
		c.Ctx.Output.SetStatus(http.StatusMethodNotAllowed)
		c.Ctx.WriteString("Method Not Allowed")
		return
	}
	user, err := models.GetAllUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	var output []models.UserDetailsRequest
	for _, user := range user {
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Country: user.Country, Age: user.Age}
		output = append(output, userDetails)
	}

	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
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

func (c *UserController) GetVerifiedUsers() {
	user, err := models.GetVerifiedUsers()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}

	var output []models.UserDetailsRequest
	for _, user := range user {
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Country: user.Country, Age: user.Age}
		output = append(output, userDetails)
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
}

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
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Country: user.Country, Age: user.Age}
		output = append(output, userDetails)
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
}
