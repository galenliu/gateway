package virtual

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
)

type Device struct {
	*devices.LightBulb
}

func NewLightDevice() *Device {
	light := &Device{
		LightBulb: devices.NewLightBulb(devices.DeviceDescription{Id: "1"}),
	}
	return light
}
