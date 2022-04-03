package devices

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type OnOffSwitchDevice interface {
	TurnOn() error
	TurnOff() error
	Toggle() error
}

type OnOffSwitch struct {
	*Device
	OnOff            properties.BooleanEntity
	Brightness       properties.IntegerEntity
	ColorMode        properties.Entity
	Color            properties.Entity
	ColorTemperature properties.Entity
}

func NewOnOffSwitch(id string, opts ...Option) *OnOffSwitch {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityOnOffSwitch},
	}, opts...)
	if device == nil {
		return nil
	}
	return &OnOffSwitch{Device: device}
}

func (device *OnOffSwitch) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeOnOffProperty:
			device.OnOff = p.(properties.BooleanEntity)
		}
		device.Device.AddProperty(p)
	}
}

func (device *OnOffSwitch) TurnOn() error {
	return device.OnOff.TurnOn()
}

func (device *OnOffSwitch) TurnOff() error {
	return device.OnOff.TurnOff()
}

func (device *OnOffSwitch) Toggle() error {
	if device.OnOff.IsOn() {
		return device.TurnOn()
	} else {
		return device.TurnOff()
	}
}
