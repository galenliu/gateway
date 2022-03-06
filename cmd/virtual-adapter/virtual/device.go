package virtual

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Device struct {
	*devices.Device
}

func NewDevice(device devices.Entity) *Device {
	return &Device{device.GetDevice()}
}

func (d *Device) SetPin(pin string) error {
	fmt.Printf("device: %s set pin: %s\n", d.GetId(), pin)
	return nil
}

func (d *Device) SetCredentials(username string, password string) error {
	fmt.Printf("device: %s set credentials: user:%s  password: %s  \t\n", d.GetId(), username, password)
	return nil
}

func (d *Device) addProperties(props ...properties.Entity) {
	for _, p := range props {
		d.AddProperty(NewProperty(p))
	}
}
