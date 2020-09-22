//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeMicrophone = "112"

type Microphone struct {
	*Service
	Volume *characteristic.Volume
	Mute   *characteristic.Mute
}

func NewMicrophone() *Microphone {

	svc := Microphone{}
	svc.Service = New(TypeMicrophone)
	svc.ServiceName = "Microphone"

	svc.Volume = characteristic.NewVolume()
	svc.AddCharacteristics(svc.Volume.Characteristic)

	svc.Mute = characteristic.NewMute()
	svc.AddCharacteristics(svc.Mute.Characteristic)

	return &svc
}
