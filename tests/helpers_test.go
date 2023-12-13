package test

import (
	"CarCrudv2/helpers"
	"CarCrudv2/helpers/common"
	"log"
	"testing"
)

func TestHelpers(t *testing.T) {

	t.Run("Check Hashdata", func(t *testing.T) {
		password := "123456"
		data, err := helpers.HashData(password)
		if err != nil {
			t.Errorf(err.Error())
		}
		log.Print("sucees to verify hasdata")

		err = helpers.VerifyHashedData(data, password)
		if err != nil {
			t.Errorf(err.Error())
		}
		log.Print("sucees to verify veiry hasdata")
	})

	t.Run("Send Otp", func(t *testing.T) {
		_, err := common.VerifyEmail("rideshnath.siliconithub@gmail.com", "ridesh")
		if err != nil {
			t.Errorf(err.Error())
		}
		log.Print("sucees to send mail")
	})
	t.Run("Generate random Otp", func(t *testing.T) {
		str := helpers.GenerateOtp()
		t.Log(str)
	})
}
