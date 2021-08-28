package models

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
)

type SettingsStore interface {
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type Settings struct {
	logger  logging.Logger
	storage SettingsStore
}

func NewSettingsModel(storage SettingsStore, logger logging.Logger) *Settings {
	settingsModel := Settings{}
	settingsModel.storage = storage
	settingsModel.logger = logger
	return &settingsModel
}

func (settings *Settings) GetTunnelInfo() string {
	token, err := settings.storage.GetSetting("tunneltoken")
	if err != nil {
		settings.logger.Info("Tunnel domain not set.")
		return "Not Set."
	}
	name := json.Get([]byte(token), "name").ToString()
	base := json.Get([]byte(token), "base").ToString()
	settings.logger.Info("Tunnel domain found. Tunnel name is: &s and tunnel domain is: %s", name, base)
	tunnelDomain := fmt.Sprintf("https://%s.%s", name, base)
	return tunnelDomain
}
