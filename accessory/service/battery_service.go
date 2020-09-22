//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeBatteryService = "96"

type BatteryService struct {
	*Service
	BatteryLevel     *characteristic.BatteryLevel
	ChargingState    *characteristic.ChargingState
	StatusLowBattery *characteristic.StatusLowBattery
}

func NewBatteryService() *BatteryService {

	svc := BatteryService{}
	svc.Service = New(TypeBatteryService)
	svc.ServiceName = "BatteryService"

	svc.BatteryLevel = characteristic.NewBatteryLevel()
	svc.AddCharacteristics(svc.BatteryLevel.Characteristic)

	svc.ChargingState = characteristic.NewChargingState()
	svc.AddCharacteristics(svc.ChargingState.Characteristic)

	svc.StatusLowBattery = characteristic.NewStatusLowBattery()
	svc.AddCharacteristics(svc.StatusLowBattery.Characteristic)

	return &svc
}
