//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeSlat = "B9"

type Slat struct {
	*Service
	SlatType         *characteristic.SlatType
	CurrentSlatState *characteristic.CurrentSlatState
}

func NewSlat() *Slat {

	svc := Slat{}
	svc.Service = New(TypeSlat)
	svc.ServiceName = "Slat"

	svc.SlatType = characteristic.NewSlatType()
	svc.AddCharacteristics(svc.SlatType.Characteristic)

	svc.CurrentSlatState = characteristic.NewCurrentSlatState()
	svc.AddCharacteristics(svc.CurrentSlatState.Characteristic)

	return &svc
}
