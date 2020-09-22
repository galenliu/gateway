//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeOutlet = "47"

type Outlet struct {
	*Service
	On          *characteristic.On
	OutletInUse *characteristic.OutletInUse
}

func NewOutlet() *Outlet {

	svc := Outlet{}
	svc.Service = New(TypeOutlet)
	svc.ServiceName = "Outlet"

	svc.On = characteristic.NewOn()
	svc.AddCharacteristics(svc.On.Characteristic)

	svc.OutletInUse = characteristic.NewOutletInUse()
	svc.AddCharacteristics(svc.OutletInUse.Characteristic)

	return &svc
}
