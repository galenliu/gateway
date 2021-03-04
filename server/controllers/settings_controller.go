package controllers

import (
	"gateway/config"
	"gateway/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type addonInfo struct {
	Urls          []string `json:"urls"`
	Architecture  string   `json:"architecture"`
	Version       string   `json:"version"`
	NodeVersion   string   `json:"nodeVersion"`
	PythonVersion []string   `json:"pythonVersion"`
}

type SettingsController struct {
}

func NewSettingController() *SettingsController {
	return &SettingsController{}
}

func (settings *SettingsController) HandleGetAddonsInfo(g *gin.Context) {
	var addonInfo = addonInfo{
		Urls:          config.GetAddonListUrls(),
		Architecture:  util.GetArch(),
		Version:       util.Version,
		NodeVersion:   util.GetNodeVersion(),
		PythonVersion: util.GetPythonVersion(),
	}
	//data, err := json.Marshal(addonInfo)
	//if err != nil {
	//	log.Error("marshal err", zap.Error(err))
	//	g.String(http.StatusInternalServerError, "marshal err")
	//	return
	//}
	g.JSON(http.StatusOK, addonInfo)
}
