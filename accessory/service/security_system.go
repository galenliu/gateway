//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeSecuritySystem = "7E"

type SecuritySystem struct {
	*Service
	SecuritySystemCurrentState *characteristic.SecuritySystemCurrentState
	SecuritySystemTargetState  *characteristic.SecuritySystemTargetState
}

func NewSecuritySystem() *SecuritySystem {

	svc := SecuritySystem{}
	svc.Service = New(TypeSecuritySystem)
	svc.ServiceName = "SecuritySystem"

	svc.SecuritySystemCurrentState = characteristic.NewSecuritySystemCurrentState()
	svc.AddCharacteristics(svc.SecuritySystemCurrentState.Characteristic)

	svc.SecuritySystemTargetState = characteristic.NewSecuritySystemTargetState()
	svc.AddCharacteristics(svc.SecuritySystemTargetState.Characteristic)

	return &svc
}
