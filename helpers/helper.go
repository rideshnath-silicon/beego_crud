package helpers

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
	"golang.org/x/crypto/bcrypt"
)

func HashData(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
func VerifyHashedData(hashedString string, dataString string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(dataString))
	if err != nil {
		return errors.New("HASHED_ERROR")
	}
	return nil
}

func GetUserDataFromTokan(c *context.Context) map[string]interface{} {
	userClaims := c.Input.GetData("user")
	userID := userClaims.(jwt.MapClaims)["ID"]
	userEmail := userClaims.(jwt.MapClaims)["Email"]
	response := map[string]interface{}{"Email": userEmail, "User_id": userID}
	return response
}

func RequestBody(ctx *context.Context, Struct interface{}) error {
	bodyData := ctx.Input.RequestBody
	// fmt.Println(bodyData)
	err := json.Unmarshal(bodyData, Struct)
	if err != nil {
		return errors.New("JSON_BODY_ERROR")
	}
	return nil
}

type ApiResponse struct {
	Message string
	Success int
	Status  int
	Data    interface{}
}

func ApiSuccess(c *context.Context, data interface{}, statusCode int, messageCode int) {
	message := Messagess(messageCode)
	message = GetLangaugeMessage(c.Input.GetData("Lang").(string), message)
	Response := ApiResponse{
		Message: message,
		Success: 1,
		Status:  statusCode,
		Data:    data,
	}
	c.Output.JSON(Response, true, false)
}

func ApiFailure(c *context.Context, data interface{}, statusCode int, messageCode int) {

	message := Messagess(messageCode)
	message = GetLangaugeMessage(c.Input.GetData("Lang").(string), message)
	Response := ApiResponse{
		Message: message,
		Success: 1,
		Status:  statusCode,
		Data:    data,
	}
	c.Output.JSON(Response, true, false)
}

// otp verification from here
var Username, _ = beego.AppConfig.String("TWILIO_ACCOUNT_SID")
var Password, _ = beego.AppConfig.String("TWILIO_AUTHTOKEN")
var service_id, _ = beego.AppConfig.String("TWILIO_SERVICES_ID")
var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: Username,
	Password: Password,
})

func TwilioSendOTP(phoneNumber string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo("+91" + phoneNumber)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(service_id, params)
	if err != nil {
		return "", errors.New("TWILIO_ERROR")
	}
	return *resp.Sid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo("+91" + phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(service_id, params)
	if err != nil {
		return errors.New("OTP_NOT_VERIFY")
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}

func UploadFile(c beego.Controller, filedName string, fileheader *multipart.FileHeader, uploadPath string) (string, error) {
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		// Create the directory and any necessary parent directories
		err := os.MkdirAll("./"+uploadPath, os.ModePerm)
		if err != nil {
			return "", errors.New("ERROR_DIRECTORY")
		}
	}
	filePath := uploadPath + strconv.FormatInt(time.Now().UnixNano(), 10) + fileheader.Filename
	err := c.SaveToFile(filedName, filePath)
	if err != nil {
		return "", errors.New("FILE_ERROR")
	}
	return filePath, nil
}

func GenereateKeyForHomeSection(str1, str2 string) string {
	combinedString := str1 + " " + str2
	underscoredString := strings.ReplaceAll(combinedString, " ", "_")

	// Convert to uppercase
	uppercaseCode := strings.ToUpper(underscoredString)

	return uppercaseCode
}

func SendMailOTp(userEmail string, name string, subject string, body string) (bool, error) {
	from, _ := beego.AppConfig.String("EMAIL")
	password, _ := beego.AppConfig.String("PASSWORD")
	// from, _ := config.String("EMAIL")
	// password, _ := config.String("PASSWORD")
	to := []string{
		userEmail,
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte("Subject: " + subject + "\r\n" + mime + "\r\n" + body)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return false, errors.New("TWILIO_ERROR")
	}
	return true, nil
}
	
func GenerateOtp() string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 4)
	n, err := io.ReadAtLeast(rand.Reader, b, 4)
	if n != 4 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GetLangaugeMessage(lang string, langCode string) string {
	o := orm.NewOrm()
	sql := "SELECT langcodeid FROM language_code WHERE language_code = ?"
	// orm.Debug = true
	var langcodeid int64
	err := o.Raw(sql, langCode).QueryRow(&langcodeid)
	// fmt.Println(langcodeid)
	if err != nil {
		return err.Error()
	}
	var message string
	switch lang {
	case "en-US":
		sql := "SELECT value FROM engilsh_lang_message WHERE langcodeid = ?"
		err := o.Raw(sql, langcodeid).QueryRow(&message)
		if err != nil {
			message = err.Error()
		}
	case "hi-IN":
		sql := "SELECT value FROM hindi_lang_message WHERE langcodeid = ?"
		err := o.Raw(sql, langcodeid).QueryRow(&message)
		if err != nil {
			message = err.Error()
		}
	}
	return message
}
