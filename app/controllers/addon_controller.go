package controllers

import (
	"github.com/gin-gonic/gin"
)

type AddonController struct {

}

func NewAddonController() *AddonController {

	return &AddonController{}
}

func (addon *AddonController) HandlerGetInstalledAddons(c *gin.Context) {

	c.Status(200)
}
