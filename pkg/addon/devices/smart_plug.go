package devices

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type SmartPlugDevice interface {
	OnOffSwitchDevice
}

type SmartPlug struct {
	*Device
	OnOff       properties.BooleanEntity
	Level       properties.NumberEntity
	Power       properties.NumberEntity
	PowerFactor properties.NumberEntity
	Voltage     properties.NumberEntity
	Current     properties.NumberEntity
	Frequency   properties.NumberEntity
}

func NewSmartPlug(id string, opts ...Option) *Light {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilitySmartPlug, CapabilityOnOffSwitch},
	}, opts...)
	if device == nil {
		return nil
	}
	return &Light{Device: device}
}

func (plug *SmartPlug) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeOnOffProperty:
			plug.OnOff = p.(properties.BooleanEntity)
		case properties.TypeLevelProperty:
			plug.Level = p.(properties.NumberEntity)
		case properties.TypeInstantaneousPowerProperty:
			plug.Power = p.(properties.NumberEntity)
		case properties.TypeInstantaneousPowerFactorProperty:
			plug.PowerFactor = p.(properties.NumberEntity)
		case properties.TypeVoltageProperty:
			plug.Voltage = p.(properties.NumberEntity)
		}
		plug.Device.AddProperty(p)
	}
}

func (plug *SmartPlug) TurnOn() error {
	return plug.OnOff.TurnOn()
}

func (plug *SmartPlug) TurnOff() error {
	return plug.OnOff.TurnOff()
}

func (plug *SmartPlug) Toggle() error {
	if plug.OnOff.IsOn() {
		return plug.TurnOn()
	} else {
		return plug.TurnOff()
	}
}
