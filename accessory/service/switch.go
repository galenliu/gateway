//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeSwitch = "49"

type Switch struct {
	*Service
	On *characteristic.On
}

func NewSwitch() *Switch {

	svc := Switch{}
	svc.Service = New(TypeSwitch)
	svc.ServiceName = "Switch"

	svc.On = characteristic.NewOn()
	svc.AddCharacteristics(svc.On.Characteristic)

	return &svc
}
