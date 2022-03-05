package devices

type MultiLevelSwitch struct {
	*Device
}

func NewMultiLevelSwitch(id string, args ...string) *MultiLevelSwitch {
	if id == "" {
		return nil
	}
	title := "light" + id
	desc := ""
	if len(args) > 0 {
		title = args[0]
	}
	if len(args) > 1 {
		desc = args[1]
	}
	return &MultiLevelSwitch{
		NewDevice(DeviceDescription{
			Id:          id,
			AtType:      []Capability{CapabilityMultiLevelSwitch, CapabilityOnOffSwitch, CapabilityMultiLevelSwitch},
			Title:       title,
			Description: desc,
		}),
	}
}
