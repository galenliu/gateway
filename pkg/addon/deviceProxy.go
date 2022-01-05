package addon

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Device struct {
	adapter *Adapter
	*devices.Device
}

func NewDevice(adapter *Adapter, atType []string, id, title string) *Device {
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

func (d Device) NotifyPropertyChanged(prop properties.PropertyDescription) {
	d.adapter.SendPropertyChangedNotification(d.GetId(), prop)
}
