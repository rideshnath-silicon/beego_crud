package controllers

import (
	"CarCrudv2/helpers"
	"CarCrudv2/helpers/common"
	"CarCrudv2/models"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/httplib"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Prepare() {
	// Set the language for the current request
	// langs := []string{"en-US", "hi-IN"} // List of supported languages
	// for _, lang := range langs {
	// 	if err := i18n.SetMessage(lang, "conf/locale/locale_"+lang+".ini"); err != nil {
	// 		// logger.Error("Fail to set message file:", err)
	// 		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang,err.Error()), http.StatusBadRequest, 1001)
	// 		return
	// 	}
	// }
	var lang string
	lang = c.Ctx.Input.Query("lang")
	lang = helpers.CorrectlanguageCode(lang)
	if len(lang) == 0 {
		lang = c.Ctx.GetCookie("lang")
		lang = helpers.CorrectlanguageCode(lang)
		if len(lang) != 0 {
			c.Data["Lang"] = lang
		} else {
			lang = c.Ctx.Input.Header("Accept-Language")
			if len(lang) > 4 {
				lang = lang[:5] // Only compare first 5 letters.
				lang = helpers.CorrectlanguageCode(lang)
				if lang == "en-US" || lang == "hi-IN" {
					c.Data["Lang"] = lang
				} else {
					lang = "en-US"
					c.Data["Lang"] = lang
				}
			}
		}
	} else {
		c.Data["Lang"] = lang
	}
	c.Ctx.SetCookie("lang", lang)
}

var key, _ = beego.AppConfig.String("JWT_SEC_KEY")
var jwtKey = []byte(key)

// Login ...
// @Title login User
// @Desciption login
// @Param body body models.UserLoginRequest true "login User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Accept-Languages   header  string  false  "Bearer YourAccessToken"
// @Success 201 {object} object
// @Failure 403
// @router /login [post]
func (c *UserController) Login() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var user models.UserLoginRequest
	err := helpers.RequestBody(c.Ctx, &user)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	HashPassWord, err := models.GetUserByEmail(user.Email)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	if HashPassWord.Password == "" {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "LOGIN_ERROR"), http.StatusBadRequest, 1001)
		return
	}

	err = helpers.VerifyHashedData(HashPassWord.Password, user.Password)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	userData, _ := models.LoginUser(user.Email, HashPassWord.Password)
	if userData.Email == "" && userData.FirstName == "" {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "UNAUTHORIZED"), http.StatusBadRequest, 5001)
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
// @Param lang query string false "use en-US or hi-IN"
// @Sucess 200 {string} string
// @Failure 403
// @router /logout [get]
func (c *UserController) Logout() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	c.DestroySession()
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, "LOGOUT"), http.StatusOK, 0000)
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param lang query string false "use en-US or hi-IN"
// @Param page query int false "For pagination"
// @Param limit query int false "per page user"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.Users
// @Failure 403
// @router /secure/users/ [get]
func (c *UserController) GetAllUser() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	users, err := models.GetAllUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	statIndex, endIndex, pagination, err := helpers.Pagination(page, limit, len(users))
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	output := users[statIndex:endIndex]
	result := map[string]interface{}{"pagination": pagination, "Data": output}
	helpers.ApiSuccess(c.Ctx, result, http.StatusOK, 1000)
	// helpers.ApiSuccess(c.Ctx, users, http.StatusOK, 1000)
}

// PostRegisterNewUser ...
// @Title Insert New User
// @Desciption new users
// @Param lang query string false "use en-US or hi-IN"
// @Param body body models.NewUserRequest true "Insert New User"
// @Success 201 {object} models.Users
// @Failure 403
// @router /register [post]
func (c *UserController) PostRegisterNewUser() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.NewUserRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "FORM_BODY"), http.StatusBadRequest, 1001)
		return
	}
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	valid, errString := models.Validate(lang, &bodyData)
	if !valid {
		helpers.ApiFailure(c.Ctx, errString, http.StatusBadRequest, 1001)
		return
	}
	err = bodyData.NewUserValidate()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, _ := models.GetUserByEmail(bodyData.Email)
	if data.Email == bodyData.Email {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "DUPLICATE_EMAIL"), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.InsertNewUser(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	go common.VerifyEmail(output.Email, output.FirstName)
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1002)
}

