package addon

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
	"sync"
)

type OnDeviceSavedFunc func(deviceId string, device *devices.Device)
type OnSetCredentialsFunc func(deviceId, username, password string)

//type OnSetPinFunc func(deivceId string, devices.P) error

type Adapter struct {
	Devices     map[string]DeviceProxy
	manager     *Manager
	locker      *sync.Mutex
	cancelChan  chan struct{}
	Id          string
	name        string
	packageName string
	IsPairing   bool
	verbose     bool
	pluginId    string
}

func NewAdapter(manager *Manager, adapterId, name string) *Adapter {

	adapter := &Adapter{}
	adapter.manager = manager
	adapter.Id = adapterId
	adapter.name = name
	adapter.locker = new(sync.Mutex)
	adapter.Devices = make(map[string]DeviceProxy)
	adapter.cancelChan = make(chan struct{})
	adapter.verbose = true
	return adapter
}

func (a *Adapter) HandleDeviceAdded(device DeviceProxy) {
	a.Devices[device.GetId()] = device
	a.manager.handleDeviceAdded(device)
}

func (a *Adapter) SendError(message string) {
	a.manager.send(messages.MessageType_PluginErrorNotification, messages.PluginErrorNotificationJsonData{
		Message:  message,
		PluginId: a.manager.packageName,
	})
}

//func (a *Adapter) ConnectedNotify(device *DeviceProxy, connected bool) {
//	a.manager.sendConnectedStateNotification(device, connected)
//}

// SendPairingPrompt 向前端UI发送提示
func (a *Adapter) SendPairingPrompt(prompt, url string, did string) {

	var u *string
	if url != "" {
		u = &url
	} else {
		u = nil
	}
	a.manager.send(messages.MessageType_AdapterPairingPromptNotification, messages.AdapterPairingPromptNotificationJsonData{
		AdapterId: a.GetId(),
		DeviceId:  &did,
		PluginId:  a.packageName,
		Prompt:    prompt,
		Url:       u,
	})
}

func (a *Adapter) SendUnpairingPrompt(prompt, url string, did string) {
	var u *string
	if url != "" {
		u = &url
	} else {
		u = nil
	}
	a.manager.send(messages.MessageType_AdapterUnpairingPromptNotification, messages.AdapterUnpairingPromptNotificationJsonData{
		AdapterId: a.GetId(),
		DeviceId:  &did,
		PluginId:  a.packageName,
		Prompt:    prompt,
		Url:       u,
	})
}

func (a *Adapter) Send(mt messages.MessageType, data map[string]any) {
	a.manager.send(mt, data)
}

func (a *Adapter) CancelPairing() {
	if a.verbose {
		log.Printf("adapter:(%s)- CancelPairing() not implemented", a.GetId())
	}
}

func (a *Adapter) GetId() string {
	return a.Id
}

func (a *Adapter) GetName() string {
	if a.name == "" {
		return a.Id
	}
	return a.name
}

func (a *Adapter) GetDevice(id string) DeviceProxy {
	device, ok := a.Devices[id]
	if !ok {
		return nil
	}
	return device
}

func (a *Adapter) Unload() {
	if a.verbose {
		log.Printf("adapter:(%s)- unloaded ", a.GetId())
	}
}

func (a *Adapter) HandleDeviceSaved(device DeviceProxy) {
	if a.verbose {
		log.Printf("adapter:(%s)- HandleDeviceSaved() not implemented", device.GetId())
	}
}

func (a *Adapter) HandleDeviceRemoved(device DeviceProxy) {
	delete(a.Devices, device.GetId())
}

func (a *Adapter) getDevice(id string) DeviceProxy {
	return a.Devices[id]
}

func (a *Adapter) close() {
	fmt.Print("do some thing while a close")
	a.manager.close()
}

func (a *Adapter) ProxyRunning() bool {
	return a.manager.running
}

func (a *Adapter) CloseProxy() {
	a.manager.close()
}

func (a *Adapter) setManager(m *Manager) {
	a.manager = m
}

func (a *Adapter) SetPin(deviceId string, pin any) {

}

func (a *Adapter) SetCredentials(deviceId, username, password string) {

}

func (a *Adapter) GetAddonManager() *Manager {
	return a.manager
}

func (a *Adapter) SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription) {
	a.manager.send(messages.MessageType_DevicePropertyChangedNotification, messages.DevicePropertyChangedNotificationJsonData{
		AdapterId: a.GetId(),
		DeviceId:  deviceId,
		PluginId:  a.pluginId,
		Property: messages.Property{
			Type:        property.Type,
			AtType:      property.AtType,
			Description: property.Description,
			Enum:        property.Enum,
			Maximum:     property.Maximum,
			Minimum:     property.Minimum,
			MultipleOf:  property.MultipleOf,
			Name:        property.Name,
			ReadOnly:    property.ReadOnly,
			Title:       property.Title,
			Unit:        property.Unit,
			Value:       property.Value,
		},
	})
}
