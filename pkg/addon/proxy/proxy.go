package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type PropertyProxy interface {
	properties.Entity
	SetValue(a any)
}

type DeviceProxy interface {
	devices.Entity
	GetProperty(id string) PropertyProxy
	properties.DeviceHandler
	SetCredentials(username, password string) error
	SetPin(pin string) error
}

type ManagerProxy interface {
	HandleDeviceAdded(device DeviceProxy)
	HandleDeviceRemoved(device DeviceProxy)
	Send(messageType messages.MessageType, v any)
	Close()
	GetPluginId() string
	IsRunning() bool
}

type AdapterProxy interface {
	adapter.Entity
	devices.AdapterHandler
	GetDevice(deviceId string) DeviceProxy
	GetName() string
	GetPackageName() string
	SetManager(m ManagerProxy)
	SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription)
	Unload()
	CancelPairing()
	StartPairing(timeout <-chan time.Time)
	HandleDeviceSaved(DeviceProxy)
	HandleDeviceRemoved(DeviceProxy)
}
