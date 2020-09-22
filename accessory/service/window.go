//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeWindow = "8B"

type Window struct {
	*Service
	CurrentPosition *characteristic.CurrentPosition
	TargetPosition  *characteristic.TargetPosition
	PositionState   *characteristic.PositionState
}

func NewWindow() *Window {

	svc := Window{}
	svc.Service = New(TypeWindow)
	svc.ServiceName = "Window"

	svc.CurrentPosition = characteristic.NewCurrentPosition()
	svc.AddCharacteristics(svc.CurrentPosition.Characteristic)

	svc.TargetPosition = characteristic.NewTargetPosition()
	svc.AddCharacteristics(svc.TargetPosition.Characteristic)

	svc.PositionState = characteristic.NewPositionState()
	svc.AddCharacteristics(svc.PositionState.Characteristic)

	return &svc
}
