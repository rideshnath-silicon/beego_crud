package common

import (
	"CarCrudv2/helpers"
	"CarCrudv2/models"
	"errors"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/beego/beego/v2/client/orm"
)

func NewCarType(input string) (models.CarType, error) {
	switch input {
	case "sedan", "hatchback", "SUV":
		return models.CarType(input), nil
	default:
		return "", errors.New("INVALID_CAR")
	}
}

func VerifyEmail(email string, name string) (string, error) {
	OTP := helpers.GenerateOtp()
	subject := "Verify your email"
	body := `<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
	<div style="margin:50px auto;width:70%;padding:20px 0">
	<div style="border-bottom:1px solid #eee">
			<a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">Hello, I am Ridesh</a>
		</div>
		<p style="font-size:1.1em">Hi, ` + name + `</p>
		<p>Thank you for Register in this app . Use the following OTP to verify your email. OTP is valid for 5 minutes</p>
		<h2 style="background: #00466a;margin: 0 auto;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;">` + OTP + `</h2>
		<p style="font-size:0.9em;">Regards,<br />Er. Ridesh Nath</p>
		<hr style="border:none;border-top:1px solid #eee" />
		<div style="float:right;padding:8px 0;color:#aaa;font-size:0.8em;line-height:1;font-weight:300">
			<p>Ridesh Nath</p>
			<p>Burhanpur M.P</p>
			<p>India</p>
		</div>
	</div>
</div>`
	o := orm.NewOrm()
	sendemail := models.EmailLogs{}
	_, err := helpers.SendMailOTp(email, name, subject, body)
	if err != nil {
		sendemail = models.EmailLogs{
			To:      email,
			Name:    name,
			Subject: subject,
			Body:    body,
			Status:  "pending",
		}
		_, err := o.Insert(&sendemail)
		if err != nil {
			return "", errors.New("DATABASE_ERROR")
		}
	}
	sendemail = models.EmailLogs{
		To:      email,
		Name:    name,
		Subject: subject,
		Body:    body,
		Status:  "success",
	}
	_, err = o.Insert(&sendemail)
	if err != nil {
		return "", err
	}
	output, err := models.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	res, err := models.UpadteOtpForEmail(output.Id, OTP)
	if err != nil {
		return "", err
	}
	return res, nil
}

func ImportData(filePath string) error {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	rows, err := xlsx.Rows("Sheet1")
	if err != nil {
		return err
	}
	isFirstRow := true
	for rows.Next() {
		var dataRow []string
		cols := rows.Columns()
		for _, col := range cols {
			dataRow = append(dataRow, col)
		}
		if isFirstRow {
			isFirstRow = false
			continue
		}
		age, err := strconv.Atoi(dataRow[5])
		if err != nil {
			return err
		}
		country, err := strconv.Atoi(dataRow[6])
		if err != nil {
			return err
		}
		data, _ := models.GetUserByEmail(dataRow[3])
		if data.Email == dataRow[3] {
			continue
		}
		// Assuming you have a model named 'YourModel' for database operations
		// Update this part based on your actual model structure
		yourModel := models.NewUserRequest{
			FirstName:   dataRow[1],
			LastName:    dataRow[2],
			Email:       dataRow[3],
			Age:         age,
			Country:     country,
			PhoneNumber: dataRow[4],
			Password:    dataRow[7],
		}
		// Save the data to the database
		_, err = models.ImportNewUser(yourModel)
		if err != nil {
			return err
		}
	}
	return nil
}

