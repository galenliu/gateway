//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeLightSensor = "84"

type LightSensor struct {
	*Service
	CurrentAmbientLightLevel *characteristic.CurrentAmbientLightLevel
}

func NewLightSensor() *LightSensor {

	svc := LightSensor{}
	svc.Service = New(TypeLightSensor)
	svc.ServiceName = "LightSensor"

	svc.CurrentAmbientLightLevel = characteristic.NewCurrentAmbientLightLevel()
	svc.AddCharacteristics(svc.CurrentAmbientLightLevel.Characteristic)

	return &svc
}
