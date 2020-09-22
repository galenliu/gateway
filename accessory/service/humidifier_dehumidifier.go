//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeHumidifierDehumidifier = "BD"

type HumidifierDehumidifier struct {
	*Service
	CurrentRelativeHumidity            *characteristic.CurrentRelativeHumidity
	CurrentHumidifierDehumidifierState *characteristic.CurrentHumidifierDehumidifierState
	TargetHumidifierDehumidifierState  *characteristic.TargetHumidifierDehumidifierState
	Active                             *characteristic.Active
}

func NewHumidifierDehumidifier() *HumidifierDehumidifier {

	svc := HumidifierDehumidifier{}
	svc.Service = New(TypeHumidifierDehumidifier)
	svc.ServiceName = "HumidifierDehumidifier"

	svc.CurrentRelativeHumidity = characteristic.NewCurrentRelativeHumidity()
	svc.AddCharacteristics(svc.CurrentRelativeHumidity.Characteristic)

	svc.CurrentHumidifierDehumidifierState = characteristic.NewCurrentHumidifierDehumidifierState()
	svc.AddCharacteristics(svc.CurrentHumidifierDehumidifierState.Characteristic)

	svc.TargetHumidifierDehumidifierState = characteristic.NewTargetHumidifierDehumidifierState()
	svc.AddCharacteristics(svc.TargetHumidifierDehumidifierState.Characteristic)

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	return &svc
}
