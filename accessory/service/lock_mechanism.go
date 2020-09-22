//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeLockMechanism = "45"

type LockMechanism struct {
	*Service
	LockCurrentState *characteristic.LockCurrentState
	LockTargetState  *characteristic.LockTargetState
}

func NewLockMechanism() *LockMechanism {

	svc := LockMechanism{}
	svc.Service = New(TypeLockMechanism)
	svc.ServiceName = "LockMechanism"

	svc.LockCurrentState = characteristic.NewLockCurrentState()
	svc.AddCharacteristics(svc.LockCurrentState.Characteristic)

	svc.LockTargetState = characteristic.NewLockTargetState()
	svc.AddCharacteristics(svc.LockTargetState.Characteristic)

	return &svc
}
