//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeIrrigationSystem = "CF"

type IrrigationSystem struct {
	*Service
	Active      *characteristic.Active
	ProgramMode *characteristic.ProgramMode
	InUse       *characteristic.InUse
}

func NewIrrigationSystem() *IrrigationSystem {

	svc := IrrigationSystem{}
	svc.Service = New(TypeIrrigationSystem)
	svc.ServiceName = "IrrigationSystem"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	svc.ProgramMode = characteristic.NewProgramMode()
	svc.AddCharacteristics(svc.ProgramMode.Characteristic)

	svc.InUse = characteristic.NewInUse()
	svc.AddCharacteristics(svc.InUse.Characteristic)

	return &svc
}
