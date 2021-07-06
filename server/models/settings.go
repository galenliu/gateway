package models

import "github.com/galenliu/gateway/pkg/logging"

type Settings  struct {
	logger logging.Logger
}

func NewSettingsModel(logger logging.Logger) *Settings {
	settingsModel := Settings{}
	settingsModel.logger = logger
	return &settingsModel
}
