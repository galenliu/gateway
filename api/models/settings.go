package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/tidwall/gjson"
)

type AddonInfo struct {
	Urls           []string `json:"urls"`
	Architecture   string   `json:"architecture"`
	Version        string   `json:"version"`
	NodeVersion    string   `json:"nodeVersion"`
	PythonVersions []string `json:"pythonVersions"`
}

type SettingsStore interface {
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type Settings struct {
	storage   SettingsStore
	addonInfo AddonInfo
}

func NewSettingsModel(addonUrl []string, storage SettingsStore) *Settings {
	s := Settings{}
	s.addonInfo = AddonInfo{
		Urls:           addonUrl,
		Architecture:   util.GetArch(),
		Version:        constant.Version,
		NodeVersion:    util.GetNodeVersion(),
		PythonVersions: util.GetPythonVersion(),
	}
	log.Debugf("settings model: %s", util.JsonIndent(s.addonInfo))
	s.storage = storage
	return &s
}

func (s *Settings) GetTunnelInfo() string {
	token, err := s.storage.GetSetting("tunneltoken")
	if err != nil {
		log.Info("Tunnel domain not set.")
		return "Not Set."
	}
	name := gjson.GetBytes([]byte(token), "name").String()
	base := gjson.GetBytes([]byte(token), "base").String()
	log.Info("Tunnel domain found. Tunnel name is: &s and tunnel domain is: %s", name, base)
	tunnelDomain := fmt.Sprintf("https://%s.%s", name, base)
	return tunnelDomain
}

func (s *Settings) GetAddonInfo() AddonInfo {
	return s.addonInfo
}

func (s *Settings) GetTemperatureUnits() (string, error) {
	return s.storage.GetSetting("localization.units.temperature")
}
