package virtual

import (
	"fmt"
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

func (d *Device) SetPin(pin string) error {
	fmt.Printf("device: %s set pin: %s\n", d.GetId(), pin)
	return nil
}
