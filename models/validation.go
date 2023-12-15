package models

import (
	"CarCrudv2/helpers"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	valid "github.com/beego/beego/v2/core/validation"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (a NewUserRequest) NewUserValidate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required, validation.Length(3, 50)),
		validation.Field(&a.LastName, validation.Required, validation.Length(3, 50)),
		validation.Field(&a.Email, validation.Required, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))),
		validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.Age, validation.Required),
		validation.Field(&a.Role, validation.Required),
		validation.Field(&a.Country, validation.Required),
		validation.Field(&a.Password, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{6}$`))),
	)
}

func (a UpdateUserRequest) UdateUserValidate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required, validation.Length(3, 50)),
		validation.Field(&a.LastName, validation.Required, validation.Length(3, 50)),
		validation.Field(&a.Email, validation.Required, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))),
		validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.Age, validation.Required),
		validation.Field(&a.Role, validation.Required),
		validation.Field(&a.Country, validation.Required),
	)
}

func (a VerifyEmailOTPRequest) ValidateUsernameOtp() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Username, validation.Required),
		validation.Field(&a.Otp, validation.Required),
	)
}

func (a GetNewCarRequest) NewCarValidate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CarName, validation.Required),
		validation.Field(&a.Model, validation.Required),
		validation.Field(&a.Type, validation.Required),
	)
}

func (a UpdateCarRequest) UpdateCarValidate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CarName, validation.Required),
		validation.Field(&a.Model, validation.Required),
		validation.Field(&a.Type, validation.Required),
	)
}

func (a InserNewHomeSettingRequest) NewHomeSettingValidate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Section, validation.Required),
		validation.Field(&a.Type, validation.Required),
	)
}

func (a UpdateHomeSetingRequest) UpdateHomeSetingValidate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Section, validation.Required),
		validation.Field(&a.Type, validation.Required),
	)
}

// Beggo defaoult vallidations keys and messages

func FormValidation() {
	valid.SetDefaultMessage(map[string]string{
		"Required":     "REQUIRED_ERROR",
		"Min":          "MINIMUM_ERROR",
		"Max":          "MAXIMUM_ERROR",
		"Range":        "RANGE_ERROR",
		"MinSize":      "MINIMUM_ERROR",
		"MaxSize":      "MAXIMUM_ERROR",
		"Length":       "LENGTH_ERROR",
		"Alpha":        "ALPHA_ERROR",
		"Numeric":      "NUMERIC_ERROR",
		"AlphaNumeric": "ALPHANUMERIC_ERROR",
		"Match":        "MATCH_ERROR",
		"NoMatch":      "NOT_MATCH_ERROR",
		"AlphaDash":    "ALPHADASH_ERROR",
		"Email":        "EMAIL_ERROR",
		"IP":           "IP_ERROR",
		"Base64":       "BASE64_ERROR",
		"Mobile":       "PHONE_ERROR",
		"Tel":          "PHONE_ERROR",
		"Phone":        "PHONE_ERROR",
		"ZipCode":      "ZIP_ERROR",
	})
	valid.AddCustomFunc("IsMobile", IsMobile)
}

func Validate(lang string, bodyData interface{}) (bool, []string) {
	valid := valid.Validation{}
	var verror []string
	isValid, _ := valid.Valid(bodyData)
	if !isValid {
		for _, err := range valid.Errors {
			words := strings.Fields(err.Error())
			var errorcode string
			var errorstr string
			if len(words) > 2 {
				if len(words) == 3 {
					integer := words[len(words)-1]
					valueStr := strings.TrimSuffix(strings.TrimPrefix(integer, "int="), ")")
					value, _ := strconv.Atoi(valueStr)
					if value == 0 {
						matchstr := strings.TrimSuffix(strings.TrimPrefix(integer, "string="), ")")
						splitResult := strings.Split(words[1], "%")
						str := fmt.Sprintf(helpers.GetLangaugeMessage(lang, splitResult[0]), matchstr)
						errorstr = words[0] + ":- " + str
					} else {
						splitResult := strings.Split(words[1], "%")
						str := fmt.Sprintf(helpers.GetLangaugeMessage(lang, splitResult[0]), value)
						errorstr = words[0] + ":- " + str
					}
				}
				if len(words) == 4 {
					integer1 := words[2]
					valueStr1 := strings.TrimSuffix(strings.TrimPrefix(integer1, "int="), ",")
					value1, _ := strconv.Atoi(valueStr1)
					integer2 := words[3]
					valueStr := strings.TrimSuffix(strings.TrimPrefix(integer2, "int="), ")")
					value, _ := strconv.Atoi(valueStr)
					splitResult := strings.Split(words[1], "%")
					str := fmt.Sprintf(helpers.GetLangaugeMessage(lang, splitResult[0]), value1, value)
					errorstr = words[0] + ":- " + str
				}
			} else {
				errorcode = words[len(words)-1]
				errorstr = words[0] + ":- " + helpers.GetLangaugeMessage(lang, errorcode)
			}
			verror = append(verror, errorstr)
		}
		return false, verror
	}
	return true, verror
}

func IsMobile(v *valid.Validation, obj interface{}, key string) {
	name, ok := obj.(string)
	if !ok {
		// wrong use case?
		return
	}
	pattern := `^[6-9][0-9]{9}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(name) {
		splitResult := strings.Split(key, ".")
		v.SetError("MobileNumber", splitResult[0]+" PHONE_ERROR")
	}
}
