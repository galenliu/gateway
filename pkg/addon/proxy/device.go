package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
)

type DeviceInstance interface {
	devices.Entity
	SetCredentials(username, password string) error
}

type Device struct {
	DeviceInstance
}

func NewDevice(d DeviceInstance) *Device {
	return &Device{
		DeviceInstance: d,
	}
}

func (d *Device) GetProperty(id string) PropertyProxy {
	prop := d.DeviceInstance.GetPropertyEntity(id)
	p, ok := prop.(PropertyProxy)
	if ok {
		return p
	}
	return nil
}

func (d *Device) SetCredentials(username, password string) error {
	return fmt.Errorf("SetCredentials not implemented")
}
