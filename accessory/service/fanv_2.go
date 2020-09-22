//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeFanv2 = "B7"

type Fanv2 struct {
	*Service
	Active *characteristic.Active
}

func NewFanv2() *Fanv2 {

	svc := Fanv2{}
	svc.Service = New(TypeFanv2)
	svc.ServiceName = "Fanv2"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	return &svc
}
