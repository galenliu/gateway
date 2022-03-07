package devices

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type BinarySensorHandler interface {
	IsOn() bool
}

type BinarySensor struct {
	*Device
	Bool properties.BooleanEntity
}

func NewBinarySensor(id string, opts ...Option) *BinarySensor {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityOnOffSwitch},
	}, opts...)
	if device == nil {
		return nil
	}
	return &BinarySensor{Device: device}
}

func (device *BinarySensor) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeOnOffProperty:
			device.Bool = p.(properties.BooleanEntity)
		}
		device.Device.AddProperty(p)
	}
}

func (device *BinarySensor) IsOn() bool {
	return device.Bool.IsOn()
}
