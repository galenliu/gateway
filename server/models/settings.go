package models

import "github.com/galenliu/gateway/pkg/logging"

type SettingsModel struct {
	logger logging.Logger
}

func NewSettingsModel(logger logging.Logger) *SettingsModel {
	settingsModel := SettingsModel{}
	settingsModel.logger = logger
	return &settingsModel
}
