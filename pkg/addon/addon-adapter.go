package addon

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
	"sync"
	"time"
)

type AddonAdapterProxy interface {
	adapter.AdapterProxy
	GetDevice(deviceId string) AddonDeviceProxy
	SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription)
	Unload()
	CancelPairing()
	StartPairing(timeout time.Duration)
	HandleDeviceSaved(AddonDeviceProxy)
	HandleDeviceRemoved(AddonDeviceProxy)
}

type AddonAdapter struct {
	adapter.Adapter
	manager     *Manager
	locker      *sync.Mutex
	cancelChan  chan struct{}
	packageName string
	IsPairing   bool
	verbose     bool
	pluginId    string
}

func NewAddonAdapter(manager *Manager, adapterId, name string) *AddonAdapter {
	a := &AddonAdapter{}
	a.Adapter = adapter.Adapter{
		Id:   adapterId,
		Name: name,
	}
	a.manager = manager
	a.locker = new(sync.Mutex)
	a.cancelChan = make(chan struct{})
	a.verbose = true
	return a
}

func (a *AddonAdapter) HandleDeviceAdded(device AddonDeviceProxy) {
	a.AddDevice(device)
	a.manager.handleDeviceAdded(device)
}

func (a *AddonAdapter) SendError(message string) {
	a.manager.send(messages.MessageType_PluginErrorNotification, messages.PluginErrorNotificationJsonData{
		Message:  message,
		PluginId: a.manager.packageName,
	})
}

//func (a *AddonAdapter) ConnectedNotify(device *DeviceProxy, connected bool) {
//	a.manager.sendConnectedStateNotification(device, connected)
//}

// SendPairingPrompt 向前端UI发送提示
func (a *AddonAdapter) SendPairingPrompt(prompt, url string, did string) {

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

func (a *AddonAdapter) SendUnpairingPrompt(prompt, url string, did string) {
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

func (a *AddonAdapter) Send(mt messages.MessageType, data map[string]any) {
	a.manager.send(mt, data)
}

func (a *AddonAdapter) CancelPairing() {
	if a.verbose {
		log.Printf("adapter:(%s)- CancelPairing() not implemented", a.GetId())
	}
}

func (a *AddonAdapter) GetId() string {
	return a.Id
}

func (a *AddonAdapter) GetName() string {
	if a.Name == "" {
		return a.Id
	}
	return a.Name
}

func (a *AddonAdapter) GetDevice(id string) AddonDeviceProxy {
	device := a.Adapter.GetDevice(id)
	if device != nil {
		d, ok := device.(AddonDeviceProxy)
		if ok {
			return d
		}
	}
	return nil
}

func (a *AddonAdapter) Unload() {
	if a.verbose {
		log.Printf("adapter:(%s)- unloaded ", a.GetId())
	}
}

func (a *AddonAdapter) HandleDeviceSaved(device AddonDeviceProxy) {
	if a.verbose {
		log.Printf("adapter:(%s)- HandleDeviceSaved() not implemented", device.GetId())
	}
}

func (a *AddonAdapter) HandleDeviceRemoved(device AddonDeviceProxy) {
	a.Adapter.RemoveDevice(device.GetId())
}

func (a *AddonAdapter) close() {
	fmt.Print("do some thing while a close")
	a.manager.close()
}

func (a *AddonAdapter) ProxyRunning() bool {
	return a.manager.running
}

func (a *AddonAdapter) CloseProxy() {
	a.manager.close()
}

func (a *AddonAdapter) setManager(m *Manager) {
	a.manager = m
}

func (a *AddonAdapter) SetPin(deviceId string, pin any) {

}

func (a *AddonAdapter) SetCredentials(deviceId, username, password string) {

}

func (a *AddonAdapter) GetAddonManager() *Manager {
	return a.manager
}

func (a *AddonAdapter) SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription) {
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
