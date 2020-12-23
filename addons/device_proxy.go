package addons

import (
	json "github.com/json-iterator/go"
)

type DeviceProxy struct {
	*Device
	adapter    *AdapterProxy
}

func NewDeviceProxy(proxy *AdapterProxy, device interface{}) *DeviceProxy {
	devProxy := &DeviceProxy{}
	return devProxy
}

func NewDevice(adapter *AdapterProxy, _id string) (dev *DeviceProxy) {
	dev = &DeviceProxy{}
	dev.adapter = adapter
	dev.ID = _id
	dev.AtContext = Context
	return
}

//Device GetId
func (proxy *DeviceProxy) GetId() string {
	return proxy.ID
}

func (proxy *DeviceProxy) GetAdapter() *AdapterProxy {
	return proxy.adapter
}

func (proxy *DeviceProxy) FindProperty(propName string) *Property {
	p := proxy.Properties[propName]
	return p
}

func (proxy *DeviceProxy) AddProperty(props ...*Property) {
	for _, prop := range props {
		prop.EventEmitter = proxy
		proxy.Properties[prop.Name] = prop
	}
}

func (proxy *DeviceProxy) AsDict() string {
	data, err := json.MarshalToString(proxy)
	if err != nil {
		return ""
	}
	return data
}

func (proxy *DeviceProxy) AppendType(ts ...string) {
	for _, t := range ts {
		proxy.AtType = append(proxy.AtType, t)
	}
}

func (proxy *DeviceProxy) OnPropertyChanged(prop interface{}) {

}

func (proxy *DeviceProxy) doPropertyChanged(message string) {
	var m Property
	_ = json.UnmarshalFromString(message, m)
	p := proxy.FindProperty(m.Name)
	p.Value = m.Value

}
