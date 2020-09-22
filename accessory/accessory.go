package accessory

import (
	"gateway/accessory/service"
)

type Service = service.Service

type Info struct {
	Name             string
	SerialNumber     string
	Manufacturer     string
	Model            string
	FirmwareRevision string
	ID               uint64
}

type Accessory struct {
	ID         uint64 `json:"aid"`
	Services   []*Service
	Info       *service.AccessoryInformation
	Type       AccessoryType
	onIdentify func()
}

func New(info Info, typ AccessoryType) *Accessory {
	sev := service.NewAccessoryInformation()
	if name := info.Name; len(name) > 0 {
		sev.Name.SetValue(name)
	}
	if serial := info.SerialNumber; len(serial) > 0 {
		sev.SerialNumber.SetValue(serial)
	} else {
		sev.SerialNumber.SetValue("undefined")
	}
	if manufacturer := info.Manufacturer; len(manufacturer) > 0 {
		sev.Manufacturer.SetValue(manufacturer)
	} else {
		sev.Manufacturer.SetValue("undefined")
	}
	if model := info.Model; len(model) > 0 {
		sev.Model.SetValue(model)
	} else {
		sev.Model.SetValue("undefined")
	}
	if version := info.FirmwareRevision; len(version) > 0 {
		sev.FirmwareRevision.SetValue(version)
	} else {
		sev.FirmwareRevision.SetValue("undefined")
	}
	acc := &Accessory{
		ID:         info.ID,
		Services:   nil,
		Info:       sev,
		Type:       typ,
		onIdentify: nil,
	}
	acc.AddServices(acc.Info.Service)
	return acc
}

func (a *Accessory) GetServices() []*Service {
	result := make([]*Service, 0)
	for _, s := range a.Services {
		result = append(result, s)
	}
	return result
}

func (a *Accessory) AddServices(ser *Service) {
	a.Services = append(a.Services, ser)
}
