//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeCarbonDioxideSensor = "97"

type CarbonDioxideSensor struct {
	*Service
	CarbonDioxideDetected *characteristic.CarbonDioxideDetected
}

func NewCarbonDioxideSensor() *CarbonDioxideSensor {

	svc := CarbonDioxideSensor{}
	svc.Service = New(TypeCarbonDioxideSensor)
	svc.ServiceName = "CarbonDioxideSensor"

	svc.CarbonDioxideDetected = characteristic.NewCarbonDioxideDetected()
	svc.AddCharacteristics(svc.CarbonDioxideDetected.Characteristic)

	return &svc
}
