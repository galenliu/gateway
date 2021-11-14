package plugin

import (
	"github.com/galenliu/gateway/pkg/addon"
)

type Device struct {
	adapter *Adapter
	*addon.Device
}

//func NewDeviceFormMessage(msg *rpc.Device, adapter *Adapter) *Device {
//	device := &Device{
//		adapter: adapter,
//		Device:  devices.NewDeviceFormMessage(msg),
//	}
//	if len(msg.Properties) > 0 {
//		for _, p := range msg.Properties {
//			device.addProperty(NewProperty(device, properties.NewPropertyFormMessage(p)))
//		}
//	}
//	if len(msg.Events) > 0 {
//		for _, e := range msg.Events {
//			device.addEvent(NewEvent(device, events.NewEventFormMessage(e)))
//		}
//	}
//	if len(msg.Actions) > 0 {
//		for _, a := range msg.Actions {
//			device.addAction(NewAction(device, actions.NewActionFormMessage(a)))
//		}
//	}
//	device.adapter = adapter
//	return device
//}
//
//func NewDeviceFormString(desc string, adapter *Adapter) *Device {
//	data := []byte(desc)
//	device := &Device{}
//	device.adapter = adapter
//	device.Device = devices.NewDeviceFormString(desc)
//	device.properties = make(map[string]*Property)
//	var p map[string]string
//	json.Get(data, "p").ToVal(&p)
//	if p != nil {
//		for name, prop := range p {
//			p := NewPropertyFormString(device, prop)
//			if p != nil {
//				p.Name = name
//				device.properties[name] = p
//			}
//		}
//	}
//	return device
//}

//func (device *Device) addProperty(property *Property) {
//	device.Device.AddProperty(property.Property)
//	if device.properties == nil {
//		device.properties = make(map[string]*Property)
//	}
//	device.properties[property.Name] = property
//}
//
//func (device *Device) addAction(action *Action) {
//	device.Device.AddAction(action.Action)
//	if device.actions == nil {
//		device.actions = make(map[string]*Action)
//	}
//	device.actions[action.Name] = action
//}
//
//func (device *Device) addEvent(event *Event) {
//	device.Device.AddEvent(event.Event)
//	if device.events == nil {
//		device.events = make(map[string]*Event)
//	}
//	device.events[event.Name] = event
//}
//
//func (device *Device) GetProperty(name string) *Property {
//	return device.properties[name]
//}

func (device *Device) notifyValueChanged(property *addon.Property) {
	p := device.GetProperty(property.Name)
	if p.ReadOnly {
		return
	}
	p.Value = property.Value
	p.Title = property.Title
	p.Description = property.Description
	device.adapter.plugin.pluginServer.manager.bus.PublishPropertyChanged(device.GetId(), property)
}

func (device *Device) connectedNotify(connected bool) {
	device.adapter.plugin.pluginServer.manager.bus.PublishConnected(device.GetId(), connected)
}

func (device *Device) actionNotify(actionDescription *addon.Action) {
	device.adapter.plugin.pluginServer.manager.bus.PublishActionStatus(actionDescription)
}
