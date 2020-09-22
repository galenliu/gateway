//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeContactSensor = "80"

type ContactSensor struct {
	*Service
	ContactSensorState *characteristic.ContactSensorState
}

func NewContactSensor() *ContactSensor {

	svc := ContactSensor{}
	svc.Service = New(TypeContactSensor)
	svc.ServiceName = "ContactSensor"

	svc.ContactSensorState = characteristic.NewContactSensorState()
	svc.AddCharacteristics(svc.ContactSensorState.Characteristic)

	return &svc
}
