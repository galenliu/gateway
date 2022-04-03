package devices

import "github.com/galenliu/gateway/pkg/addon/properties"

type MultiLevelSensor struct {
	*Device
	motion properties.NumberEntity
}

func NewMultiLevelSensor(id string, opts ...Option) *MultiLevelSensor {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityMultiLevelSensor},
	}, opts...)
	if device == nil {
		return nil
	}
	return &MultiLevelSensor{Device: device}
}

func (sensor *MultiLevelSensor) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeLevelProperty:
			sensor.motion = p.(properties.NumberEntity)
		}
		sensor.Device.AddProperty(p)
	}
}
