//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeLockManagement = "44"

type LockManagement struct {
	*Service
	LockControlPoint *characteristic.LockControlPoint
	Version          *characteristic.Version
}

func NewLockManagement() *LockManagement {

	svc := LockManagement{}
	svc.Service = New(TypeLockManagement)
	svc.ServiceName = "LockManagement"

	svc.LockControlPoint = characteristic.NewLockControlPoint()
	svc.AddCharacteristics(svc.LockControlPoint.Characteristic)

	svc.Version = characteristic.NewVersion()
	svc.AddCharacteristics(svc.Version.Characteristic)

	return &svc
}
