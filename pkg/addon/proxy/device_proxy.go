package proxy

import (
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type DeviceProxy interface {
	devices.Entity
	properties.DeviceHandler
	GetProperty(id string) properties.Entity
	SetCredentials(username, password string) error
	SetPin(pin string) error
}

//type DeviceInstance interface {
//	devices.Entity
//	SetCredentials(username, password string) error
//}
//
//type Device struct {
//	DeviceInstance
//}
//
//func NewDevice(d DeviceInstance) *Device {
//	return &Device{
//		DeviceInstance: d,
//	}
//}
//
//func (d *Device) GetProperty(id string) properties.Entity {
//	prop := d.DeviceInstance.GetProperty(id)
//	p, ok := prop.(properties.Entity)
//	if ok {
//		return p
//	}
//	return nil
//}
//
//func (d *Device) SetCredentials(username, password string) error {
//	return fmt.Errorf("device:%s SetCredentials not implemented", d.GetId())
//}
//
//func (d *Device) SetPin(pin string) error {
//	return fmt.Errorf("device:%s SetPin not implemented", d.GetId())
//}
