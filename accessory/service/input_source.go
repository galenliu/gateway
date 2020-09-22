//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeInputSource = "D9"

type InputSource struct {
	*Service
	ConfiguredName         *characteristic.ConfiguredName
	InputSourceType        *characteristic.InputSourceType
	IsConfigured           *characteristic.IsConfigured
	CurrentVisibilityState *characteristic.CurrentVisibilityState
}

func NewInputSource() *InputSource {

	svc := InputSource{}
	svc.Service = New(TypeInputSource)
	svc.ServiceName = "InputSource"

	svc.ConfiguredName = characteristic.NewConfiguredName()
	svc.AddCharacteristics(svc.ConfiguredName.Characteristic)

	svc.InputSourceType = characteristic.NewInputSourceType()
	svc.AddCharacteristics(svc.InputSourceType.Characteristic)

	svc.IsConfigured = characteristic.NewIsConfigured()
	svc.AddCharacteristics(svc.IsConfigured.Characteristic)

	svc.CurrentVisibilityState = characteristic.NewCurrentVisibilityState()
	svc.AddCharacteristics(svc.CurrentVisibilityState.Characteristic)

	return &svc
}
