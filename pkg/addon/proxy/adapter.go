package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
	"time"
)

type Adapter struct {
	*adapter.Adapter
	name      string
	manager   ManagerProxy
	IsPairing bool
	verbose   bool
	pluginId  string
	client    *IpcClient
}

func NewAdapter(adapterId, name string) *Adapter {
	a := &Adapter{}
	a.Adapter = adapter.NewAdapter(adapterId)
	a.name = name
	a.verbose = true
	return a
}

func (a *Adapter) HandleDeviceAdded(devices ...DeviceProxy) {
	for _, device := range devices {
		if device != nil {
			a.handleDeviceAdded(device)
			a.manager.handleDeviceAdded(device)
		}
	}
}

func (a *Adapter) SendError(message string) {
	a.manager.send(messages.MessageType_PluginErrorNotification, messages.PluginErrorNotificationJsonData{
		Message:  message,
		PluginId: a.GetPackageName(),
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
		PluginId:  a.GetPackageName(),
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
		PluginId:  a.GetPackageName(),
		Prompt:    prompt,
		Url:       u,
	})
}

func (a *Adapter) Send(mt messages.MessageType, data any) {
	a.manager.send(mt, data)
}

func (a *Adapter) CancelPairing() {
	if a.verbose {
		log.Printf("adapter:(%s)- CancelPairing() not implemented", a.GetId())
	}
}

func (a *Adapter) getDevice(id string) DeviceProxy {
	device := a.Adapter.GetDevice(id)
	if device != nil {
		d, ok := device.(DeviceProxy)
		if ok {
			return d
		}
	}
	return nil
}

func (a *Adapter) unload() {
	if a.verbose {
		log.Printf("adapter:(%s)- unloaded ", a.GetId())
	}
}

func (a *Adapter) HandleDeviceSaved(msg messages.DeviceSavedNotificationJsonData) {
	if a.verbose {
		log.Printf("adapter: %s HandleDeviceSaved not implemented", msg.AdapterId)
	}
}

func (a *Adapter) HandleDeviceRemoved(device DeviceProxy) {
	if a.verbose {
		log.Printf("adapter: %s HandleDeviceRemoved not implemented", a.GetId())
	}
}

func (a *Adapter) Close() {
	fmt.Print("do some thing while a close")
	a.manager.Close()
}

func (a *Adapter) ProxyRunning() bool {
	return a.manager.IsRunning()
}

func (a *Adapter) SetPin(deviceId string, pin any) {

}

func (a *Adapter) SendPropertyChangedNotification(deviceId string, property properties.PropertyDescription) {
	a.manager.send(messages.MessageType_DevicePropertyChangedNotification, messages.DevicePropertyChangedNotificationJsonData{
		AdapterId: a.GetId(),
		DeviceId:  deviceId,
		PluginId:  a.pluginId,
		Property: messages.Property{
			Type:        property.Type,
			AtType:      &property.AtType,
			Description: &property.Description,
			Enum:        property.Enum,
			Maximum: func() *float64 {
				if v := property.Maximum; v != nil {
					f, ok := v.(float64)
					if ok {
						return &f
					}
				}
				return nil
			}(),
			Minimum: func() *float64 {
				if v := property.Minimum; v != nil {
					f, ok := v.(float64)
					if ok {
						return &f
					}
				}
				return nil
			}(),
			MultipleOf: func() *float64 {
				if v := property.MultipleOf; v != nil {
					f, ok := v.(float64)
					if ok {
						return &f
					}
				}
				return nil
			}(),
			Name:     &property.Name,
			ReadOnly: &property.ReadOnly,
			Title:    &property.Title,
			Unit:     &property.Unit,
			Value:    property.Value,
		},
	})
}

func (a *Adapter) StartPairing(timeout <-chan time.Time) {
	fmt.Printf("Adapter:%s StartPairing func not implemented\t\n", a.GetId())
}

func (a *Adapter) CancelRemoveThing(id string) {
	device := a.GetDevice(id)
	if device == nil {
		fmt.Printf("no device found")
		return
	}
	fmt.Printf("Adapter:%s CancelRemoveThing func not implemented\t\n", a.GetId())
}

func (a *Adapter) registered(manager ManagerProxy) {
	a.manager = manager
	a.pluginId = manager.getPluginId()
}

func (a *Adapter) GetPackageName() string {
	return a.pluginId
}

func (a *Adapter) GetName() string {
	return a.name
}

func (a *Adapter) handleDeviceAdded(device DeviceProxy) {
	device.SetHandler(a)
	a.AddDevice(device)
}
