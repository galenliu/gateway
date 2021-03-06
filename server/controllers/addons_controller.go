package controllers

import (
	"gateway/log"
	"gateway/pkg/database"
	"gateway/plugin"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
)

type AddonController struct {
}

func NewAddonController() *AddonController {
	return &AddonController{}
}

//  GET /addons
func (addon *AddonController) handlerGetAddons(c *gin.Context) {
	data, err := plugin.GetInstallAddons()
	if err != nil {
		c.String(500, "")
		return
	}
	c.String(http.StatusOK, string(data))
}

// PUT /addon/:id
func (addon *AddonController) handlerSetAddon(c *gin.Context) {
	addonId := c.Param("addon_id")
	var body map[string]bool
	data, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	if body["enabled"] == true {
		err = plugin.EnableAddon(addonId)
	} else {
		err = plugin.DisableAddon(addonId)
	}
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, body)
}

// Post /addons
func (addon *AddonController) handlerInstallAddon(c *gin.Context) {

	d, err := io.ReadAll(c.Request.Body)
	id := json.Get(d, "id").ToString()
	url := json.Get(d, "url").ToString()
	checksum := json.Get(d, "checksum").ToString()
	if id == "" || url == "" || checksum == "" || err != nil {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}
	e := plugin.InstallAddonFromUrl(id, url, checksum, true)
	if e != nil {
		c.String(http.StatusInternalServerError, "install addon err:  %v", e.Error())
		log.Error("install add-on err :", e.Error())
		return
	}
	key := "addons." + id
	setting, ee := database.GetSetting(key)
	if ee != nil {
		c.String(http.StatusInternalServerError, "install addon err: "+ee.Error())
		log.Error("install add-on err : %v", ee.Error())
		return
	}
	c.String(http.StatusOK, setting)
	return
}

//GET /addon/:addonId/config
func (addon *AddonController) HandlerGetAddonConfig(c *gin.Context) {

}

func (addon *AddonController) HandlerSetAddonConfig(c *gin.Context) {

}
