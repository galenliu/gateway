package plugin

import (
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
)

type Device struct {
	adapter *Adapter
	*internal.Device
	Properties map[string]*Property
}

func NewDeviceFormString(desc string, adapter *Adapter) *Device {
	date := []byte(desc)
	device := &Device{}
	device.adapter = adapter
	device.Properties = make(map[string]*Property)
	device.Device = internal.NewDeviceFormString(desc)
	if device.Device == nil {
		return nil
	}
	var properties map[string]string
	json.Get(date, "properties").ToVal(&properties)
	if properties != nil {
		for name, prop := range properties {
			p := NewPropertyFormString(prop)
			if p != nil {
				p.NotifyValueChanged = device.NotifyValueChanged
				p.Name = name
				device.Properties[name] = p
			}
		}
	}
	return device
}

func (device Device) NotifyValueChanged(property *internal.Property) {
	data, err := json.Marshal(property)
	if err != nil {
		return
	}
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishPropertyChanged(data)
}
