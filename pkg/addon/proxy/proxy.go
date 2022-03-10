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
	registered(m ManagerProxy)
	unload()

	// CancelPairing ## 取消配对 子类可改写业务逻辑
	CancelPairing()
	// StartPairing ## 开始配对 子类可改写业务逻辑
	StartPairing(timeout <-chan time.Time)
	// HandleDeviceSaved 当网关把一个设备进行了保存
	HandleDeviceSaved(data messages.DeviceSavedNotificationJsonData)
	// HandleDeviceRemoved 当一个设备被网关移除
	HandleDeviceRemoved(DeviceProxy)
	// CancelRemoveThing 当一个设备重新添加
	CancelRemoveThing(id string)
}

// DeviceProxy 所有Addon所有Device的抽象接口，
//   addons/devices下所有Device均实现了此接口
type DeviceProxy interface {
	// Entity Device的公共接口
	devices.Entity
	// DeviceHandler 处理Device下Property业务的抽象接口，
	properties.DeviceHandler

	// SetCredentials 设备设置用户名t 密码时
	//子类需重写这个方法，来完成设备的用户名密码设置的逻辑
	SetCredentials(username, password string) error
	// SetPin  设备需要Pin码时时调用
	//子类需重写这个方法，来完成设备的用户名密码设置的逻辑
	SetPin(pin string) error
}
