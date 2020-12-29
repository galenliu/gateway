package addons

import (
	json "github.com/json-iterator/go"
)

type DeviceProxy struct {
	*Device
	adapter *AdapterProxy
}

func NewDeviceProxy(proxy *AdapterProxy, device *Device) *DeviceProxy {
	devProxy := &DeviceProxy{}
	devProxy.Device = device
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

func (proxy *DeviceProxy) OnPropertyChanged(prop interface{}) {

}

func (proxy *DeviceProxy) doPropertyChanged(propInfo string) {
	var prop Property
	_ = json.UnmarshalFromString(propInfo, &prop)

}

func (proxy *DeviceProxy)setValue(prop Property){

}

