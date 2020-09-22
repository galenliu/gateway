//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeCarbonMonoxideSensor = "7F"

type CarbonMonoxideSensor struct {
	*Service
	CarbonMonoxideDetected *characteristic.CarbonMonoxideDetected
}

func NewCarbonMonoxideSensor() *CarbonMonoxideSensor {

	svc := CarbonMonoxideSensor{}
	svc.Service = New(TypeCarbonMonoxideSensor)
	svc.ServiceName = "CarbonMonoxideSensor"

	svc.CarbonMonoxideDetected = characteristic.NewCarbonMonoxideDetected()
	svc.AddCharacteristics(svc.CarbonMonoxideDetected.Characteristic)

	return &svc
}
