//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeLightbulb = "43"

type Lightbulb struct {
	*Service
	On *characteristic.On
}

func NewLightbulb() *Lightbulb {

	svc := Lightbulb{}
	svc.Service = New(TypeLightbulb)
	svc.ServiceName = "Lightbulb"

	svc.On = characteristic.NewOn()
	svc.AddCharacteristics(svc.On.Characteristic)

	return &svc
}
