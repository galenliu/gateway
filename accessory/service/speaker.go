//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeSpeaker = "113"

type Speaker struct {
	*Service
	Mute *characteristic.Mute
}

func NewSpeaker() *Speaker {

	svc := Speaker{}
	svc.Service = New(TypeSpeaker)
	svc.ServiceName = "Speaker"

	svc.Mute = characteristic.NewMute()
	svc.AddCharacteristics(svc.Mute.Characteristic)

	return &svc
}
