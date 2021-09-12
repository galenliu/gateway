package plugin

import (
	"github.com/galenliu/gateway-addon/devices"
	"github.com/galenliu/gateway-grpc"
	json "github.com/json-iterator/go"
)

type Device struct {
	adapter *Adapter
	*devices.Device
	Properties map[string]*Property `json:"properties"`
}

func NewDeviceFormString(desc string, adapter *Adapter) *Device {
	data := []byte(desc)
	device := &Device{}
	device.adapter = adapter
	device.Device = devices.NewDeviceFormString(desc)
	device.AdapterId = adapter.ID
	device.Properties = make(map[string]*Property)
	var properties map[string]string
	json.Get(data, "properties").ToVal(&properties)
	if properties != nil {
		for name, prop := range properties {
			p := NewPropertyFormString(device, prop)
			if p != nil {
				p.Name = name
				device.Properties[name] = p
			}
		}
	}
	return device
}

func (device *Device) GetProperty(name string) *Property {
	return device.Properties[name]
}

func (device *Device) notifyValueChanged(property *gateway_grpc.DevicePropertyChangedNotificationMessage_Data) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishPropertyChanged(property)
}

func (device *Device) connectedNotify(connected bool) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishConnected(device.ID, connected)
}

func (device *Device) actionNotify(message *gateway_grpc.DeviceActionStatusNotificationMessage_Data) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishActionStatus(message)
}
