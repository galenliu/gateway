package models

import (
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/logging"
)

type NewThingsManager interface {
	GetDeviceMaps() map[string]*addon.Device
}

type NewThingsModel struct {
	Manager NewThingsManager
	logger  logging.Logger
}

func NewNewThingsModel(manager NewThingsManager, log logging.Logger) *NewThingsModel {
	n := &NewThingsModel{}
	n.Manager = manager
	n.logger = log
	return n
}
