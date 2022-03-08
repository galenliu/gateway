package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type ManagerProxy interface {
	HandleDeviceAdded(device DeviceProxy)
	HandleDeviceRemoved(device DeviceProxy)
	Send(messageType messages.MessageType, v any)
	GetUserProfile() *messages.PluginRegisterResponseJsonDataUserProfile
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
	Registered(m *Manager)
	SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription)
	Unload()
	CancelPairing()
	StartPairing(timeout <-chan time.Time)
	HandleDeviceSaved(data messages.DeviceSavedNotificationJsonData)
	HandleDeviceRemoved(DeviceProxy)
	CancelRemoveThing(id string)
}
