//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeAccessoryInformation = "3E"

type AccessoryInformation struct {
	*Service
	Identify         *characteristic.Identify
	Manufacturer     *characteristic.Manufacturer
	Model            *characteristic.Model
	Name             *characteristic.Name
	SerialNumber     *characteristic.SerialNumber
	FirmwareRevision *characteristic.FirmwareRevision
}

func NewAccessoryInformation() *AccessoryInformation {

	svc := AccessoryInformation{}
	svc.Service = New(TypeAccessoryInformation)
	svc.ServiceName = "AccessoryInformation"

	svc.Identify = characteristic.NewIdentify()
	svc.AddCharacteristics(svc.Identify.Characteristic)

	svc.Manufacturer = characteristic.NewManufacturer()
	svc.AddCharacteristics(svc.Manufacturer.Characteristic)

	svc.Model = characteristic.NewModel()
	svc.AddCharacteristics(svc.Model.Characteristic)

	svc.Name = characteristic.NewName()
	svc.AddCharacteristics(svc.Name.Characteristic)

	svc.SerialNumber = characteristic.NewSerialNumber()
	svc.AddCharacteristics(svc.SerialNumber.Characteristic)

	svc.FirmwareRevision = characteristic.NewFirmwareRevision()
	svc.AddCharacteristics(svc.FirmwareRevision.Characteristic)

	return &svc
}
