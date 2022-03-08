package devices

type MultiLevelSwitch struct {
	*Device
}

func NewMultiLevelSwitch(id string, opts ...Option) *MultiLevelSwitch {
	device := NewDevice(DeviceDescription{
		Id:     id,
		AtType: []Capability{CapabilityOnOffSwitch, CapabilityMultiLevelSwitch},
	}, opts...)
	if device == nil {
		return nil
	}
	return &MultiLevelSwitch{device}
}
