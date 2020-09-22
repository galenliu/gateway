//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeAirQualitySensor = "8D"

type AirQualitySensor struct {
	*Service
	AirQuality *characteristic.AirQuality
}

func NewAirQualitySensor() *AirQualitySensor {

	svc := AirQualitySensor{}
	svc.Service = New(TypeAirQualitySensor)
	svc.ServiceName = "AirQualitySensor"

	svc.AirQuality = characteristic.NewAirQuality()
	svc.AddCharacteristics(svc.AirQuality.Characteristic)

	return &svc
}
