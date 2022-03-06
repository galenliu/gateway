package devices

type MultiLevelSwitch struct {
	*Device
}

func NewMultiLevelSwitch(id string, opts ...Option) *MultiLevelSwitch {
	return &MultiLevelSwitch{
		NewDevice(DeviceDescription{
			Id:     id,
			AtType: []Capability{CapabilityMultiLevelSwitch, CapabilityOnOffSwitch, CapabilityMultiLevelSwitch},
		}),
	}
}
