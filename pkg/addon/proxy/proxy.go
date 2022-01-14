package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type PropertyProxy interface {
	devices.PropertyEntity
	SetValue(a any)
	SetHandler(DeviceProxy)
	NotifyChanged()
}

type DeviceProxy interface {
	properties.DeviceHandler
	SetHandler(proxy AdapterProxy)
	GetId() string
	GetAdapter() AdapterProxy
	ToMessage() messages.Device
	GetProperty(id string) PropertyProxy
	SetCredentials(username, password string) error
	NotifyPropertyChanged(p properties.PropertyDescription)
}

type AdapterProxy interface {
	GetId() string
	GetPackageName() string
	GetName() string
	GetDevice(deviceId string) DeviceProxy
	SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription)
	Unload()
	CancelPairing()
	StartPairing(timeout time.Duration)
	HandleDeviceSaved(DeviceProxy)
	HandleDeviceRemoved(DeviceProxy)
}
