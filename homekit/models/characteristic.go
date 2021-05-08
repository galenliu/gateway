package models

import "github.com/brutella/hc/characteristic"

type HomeKitCharacteristic interface {
}

type HomeKitCharacteristicProxy struct {
	characteristic *characteristic.Characteristic
}

func NewHomeKitCharacteristic(data []byte) *HomeKitCharacteristic {
	return nil
}

func (c *HomeKitCharacteristicProxy) name() {

}
