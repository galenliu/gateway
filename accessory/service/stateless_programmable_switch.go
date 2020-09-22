//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeStatelessProgrammableSwitch = "89"

type StatelessProgrammableSwitch struct {
	*Service
	ProgrammableSwitchEvent *characteristic.ProgrammableSwitchEvent
}

func NewStatelessProgrammableSwitch() *StatelessProgrammableSwitch {

	svc := StatelessProgrammableSwitch{}
	svc.Service = New(TypeStatelessProgrammableSwitch)
	svc.ServiceName = "StatelessProgrammableSwitch"

	svc.ProgrammableSwitchEvent = characteristic.NewProgrammableSwitchEvent()
	svc.AddCharacteristics(svc.ProgrammableSwitchEvent.Characteristic)

	return &svc
}
