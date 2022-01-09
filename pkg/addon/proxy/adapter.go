package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/adapter"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
)

type Adapter struct {
	adapter.Adapter
	manager     *Manager
	packageName string
	IsPairing   bool
	verbose     bool
	pluginId    string
}

func NewAdapter(manager *Manager, adapterId, name string) *Adapter {
	a := &Adapter{}
	a.Adapter = adapter.Adapter{
		Id:   adapterId,
		Name: name,
	}
	a.manager = manager
	a.verbose = true
	return a
}

func (a *Adapter) HandleDeviceAdded(device addon.DeviceProxy) {
	a.AddDevice(device)
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
	if a.Name == "" {
		return a.Id
	}
	return a.Name
}

func (a *Adapter) GetDevice(id string) addon.DeviceProxy {
	device := a.Adapter.GetDevice(id)
	if device != nil {
		d, ok := device.(addon.DeviceProxy)
		if ok {
			return d
		}
	}
	return nil
}

func (a *Adapter) Unload() {
	if a.verbose {
		log.Printf("adapter:(%s)- unloaded ", a.GetId())
	}
}

func (a *Adapter) HandleDeviceSaved(device addon.DeviceProxy) {
	if a.verbose {
		log.Printf("adapter:(%s)- HandleDeviceSaved() not implemented", device.GetId())
	}
}

func (a *Adapter) HandleDeviceRemoved(device addon.DeviceProxy) {
	a.Adapter.RemoveDevice(device.GetId())
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

func (a *Adapter) SendPropertyChangedNotification(deviceId string, property addon.PropertyDescription) {
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