// UpdateUser .
// @Title update User
// @Desciption update users
// @Param body body models.UpdateUserRequest true "update New User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/update [put]
func (c *UserController) UpdateUser() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.UpdateUserRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	valid, errString := models.Validate(lang, &bodyData)
	if !valid {
		helpers.ApiFailure(c.Ctx, errString, http.StatusBadRequest, 1001)
		return
	}
	err = bodyData.UdateUserValidate()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.GetUserDetails(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.Email != data.Email {
		res, _ := models.GetUserByEmail(bodyData.Email)
		if res.Email == bodyData.Email {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "DUPLICATE_EMAIL"), http.StatusBadRequest, 10001)
			return
		}
	}
	output, err := models.UpdateUser(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, output.(string)), http.StatusOK, 1003)
}

// ResetPassword ...
// @Title Reset password
// @Desciption Reset password
// @Param body body models.ResetUserPassword true "reset password"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/reset_pass [post]
func (c *UserController) ResetPassword() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	claims := helpers.GetUserDataFromTokan(c.Ctx)
	id := claims["User_id"].(float64)
	output, err := models.GetUserDetails(id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	var bodyData models.ResetUserPassword
	err = helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	err = helpers.VerifyHashedData(output.Password, bodyData.CurrentPass)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.ConfirmPass != bodyData.NewPass {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "Not_match"), http.StatusBadRequest, 1001)
		return
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, uppass.(string)), http.StatusOK, 1003)
}

// SendOtp ...
// @Title forgot password
// @Desciption forgot password
// @Param body body models.SendOtpData true "forgot password this is send otp on mobile and email"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/forgot_pass [post]
func (c *UserController) ForgotPassword() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.SendOtpData
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Username)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	_, err = helpers.TwilioSendOTP(output.PhoneNumber)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	res, err := common.VerifyEmail(output.Email, output.FirstName)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, res), http.StatusOK, 1000)
	go func() {
		newOtp := helpers.GenerateOtp()
		models.UpdateColumnOTP(output.Id, newOtp)
	}()
}

// VerifyOtpResetpassword ...
// @Title verify otp
// @Desciption otp verification for forgot password
// @Param body body models.ResetUserPasswordOtp true "otp verification for forgot password"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/reset_pass_otp [post]
func (c *UserController) VerifyOtpResetpassword() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.ResetUserPasswordOtp
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Email)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	err = helpers.TwilioVerifyOTP(output.PhoneNumber, bodyData.Otp)
	if err != nil {
		data, err := models.VerifyEmailOTP(bodyData.Email, bodyData.Otp)
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
			return
		}
		if data.Otp != bodyData.Otp {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "OTP_NOT_VERIFY"), http.StatusBadRequest, 5001)
		}
		err = models.UpdateVerified(data.Id)
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
			return
		}
		uppass, err := models.ResetPassword(bodyData.NewPass, float64(output.Id))
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
			return
		}
		helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, uppass.(string)), http.StatusOK, 1003)
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, float64(output.Id))
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, uppass.(string)), http.StatusOK, 1003)

}

// VerifyUserEmail ...
// @Title verify email
// @Desciption Verify email
// @Param body body models.SendOtpData true "verify email"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verify_email [post]
func (c *UserController) VerifyUserEmail() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.SendOtpData
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Username)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	res, err := common.VerifyEmail(output.Email, output.FirstName)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, res), http.StatusOK, 1000)
	go func() {
		<-time.After(5 * time.Minute)
		newOtp := helpers.GenerateOtp()
		models.UpdateColumnOTP(output.Id, newOtp)
	}()
}

