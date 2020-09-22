//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeWindowCovering = "8C"

type WindowCovering struct {
	*Service
	CurrentPosition *characteristic.CurrentPosition
	TargetPosition  *characteristic.TargetPosition
	PositionState   *characteristic.PositionState
}

func NewWindowCovering() *WindowCovering {

	svc := WindowCovering{}
	svc.Service = New(TypeWindowCovering)
	svc.ServiceName = "WindowCovering"

	svc.CurrentPosition = characteristic.NewCurrentPosition()
	svc.AddCharacteristics(svc.CurrentPosition.Characteristic)

	svc.TargetPosition = characteristic.NewTargetPosition()
	svc.AddCharacteristics(svc.TargetPosition.Characteristic)

	svc.PositionState = characteristic.NewPositionState()
	svc.AddCharacteristics(svc.PositionState.Characteristic)

	return &svc
}
