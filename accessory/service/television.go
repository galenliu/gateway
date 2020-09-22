//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeTelevision = "D8"

type Television struct {
	*Service
	Active             *characteristic.Active
	ActiveIdentifier   *characteristic.ActiveIdentifier
	ConfiguredName     *characteristic.ConfiguredName
	SleepDiscoveryMode *characteristic.SleepDiscoveryMode
}

func NewTelevision() *Television {

	svc := Television{}
	svc.Service = New(TypeTelevision)
	svc.ServiceName = "Television"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	svc.ActiveIdentifier = characteristic.NewActiveIdentifier()
	svc.AddCharacteristics(svc.ActiveIdentifier.Characteristic)

	svc.ConfiguredName = characteristic.NewConfiguredName()
	svc.AddCharacteristics(svc.ConfiguredName.Characteristic)

	svc.SleepDiscoveryMode = characteristic.NewSleepDiscoveryMode()
	svc.AddCharacteristics(svc.SleepDiscoveryMode.Characteristic)

	return &svc
}
