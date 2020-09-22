//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeDoor = "81"

type Door struct {
	*Service
	CurrentPosition *characteristic.CurrentPosition
	PositionState   *characteristic.PositionState
	TargetPosition  *characteristic.TargetPosition
}

func NewDoor() *Door {

	svc := Door{}
	svc.Service = New(TypeDoor)
	svc.ServiceName = "Door"

	svc.CurrentPosition = characteristic.NewCurrentPosition()
	svc.AddCharacteristics(svc.CurrentPosition.Characteristic)

	svc.PositionState = characteristic.NewPositionState()
	svc.AddCharacteristics(svc.PositionState.Characteristic)

	svc.TargetPosition = characteristic.NewTargetPosition()
	svc.AddCharacteristics(svc.TargetPosition.Characteristic)

	return &svc
}
