package plugin

import (
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
)

type Device struct {
	AdapterId string `json:"adapterId"`
	adapter   *Adapter
	*internal.Device
}

func NewDeviceFormString(desc string, adapter *Adapter) *Device {
	data := []byte(desc)
	device := &Device{}
	device.adapter = adapter
	device.AdapterId = adapter.ID
	device.Properties = make(map[string]internal.Prop)

	device.AtContext = json.Get(data, "@context").ToString()
	device.AtType = json.Get(data, "@type").ToString()
	device.Name = json.Get(data, "name").ToString()
	device.Description = json.Get(data, "description").ToString()
	device.Device = internal.NewDeviceFormString(adapter, json.Get(data, "id").ToString())
	if device.Device == nil {
		return nil
	}
	var properties map[string]string
	json.Get(data, "properties").ToVal(&properties)
	if properties != nil {
		for name, prop := range properties {
			p := NewPropertyFormString(prop, device)
			if p != nil {
				p.Name = name
				device.Properties[name] = p
			}
		}
	}
	return device
}

func (device *Device) NotifyValueChanged(property *internal.Property) {
	data, err := json.Marshal(property)
	if err != nil {
		return
	}
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishPropertyChanged(data)
}

func (device *Device) connectedNotify(connected bool) {
	device.adapter.plugin.pluginServer.manager.Eventbus.PublishConnected(device.ID, connected)
}
