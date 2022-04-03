package devices

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Thermostat struct {
	*Device
	Temperature       properties.NumberEntity
	TargetTemperature properties.NumberEntity
	HeatingCooling    properties.StringEntity
	ThermostatMode    properties.StringEntity
}

func NewThermostat(id string, opts ...Option) *OnOffSwitch {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityThermostat},
	}, opts...)
	if device == nil {
		return nil
	}
	return &OnOffSwitch{Device: device}
}

func (device *Thermostat) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeTemperatureProperty:
			device.Temperature = p.(properties.NumberEntity)
		case properties.TypeTargetTemperatureProperty:
			device.Temperature = p.(properties.NumberEntity)
		case properties.TypeHeatingCoolingProperty:
			device.HeatingCooling = p.(properties.StringEntity)
		case properties.TypeThermostatModeProperty:
			device.ThermostatMode = p.(properties.StringEntity)
		}
		device.Device.AddProperty(p)
	}
}
