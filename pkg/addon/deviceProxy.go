package addon

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
)

type DeviceProxy interface {
	GetId() string
	GetAdapter() AdapterProxy
	GetProperty(name string) devices.PropertyProxy
	SetCredentials(username, password string) error
}
