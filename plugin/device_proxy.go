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
	Actions    map[string]*Action   `json:"actions"`
	Events     map[string]*Event    `json:"events"`
}

func NewDeviceFormMessage(dev *rpc.Device, adapter *Adapter) *Device {
	device := &Device{
		adapter:    adapter,
		Device:     devices.NewDeviceFormMessage(dev),
		Properties: nil,
	}
	if len(device.Device.Properties) > 0 {
		device.Properties = make(map[string]*Property)
		for name, p := range device.Device.Properties {
			device.Properties[name] = NewProperty(device, p)
		}
	}

	if len(device.Device.Events) > 0 {
		device.Events = make(map[string]*Event)
		for name, e := range device.Device.Events {
			device.Events[name] = NewEvent(device, e)
		}
	}

	if len(device.Device.Actions) > 0 {
		device.Actions = make(map[string]*Action)
		for name, a := range device.Device.Actions {
			device.Actions[name] = NewAction(device, a)
		}
	}

	device.adapter = adapter
	return device
}

func NewDeviceFormString(desc string, adapter *Adapter) *Device {
	data := []byte(desc)
	device := &Device{}
	device.adapter = adapter
	device.Device = devices.NewDeviceFormString(desc)
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

func (device *Device) notifyValueChanged(property *rpc.DevicePropertyChangedNotificationMessage_Data) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishPropertyChanged(property)
}

func (device *Device) connectedNotify(connected bool) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishConnected(device.ID, connected)
}

func (device *Device) actionNotify(message *rpc.DeviceActionStatusNotificationMessage_Data) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishActionStatus(message)
}