// VerifyEmailOTP ...
// @Title verify otp for email
// @Desciption otp verification for eamil
// @Param body body models.VerifyEmailOTPRequest true "otp verification for email"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verify_email_otp [post]
func (c *UserController) VerifyEmailOTP() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.VerifyEmailOTPRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.VerifyEmailOTP(bodyData.Username, bodyData.Otp)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	if data.Otp != bodyData.Otp {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "OTP_NOT_VERIFY"), http.StatusBadRequest, 5001)
		return
	}
	err = models.UpdateVerified(data.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, "LOGIN"), http.StatusOK, 5000)
}

func (c *UserController) GetCountryWiseCountUser() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	res, err := models.GetCountryWiseCountUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, res, http.StatusOK, 1000)
}

// GetVerifiedUsers ...
// @Title verifid users
// @Desciption Get all verified user
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verified_user [get]
func (c *UserController) GetVerifiedUsers() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	user, err := models.GetVerifiedUsers()
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}

	var output []models.UserDetailsRequest
	for _, user := range user {
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Age: user.Age, CreatedAt: helpers.GetFormatedDate(user.CreatedAt, "yyyy-mm-dd")}
		output = append(output, userDetails)
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
}

// SearchUser ...
// @Title Search User
// @Desciption SearchUser
// @Param body body models.SearchRequest true "otp verification for email"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/search [post]
func (c *UserController) SearchUser() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.SearchRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	user, err := models.SearchUser(bodyData.Search)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	var output []models.UserDetailsRequest
	for _, user := range user {
		userDetails := models.UserDetailsRequest{Id: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Age: user.Age, Country: user.Country, CreatedAt: helpers.GetFormatedDate(user.CreatedAt, "yyyy-mm-dd")}
		output = append(output, userDetails)
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1000)
}

// @Title export Data
// Description export in xlsx
// @Param lang query string false "use en-US or hi-IN"
// @Success 201 {string} string
// @Failure 403
// @router /export/xlxs [get]
func (c *UserController) FetchAndExportToXLS() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	users, err := models.GetAllUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	// Fetch data from the database
	// Create an Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		c.Ctx.WriteString("Error creating sheet: " + err.Error())
		return
	}
	columnNames := []string{"user_id", "name", "last_name", "user_email", "user_age", "country_name"}

	// Add the header row to the sheet
	headerRow := sheet.AddRow()
	for _, columnName := range columnNames {
		cell := headerRow.AddCell()
		cell.SetValue(strings.ToUpper(columnName))
	}
	for _, rowParams := range users {
		dataRow := sheet.AddRow()
		for _, columnName := range columnNames {
			cell := dataRow.AddCell()
			cellValue, ok := rowParams[columnName]
			if !ok {
				fmt.Printf("Column '%s' not found in row data\n", columnName)
				return
			}
			cell.SetValue(cellValue)
		}
	}
	err = os.MkdirAll("./exported/XLS", os.ModePerm)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	// Save the Excel file
	filePath := "./exported/XLS/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "users.xlsx"
	err = file.Save(filePath)
	if err != nil {
		c.Ctx.WriteString("Error saving Excel file: " + err.Error())
		return
	}

	c.Ctx.WriteString("Data exported to Excel successfully.")
}

// @Title export Data
// Description export in xlsx
// @Param lang query string false "use en-US or hi-IN"
// @Success 201 {string} string
// @Failure 403
// @router /export/pdf/ [get]
func (c *UserController) FetchAndExportToPDF() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	users, err := models.GetAllUser()
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	// Fetch data from the database
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 12)
	header := []string{"user_id", "name", "last_name", "user_email", "user_age", "country_name"}
	for _, col := range header {
		pdf.CellFormat(40, 10, strings.ToUpper(col), "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	// Add data to the PDF
	for _, row := range users {
		for _, col := range header {
			pdf.CellFormat(40, 10, fmt.Sprintf("%v", row[col]), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
	// Save the PDF
	err = os.MkdirAll("./exported/PDF", os.ModePerm)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "ERROR_DIRECTORY"), http.StatusBadRequest, 1001)
		return
	}
	filePath := "./exported/PDF/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "users.pdf"
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		c.Ctx.WriteString("Error generating PDF: " + err.Error())
		return
	}
	c.Ctx.WriteString("Data exported to pdf successfully.")
}

