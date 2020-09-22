//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeDoorbell = "121"

type Doorbell struct {
	*Service
	ProgrammableSwitchEvent *characteristic.ProgrammableSwitchEvent
}

func NewDoorbell() *Doorbell {

	svc := Doorbell{}
	svc.Service = New(TypeDoorbell)
	svc.ServiceName = "Doorbell"

	svc.ProgrammableSwitchEvent = characteristic.NewProgrammableSwitchEvent()
	svc.AddCharacteristics(svc.ProgrammableSwitchEvent.Characteristic)

	return &svc
}
