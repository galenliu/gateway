//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeServiceLabel = "CC"

type ServiceLabel struct {
	*Service
	ServiceLabelNamespace *characteristic.ServiceLabelNamespace
}

func NewServiceLabel() *ServiceLabel {

	svc := ServiceLabel{}
	svc.Service = New(TypeServiceLabel)
	svc.ServiceName = "ServiceLabel"

	svc.ServiceLabelNamespace = characteristic.NewServiceLabelNamespace()
	svc.AddCharacteristics(svc.ServiceLabelNamespace.Characteristic)

	return &svc
}
