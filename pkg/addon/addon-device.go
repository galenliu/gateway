package addon

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type AddonDeviceProxy interface {
	properties.DeviceProxy
	SetCredentials(username, password string) error
	ToMessage() messages.Device
}

type AddonDevice struct {
	adapter AddonAdapterProxy
	*devices.Device
}

func NewDevice(adapter AddonAdapterProxy, atType []string, id, title string) *AddonDevice {
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

func (d AddonDevice) SetCredentials(username, password string) error {
	return nil
}
