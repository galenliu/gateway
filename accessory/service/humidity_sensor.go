//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeHumiditySensor = "82"

type HumiditySensor struct {
	*Service
	CurrentRelativeHumidity *characteristic.CurrentRelativeHumidity
}

func NewHumiditySensor() *HumiditySensor {

	svc := HumiditySensor{}
	svc.Service = New(TypeHumiditySensor)
	svc.ServiceName = "HumiditySensor"

	svc.CurrentRelativeHumidity = characteristic.NewCurrentRelativeHumidity()
	svc.AddCharacteristics(svc.CurrentRelativeHumidity.Characteristic)

	return &svc
}
