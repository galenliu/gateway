package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
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
	logger    logging.Logger
	storage   SettingsStore
	addonInfo AddonInfo
}

func NewSettingsModel(addonUrl []string, storage SettingsStore, logger logging.Logger) *Settings {
	s := Settings{}
	s.addonInfo = AddonInfo{
		Urls:           addonUrl,
		Architecture:   util.GetArch(),
		Version:        constant.Version,
		NodeVersion:    util.GetNodeVersion(),
		PythonVersions: util.GetPythonVersion(),
	}
	logger.Debugf("settings model: %s", util.JsonIndent(s.addonInfo))
	s.storage = storage
	s.logger = logger
	return &s
}

func (s *Settings) GetTunnelInfo() string {
	token, err := s.storage.GetSetting("tunneltoken")
	if err != nil {
		s.logger.Info("Tunnel domain not set.")
		return "Not Set."
	}
	name := json.Get([]byte(token), "name").ToString()
	base := json.Get([]byte(token), "base").ToString()
	s.logger.Info("Tunnel domain found. Tunnel name is: &s and tunnel domain is: %s", name, base)
	tunnelDomain := fmt.Sprintf("https://%s.%s", name, base)
	return tunnelDomain
}

func (s *Settings) GetAddonInfo() AddonInfo {
	return s.addonInfo
}

func (s *Settings) GetTemperatureUnits() (string, error) {
	return s.storage.GetSetting("localization.units.temperature")
}
