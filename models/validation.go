package models

import (
	"regexp"

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