// ImportData
// @Title export Data
// Description export in xlsx
// @Param lang query string false "use en-US or hi-IN"
// @Param file formData file true "File to be uploaded"
// @Success 201 {string} string
// @Failure 403
// @router /import/xlsx [post]
func (c *UserController) ImportData() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	_, fileheader, err := c.GetFile("file")
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}

	// Save the uploaded file
	filedName := "file"
	uploadDir := "./uploads/imported/XLS/"
	filepaths, err := helpers.UploadFile(c.Controller, filedName, fileheader, uploadDir)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	// Import data from the file into the database
	err = common.ImportData(filepaths)
	if err != nil {
		c.Ctx.WriteString("Error importing data: " + err.Error())
		return
	}

	c.Ctx.WriteString("Data imported successfully.")
}

// @Title delete user
// @Description delete user
// @Param id path int true "user id to delete recode"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {string} string
// @Failure 403
// @router /secure/delete/:id([0-9]+) [delete]
func (c *UserController) DeleteUser() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	idString := c.Ctx.Input.Params()
	id, err := strconv.Atoi(idString["1"])
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.DeleteUser(uint(id))
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, data), http.StatusOK, 1000)
}

// @Title Insert Update
// @Description Insert update both
// @Param body body models.NewUserRequest true "update New User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/inup [post]
func (c *UserController) InsertUpdate() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.NewUserRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	fmt.Printf("%T", bodyData)
	Tokan := c.Ctx.Input.Header("Authorization")
	data := helpers.GetUserDataFromTokan(c.Ctx)
	if bodyData.Email == data["Email"] {
		req := httplib.Post("http://localhost:8080/v1/user/secure/update")
		req.Header("Authorization", Tokan)
		req.Header("Content-Type", "application/json")
		req.Body(bodyData)
		str, err := req.String()
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "Not_Found"), http.StatusBadRequest, 1001)
			return
		}
		c.Ctx.WriteString(str)
	} else {
		req := httplib.Post("http://localhost:8080/v1/user/register")
		req.Header("Authorization", Tokan)
		req.Header("Content-Type", "application/json")
		req.Body(bodyData)
		str, err := req.String()
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "Not_Found"), http.StatusBadRequest, 1001)
			return
		}
		c.Ctx.WriteString(str)
	}
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<Create Message Code >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// @Title Create message code
// @Description generate new message
// @Param Code 		formData string true "Write Message Key without space you can you ( _ ) place of space"
// @Param English 	formData string true "Write Message in English"
// @Param Hindi 	formData string true "Write Message in Hindi"
// @Param lang query string false "use en-US or hi-IN"
// @Success 200 {string} string
// @Failure 403
// @router /newMessage [post]
func (c *UserController) CreateMessageCode() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var message models.NewMessageCode
	if err := c.ParseForm(&message); err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	res, err := models.InsertNewMessage(message)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, res, http.StatusOK, 1002)
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
	lang := c.Ctx.Input.GetData("Lang").(string)
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
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
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
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	c.Ctx.WriteString(fmt.Sprintf("\nValue: %v,sesssion value: %v\n", str, userData))
	req = httplib.Post("http://localhost:8080/v1/user/secure/search")
	req.Header("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InJpZGVzaG5hdGguc2lsaWNvbml0aHViQGdtYWlsLmNvbSIsIklEIjoxLCJleHAiOjE3MDE3NzIwMDR9.5FHeDnH39qG_kUfM-9arLEZP_F9FMJAqIXi_8akQRFI")
	req.Header("Content-Type", "application/json")
	req.Body(`{"search":"de"}`)
	str, err = req.String()
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	c.Ctx.WriteString(fmt.Sprintf("\nValue: %v,sesssion value: %v\n", str, userData))

	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<httplib>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< checking end >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
