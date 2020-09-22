//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeAirPurifier = "BB"

type AirPurifier struct {
	*Service
	Active                  *characteristic.Active
	CurrentAirPurifierState *characteristic.CurrentAirPurifierState
	TargetAirPurifierState  *characteristic.TargetAirPurifierState
}

func NewAirPurifier() *AirPurifier {

	svc := AirPurifier{}
	svc.Service = New(TypeAirPurifier)
	svc.ServiceName = "AirPurifier"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	svc.CurrentAirPurifierState = characteristic.NewCurrentAirPurifierState()
	svc.AddCharacteristics(svc.CurrentAirPurifierState.Characteristic)

	svc.TargetAirPurifierState = characteristic.NewTargetAirPurifierState()
	svc.AddCharacteristics(svc.TargetAirPurifierState.Characteristic)

	return &svc
}
