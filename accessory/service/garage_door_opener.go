//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeGarageDoorOpener = "41"

type GarageDoorOpener struct {
	*Service
	CurrentDoorState    *characteristic.CurrentDoorState
	TargetDoorState     *characteristic.TargetDoorState
	ObstructionDetected *characteristic.ObstructionDetected
}

func NewGarageDoorOpener() *GarageDoorOpener {

	svc := GarageDoorOpener{}
	svc.Service = New(TypeGarageDoorOpener)
	svc.ServiceName = "GarageDoorOpener"

	svc.CurrentDoorState = characteristic.NewCurrentDoorState()
	svc.AddCharacteristics(svc.CurrentDoorState.Characteristic)

	svc.TargetDoorState = characteristic.NewTargetDoorState()
	svc.AddCharacteristics(svc.TargetDoorState.Characteristic)

	svc.ObstructionDetected = characteristic.NewObstructionDetected()
	svc.AddCharacteristics(svc.ObstructionDetected.Characteristic)

	return &svc
}
