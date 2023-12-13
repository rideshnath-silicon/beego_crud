package controllers

import (
	"CarCrudv2/helpers"
	"CarCrudv2/models"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type HomeSettingController struct {
	beego.Controller
}


func (c *HomeSettingController) Prepare() {
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
	if len(lang) == 0 {
		lang = c.Ctx.GetCookie("lang")
		if len(lang) != 0 {
			c.Data["Lang"] = lang
		} else {
			lang = c.Ctx.Input.Header("Accept-Language")
			if len(lang) > 4 {
				lang := lang[:5] // Only compare first 5 letters.
				if lang == "en-US" || lang == "hi-IN" {
					c.Data["Lang"] = lang
				} else {
					c.Data["Lang"] = "en-US"
				}
			}
		}
	} else {
		c.Data["Lang"] = lang
	}
	c.Ctx.SetCookie("lang", lang)
}

// GetHomeSetting ...
// @Title get home settingd
// @Desciption Get settings
// @Param body body models.GetHomeSettingRequest true "Get home settings"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router / [post]
func (c *HomeSettingController) GetHomeSetting() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.GetHomeSettingRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.GetHomeSetting(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, data, http.StatusOK, 1000)
}

// InsertNewHomeSetting ...
// @Title insert home settingd
// @Desciption insert settings
// @Param section formData string true "section"
// @Param type formData string true "types are only :-'Banner url','Logo url','Title','Description'"
// @Param value formData string false "insert when type is Title or description"
// @Param file formData file false "File to be uploaded"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /create [post]
func (c *HomeSettingController) InsertNewHomeSetting() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.InserNewHomeSettingRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	// fmt.Println(bodyData)
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	err := bodyData.NewHomeSettingValidate()
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.Type == "Banner url" || bodyData.Type == "Logo url" {
		_, fileheader, err := c.GetFile("file")
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "ERROR_DIRECTORY"), http.StatusBadRequest, 1001)
			return
		}
		filedName := "file"
		uploadDir := "./uploads/Homesetings/images/"
		filepaths, err := helpers.UploadFile(c.Controller, filedName, fileheader, uploadDir)
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
			return
		}
		bodyData.Value = filepaths
	}
	bodyData.Key = helpers.GenereateKeyForHomeSection(bodyData.Section, bodyData.Type)
	output, err := models.InsertNewHomeSetting(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1002)
}

// UpdateHomeSeting ...
// @Title update home settingd
// @Desciption update settings
// @Param home_seting_id formData string true "section"
// @Param section formData string false "section"
// @Param type formData string false "types are only :-'Banner url','Logo url','Title','Description'"
// @Param value formData string false "insert when type is Title or description"
// @Param file formData file false "File to be uploaded"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /update [put]
func (c *HomeSettingController) UpdateHomeSeting() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.UpdateHomeSetingRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	err := bodyData.UpdateHomeSetingValidate()
	if err != nil {
		helpers.ApiFailure(c.Ctx,  err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.GetHomeSetting(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.Type == "Banner url" || bodyData.Type == "Logo url" {
		_, fileheader, err := c.GetFile("file")
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "ERROR_DIRECTORY"), http.StatusBadRequest, 1001)
			return
		}
		filedName := "file"
		uploadDir := "./uploads/Homesetings/images/"
		filepaths, err := helpers.UploadFile(c.Controller, filedName, fileheader, uploadDir)
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
			return
		}
		bodyData.Value = filepaths
		err = os.Remove(data.Value)
		if err != nil {
			helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "ERROR_DIRECTORY"), http.StatusBadRequest, 1001)
			return
		}
	}
	bodyData.Key = helpers.GenereateKeyForHomeSection(bodyData.Section, bodyData.Type)
	output, err := models.UpdateHomeSeting(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, output, http.StatusOK, 1003)
}

// GetUserWiseHome ..
// @Title userwise settins
// @Description users homesettion
// @Param user_id formData string true "enter user id to search"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} object
// @Failure 403
// @router /userwise [post]
func (c *HomeSettingController) GetUserWiseHome() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData models.GetHomeSettingRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, "FORM_BODY"), http.StatusBadRequest, 1001)
		return
	}
	res, err := models.UserWiseHomeseting(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, res, http.StatusOK, 1000)
}
// DeleteHomeSetting ...
// @Title Delete home setting
// @Discription Delete home settings
// @Param id path in true "Enter id for delete home settings"
// @Success 201 {string} string
// @Failure 403
// @router /delete/:id([0-9]+) [delete]
func (c *HomeSettingController) DeleteHomeSetting() {
	lang := c.Ctx.Input.GetData("Lang").(string)
	idString := c.Ctx.Input.Params()
	id, err := strconv.Atoi(idString["1"])	
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.DeleteHomeSetting(id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, helpers.GetLangaugeMessage(lang, err.Error()), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, helpers.GetLangaugeMessage(lang, data.(string)), http.StatusOK, 1000)
}
