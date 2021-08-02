package internal

type DeviceProxy interface {
}

type Device struct {
	ID string

	AdapterId string
}

func NewDeviceFormString(des string) *Device {
	return nil
}

func (d *Device) GetId() string {
	return d.ID
}
