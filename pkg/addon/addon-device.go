package addon

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type AddonDeviceProxy interface {
	properties.DeviceProxy
	SetCredentials(username, password string) error
}

type AddonDevice struct {
	adapter *AddonAdapter
	*devices.Device
}

func NewDevice(adapter *AddonAdapter, atType []string, id, title string) *AddonDevice {
	return &AddonDevice{
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

func (d AddonDevice) NotifyPropertyChanged(prop properties.PropertyDescription) {
	d.adapter.SendPropertyChangedNotification(d.GetId(), prop)
}
