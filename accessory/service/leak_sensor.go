//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeLeakSensor = "83"

type LeakSensor struct {
	*Service
	LeakDetected *characteristic.LeakDetected
}

func NewLeakSensor() *LeakSensor {

	svc := LeakSensor{}
	svc.Service = New(TypeLeakSensor)
	svc.ServiceName = "LeakSensor"

	svc.LeakDetected = characteristic.NewLeakDetected()
	svc.AddCharacteristics(svc.LeakDetected.Characteristic)

	return &svc
}
