//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeFan = "40"

type Fan struct {
	*Service
	On *characteristic.On
}

func NewFan() *Fan {

	svc := Fan{}
	svc.Service = New(TypeFan)
	svc.ServiceName = "Fan"

	svc.On = characteristic.NewOn()
	svc.AddCharacteristics(svc.On.Characteristic)

	return &svc
}
