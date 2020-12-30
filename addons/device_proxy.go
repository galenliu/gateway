package addons

import (
	"fmt"
	"gateway/event"
	addon "gitee.com/liu_guilin/gateway-addon-golang"
)

type NewDeviceModel struct {
	*addon.Device
	Properties map[string]*addon.Property
	Actions    map[string]*addon.Action
}

type DeviceProxy struct {
	*addon.Device
	adapter    *AdapterProxy
	Properties map[string]*PropertyProxy `json:"properties"`
	Actions    map[string]*ActionProxy   `json:"action"`
}

func NewDeviceProxy(adapter *AdapterProxy, model *NewDeviceModel) *DeviceProxy {
	devProxy := &DeviceProxy{}
	devProxy.Device = model.Device
	devProxy.adapter = adapter
	devProxy.Properties = make(map[string]*PropertyProxy)
	devProxy.Actions = make(map[string]*ActionProxy)
	for name, p := range model.Properties {
		devProxy.Properties[name] = NewPropertyProxy(devProxy, p)
		p.DoPropertyChanged = devProxy.OnPropertyChanged
	}
	for name, a := range model.Actions {
		devProxy.Actions[name] = NewActionProxy(devProxy, a)
	}
	return devProxy
}

func (proxy *DeviceProxy) GetAdapter() *AdapterProxy {
	return proxy.adapter
}

func (proxy *DeviceProxy) AppendType(ts ...string) {
	for _, t := range ts {
		proxy.AtType = append(proxy.AtType, t)
	}
}

func (proxy *DeviceProxy) AddProperties(props ...*PropertyProxy) {
	for _, p := range props {
		proxy.Properties[p.Name] = p
	}
}

func (proxy *DeviceProxy) AddActions(actions ...*ActionProxy) {
	for _, a := range actions {
		proxy.Actions[a.Name] = a
	}
}

func (proxy *DeviceProxy) OnPropertyChanged(prop *addon.Property) {

}

func (proxy *DeviceProxy) notifyPropertyChanged(prop *PropertyProxy) {
	event.FirePropertyChanged(prop)
}

func (proxy *DeviceProxy) SetProperty(propName string, value interface{}) (interface{}, error) {
	prop, ok := proxy.Properties[propName]
	if !ok {
		return nil, fmt.Errorf("property name(%s)  err", propName)
	}

	return prop.Value, nil

}
