//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeThermostat = "4A"

type Thermostat struct {
	*Service
	CurrentHeatingCoolingState *characteristic.CurrentHeatingCoolingState
	TargetHeatingCoolingState  *characteristic.TargetHeatingCoolingState
	CurrentTemperature         *characteristic.CurrentTemperature
	TargetTemperature          *characteristic.TargetTemperature
	TemperatureDisplayUnits    *characteristic.TemperatureDisplayUnits
}

func NewThermostat() *Thermostat {

	svc := Thermostat{}
	svc.Service = New(TypeThermostat)
	svc.ServiceName = "Thermostat"

	svc.CurrentHeatingCoolingState = characteristic.NewCurrentHeatingCoolingState()
	svc.AddCharacteristics(svc.CurrentHeatingCoolingState.Characteristic)

	svc.TargetHeatingCoolingState = characteristic.NewTargetHeatingCoolingState()
	svc.AddCharacteristics(svc.TargetHeatingCoolingState.Characteristic)

	svc.CurrentTemperature = characteristic.NewCurrentTemperature()
	svc.AddCharacteristics(svc.CurrentTemperature.Characteristic)

	svc.TargetTemperature = characteristic.NewTargetTemperature()
	svc.AddCharacteristics(svc.TargetTemperature.Characteristic)

	svc.TemperatureDisplayUnits = characteristic.NewTemperatureDisplayUnits()
	svc.AddCharacteristics(svc.TemperatureDisplayUnits.Characteristic)

	return &svc
}
