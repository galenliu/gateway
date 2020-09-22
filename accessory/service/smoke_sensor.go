//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeSmokeSensor = "87"

type SmokeSensor struct {
	*Service
	SmokeDetected *characteristic.SmokeDetected
}

func NewSmokeSensor() *SmokeSensor {

	svc := SmokeSensor{}
	svc.Service = New(TypeSmokeSensor)
	svc.ServiceName = "SmokeSensor"

	svc.SmokeDetected = characteristic.NewSmokeDetected()
	svc.AddCharacteristics(svc.SmokeDetected.Characteristic)

	return &svc
}
