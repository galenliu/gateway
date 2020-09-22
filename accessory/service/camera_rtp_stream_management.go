//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeCameraRTPStreamManagement = "110"

type CameraRTPStreamManagement struct {
	*Service
	SupportedVideoStreamConfiguration *characteristic.SupportedVideoStreamConfiguration
	SupportedAudioStreamConfiguration *characteristic.SupportedAudioStreamConfiguration
	SupportedRTPConfiguration         *characteristic.SupportedRTPConfiguration
	SelectedRTPStreamConfiguration    *characteristic.SelectedRTPStreamConfiguration
	StreamingStatus                   *characteristic.StreamingStatus
	SetupEndpoints                    *characteristic.SetupEndpoints
}

func NewCameraRTPStreamManagement() *CameraRTPStreamManagement {

	svc := CameraRTPStreamManagement{}
	svc.Service = New(TypeCameraRTPStreamManagement)
	svc.ServiceName = "CameraRTPStreamManagement"

	svc.SupportedVideoStreamConfiguration = characteristic.NewSupportedVideoStreamConfiguration()
	svc.AddCharacteristics(svc.SupportedVideoStreamConfiguration.Characteristic)

	svc.SupportedAudioStreamConfiguration = characteristic.NewSupportedAudioStreamConfiguration()
	svc.AddCharacteristics(svc.SupportedAudioStreamConfiguration.Characteristic)

	svc.SupportedRTPConfiguration = characteristic.NewSupportedRTPConfiguration()
	svc.AddCharacteristics(svc.SupportedRTPConfiguration.Characteristic)

	svc.SelectedRTPStreamConfiguration = characteristic.NewSelectedRTPStreamConfiguration()
	svc.AddCharacteristics(svc.SelectedRTPStreamConfiguration.Characteristic)

	svc.StreamingStatus = characteristic.NewStreamingStatus()
	svc.AddCharacteristics(svc.StreamingStatus.Characteristic)

	svc.SetupEndpoints = characteristic.NewSetupEndpoints()
	svc.AddCharacteristics(svc.SetupEndpoints.Characteristic)

	return &svc
}
