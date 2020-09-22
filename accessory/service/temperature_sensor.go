//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeTemperatureSensor = "8A"

type TemperatureSensor struct {
	*Service
	CurrentTemperature *characteristic.CurrentTemperature
}

func NewTemperatureSensor() *TemperatureSensor {

	svc := TemperatureSensor{}
	svc.Service = New(TypeTemperatureSensor)
	svc.ServiceName = "TemperatureSensor"

	svc.CurrentTemperature = characteristic.NewCurrentTemperature()
	svc.AddCharacteristics(svc.CurrentTemperature.Characteristic)

	return &svc
}
