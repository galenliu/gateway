package controllers

import (
	"gateway/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddonController struct {
	IAddon models.IAddonManager
}

func NewAddonController(iAddon models.IAddonManager) *AddonController {

	return &AddonController{IAddon: iAddon}
}

func (addon *AddonController) HandlerGetInstalledAddons(c *gin.Context) {
	addons := addon.IAddon.GetInstallAddons()
	c.JSON(http.StatusOK, addons)
}
