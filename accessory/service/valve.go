//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeValve = "D0"

type Valve struct {
	*Service
	Active    *characteristic.Active
	InUse     *characteristic.InUse
	ValveType *characteristic.ValveType
}

func NewValve() *Valve {

	svc := Valve{}
	svc.Service = New(TypeValve)
	svc.ServiceName = "Valve"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	svc.InUse = characteristic.NewInUse()
	svc.AddCharacteristics(svc.InUse.Characteristic)

	svc.ValveType = characteristic.NewValveType()
	svc.AddCharacteristics(svc.ValveType.Characteristic)

	return &svc
}
