package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Device struct {
	handler AdapterProxy
	*devices.Device
}

func NewDevice(atType []string, id, title string) *Device {
	return &Device{
		handler: nil,
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

func (d *Device) NotifyPropertyChanged(prop properties.PropertyDescription) {
	d.handler.SendPropertyChangedNotification(d.GetId(), prop)
}

func (d *Device) SetCredentials(username, password string) error {
	return fmt.Errorf("SetCredentials not implemented")
}

func (d *Device) AddProperty(p PropertyProxy) {
	p.SetHandler(d)
	d.Device.AddProperty(p.GetName(), p)
}

func (d *Device) GetAdapter() AdapterProxy {
	return d.handler
}

func (d *Device) GetProperty(id string) PropertyProxy {
	p := d.Device.GetProperty(id)
	if p != nil {
		p, ok := p.(PropertyProxy)
		if ok {
			return p
		}
	}
	return nil
}

func (d *Device) SetHandler(h AdapterProxy) {
	d.handler = h
}
