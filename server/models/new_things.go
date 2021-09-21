package models

import (
	"github.com/galenliu/gateway-addon/devices"
	"github.com/galenliu/gateway/pkg/logging"
)

type NewThingsManager interface {
	GetDevicesBytes() map[string][]byte
	GetDevicesMaps() map[string]*devices.Device
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
