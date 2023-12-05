package controllers

import (
	"CarCrudv2/helpers"
	"CarCrudv2/models"
	"encoding/json"
	"net/http"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

type HomeSettingController struct {
	beego.Controller
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
	var bodyData models.GetHomeSettingRequest
	err := helpers.RequestBody(c.Ctx, &bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	data, err := models.GetHomeSetting(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
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
	var bodyData models.InserNewHomeSettingRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	// fmt.Println(bodyData)
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	if bodyData.Type == "Banner url" || bodyData.Type == "Logo url" {
		_, fileheader, err := c.GetFile("file")
		if err != nil {
			helpers.ApiFailure(c.Ctx, "File Getting Error", http.StatusBadRequest, 1001)
			return
		}
		filedName := "file"
		uploadDir := "./uploads/Homesetings/images/"
		filepaths, err := helpers.UploadFile(c.Controller, filedName, fileheader, uploadDir)
		if err != nil {
			helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
			return
		}
		bodyData.Value = filepaths
	}
	bodyData.Key = helpers.GenereateKeyForHomeSection(bodyData.Section, bodyData.Type)
	output, err := models.InsertNewHomeSetting(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
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
	var bodyData models.UpdateHomeSetingRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	data, err := models.GetHomeSetting(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	if bodyData.Type == "Banner url" || bodyData.Type == "Logo url" {
		_, fileheader, err := c.GetFile("file")
		if err != nil {
			helpers.ApiFailure(c.Ctx, "File Getting Error", http.StatusBadRequest, 1001)
			return
		}
		filedName := "file"
		uploadDir := "./uploads/Homesetings/images/"
		filepaths, err := helpers.UploadFile(c.Controller, filedName, fileheader, uploadDir)
		if err != nil {
			helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
			return
		}
		bodyData.Value = filepaths
		err = os.Remove(data.Value)
		if err != nil {
			helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
			return
		}
	}
	bodyData.Key = helpers.GenereateKeyForHomeSection(bodyData.Section, bodyData.Type)
	output, err := models.UpdateHomeSeting(bodyData)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
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
	var bodyData models.GetHomeSettingRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	res, err := models.UserWiseHomeseting(bodyData.Id)
	if err != nil {
		helpers.ApiFailure(c.Ctx, err.Error(), http.StatusBadRequest, 1001)
		return
	}
	helpers.ApiSuccess(c.Ctx, res, http.StatusOK, 1000)
}
