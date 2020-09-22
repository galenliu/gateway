//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeFilterMaintenance = "BA"

type FilterMaintenance struct {
	*Service
	FilterChangeIndication *characteristic.FilterChangeIndication
}

func NewFilterMaintenance() *FilterMaintenance {

	svc := FilterMaintenance{}
	svc.Service = New(TypeFilterMaintenance)
	svc.ServiceName = "FilterMaintenance"

	svc.FilterChangeIndication = characteristic.NewFilterChangeIndication()
	svc.AddCharacteristics(svc.FilterChangeIndication.Characteristic)

	return &svc
}
