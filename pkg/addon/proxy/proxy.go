package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"time"
)

type ManagerProxy interface {
	getAdapter(adapterId string) AdapterProxy
	getDevice(deviceId string) DeviceProxy

	handleDeviceAdded(device DeviceProxy)
	handleDeviceRemoved(device DeviceProxy)
	Send(messageType messages.MessageType, v any)
	GetUserProfile() *messages.PluginRegisterResponseJsonDataUserProfile
	Close()
	GetPluginId() string
	IsRunning() bool
}

// AdapterProxy Adapter 的抽象接口，
type AdapterProxy interface {
	adapter.Entity
	devices.AdapterHandler

	getDevice(id string) DeviceProxy
	GetName() string
	GetPackageName() string
	registered(m *Manager)
	unload()
	CancelPairing()
	StartPairing(timeout <-chan time.Time)
	HandleDeviceSaved(data messages.DeviceSavedNotificationJsonData)
	HandleDeviceRemoved(DeviceProxy)
	CancelRemoveThing(id string)
}

// DeviceProxy 所有Addon所有Device的抽象接口，
//   addons/devices下所有Device均实现了此接口
type DeviceProxy interface {
	devices.Entity
	properties.DeviceHandler
	SetCredentials(username, password string) error
	SetPin(pin string) error
}
