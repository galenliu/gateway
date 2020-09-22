//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeHeaterCooler = "BC"

type HeaterCooler struct {
	*Service
	Active                   *characteristic.Active
	CurrentHeaterCoolerState *characteristic.CurrentHeaterCoolerState
	TargetHeaterCoolerState  *characteristic.TargetHeaterCoolerState
	CurrentTemperature       *characteristic.CurrentTemperature
}

func NewHeaterCooler() *HeaterCooler {

	svc := HeaterCooler{}
	svc.Service = New(TypeHeaterCooler)
	svc.ServiceName = "HeaterCooler"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	svc.CurrentHeaterCoolerState = characteristic.NewCurrentHeaterCoolerState()
	svc.AddCharacteristics(svc.CurrentHeaterCoolerState.Characteristic)

	svc.TargetHeaterCoolerState = characteristic.NewTargetHeaterCoolerState()
	svc.AddCharacteristics(svc.TargetHeaterCoolerState.Characteristic)

	svc.CurrentTemperature = characteristic.NewCurrentTemperature()
	svc.AddCharacteristics(svc.CurrentTemperature.Characteristic)

	return &svc
}
