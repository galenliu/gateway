package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/devices"
)

type Device struct {
	adapter addon.AdapterProxy
	*devices.Device
}

func NewDevice(adapter addon.AdapterProxy, atType []string, id, title string) *Device {
	return &Device{
		adapter: adapter,
		Device: &devices.Device{
			Context:             "https://webthings.io/schemas",
			AtType:              atType,
			Id:                  id,
			Title:               title,
			Description:         "",
			Links:               nil,
			Forms:               nil,
			BaseHref:            "",
			Pin:                 nil,
			Properties:          nil,
			Actions:             nil,
			Events:              nil,
			CredentialsRequired: false,
		},
	}
}

func (d Device) NotifyPropertyChanged(prop addon.PropertyDescription) {
	d.adapter.SendPropertyChangedNotification(d.GetId(), prop)
}

func (d Device) SetCredentials(username, password string) error {
	return fmt.Errorf("SetCredentials not implemented")
}

func (d Device) GetAdapter() addon.AdapterProxy {
	return d.adapter
}

func (d Device) GetProperty(id string) addon.PropertyProxy {
	p := d.Device.GetProperty(id)
	if p != nil {
		p, ok := p.(addon.PropertyProxy)
		if ok {
			return p
		}
	}
	return nil
}
