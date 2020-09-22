//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeMotionSensor = "85"

type MotionSensor struct {
	*Service
	MotionDetected *characteristic.MotionDetected
}

func NewMotionSensor() *MotionSensor {

	svc := MotionSensor{}
	svc.Service = New(TypeMotionSensor)
	svc.ServiceName = "MotionSensor"

	svc.MotionDetected = characteristic.NewMotionDetected()
	svc.AddCharacteristics(svc.MotionDetected.Characteristic)

	return &svc
}
