package controllers

import (
	"gateway/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AddonController struct {
	IAddon models.IAddonManager
}

func NewAddonController(iAddon models.IAddonManager, _log *zap.Logger) *AddonController {
	if _log != nil {
		log = _log
	} else {
		log = zap.L()
	}
	return &AddonController{IAddon: iAddon}
}

func (addon *AddonController) HandlerGetInstalledAddons(c *gin.Context) {
	addons := addon.IAddon.GetInstallAddons()
	c.JSON(http.StatusOK, addons)
}
