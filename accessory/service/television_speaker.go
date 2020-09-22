//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeTelevisionSpeaker = "113"

type TelevisionSpeaker struct {
	*Service
	Mute *characteristic.Mute
}

func NewTelevisionSpeaker() *TelevisionSpeaker {

	svc := TelevisionSpeaker{}
	svc.Service = New(TypeTelevisionSpeaker)
	svc.ServiceName = "TelevisionSpeaker"

	svc.Mute = characteristic.NewMute()
	svc.AddCharacteristics(svc.Mute.Characteristic)

	return &svc
}
