//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeOccupancySensor = "86"

type OccupancySensor struct {
	*Service
	OccupancyDetected *characteristic.OccupancyDetected
}

func NewOccupancySensor() *OccupancySensor {

	svc := OccupancySensor{}
	svc.Service = New(TypeOccupancySensor)
	svc.ServiceName = "OccupancySensor"

	svc.OccupancyDetected = characteristic.NewOccupancyDetected()
	svc.AddCharacteristics(svc.OccupancyDetected.Characteristic)

	return &svc
}
