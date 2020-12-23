package controllers

import (
	"gateway/addons"
	"gateway/pkg/log"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type AddonController struct {
}

func NewAddonController() *AddonController {
	return &AddonController{}
}

//  GET /addons
func (addon *AddonController) HandlerGetAddons(c *gin.Context) {
	data, err := addons.GetInstallAddons()
	if err != nil {
		c.String(500, "")
		return
	}
	c.String(http.StatusOK, string(data))
}

// PUT /addon/:id
func (addon *AddonController) HandlerSetAddon(c *gin.Context) {
	addonId := c.Param("addon_id")
	var body map[string]bool
	data, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	if body["enabled"] == true {
		err = addons.EnableAddon(addonId)
	} else {
		err = addons.DisableAddon(addonId)
	}
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, body)
}

func (addon *AddonController) HandlerInstallAddon(c *gin.Context) {

	var data struct {
		ID       string `json:"id"`
		Url      string `json:"url"`
		Checksum string `json:"checksum"`
	}
	b, err := ioutil.ReadAll(c.Request.Body)
	str := string(b)
	log.Info(str)
	err = json.Unmarshal(b, &data)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	go addons.InstallAddonFromUrl(data.ID, data.Url, data.Checksum, true)
	c.Status(200)
}

//GET /addon/:addonId/config
func (addon *AddonController) HandlerGetAddonConfig(c *gin.Context) {

}

func (addon *AddonController) HandlerSetAddonConfig(c *gin.Context) {

}
