package devices

import "github.com/galenliu/gateway/pkg/addon/properties"

type MotionSensor struct {
	*Device
	motion properties.BooleanEntity
}

func NewMotionSensor(id string, opts ...Option) *MotionSensor {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityMotionSensor},
	}, opts...)
	if device == nil {
		return nil
	}
	return &MotionSensor{Device: device}
}

func (sensor *MotionSensor) AddProperties(props ...properties.Entity) {
	for _, p := range props {
		switch p.GetAtType() {
		case properties.TypeMotionProperty:
			sensor.motion = p.(properties.BooleanEntity)
		}
		sensor.Device.AddProperty(p)
	}
}
