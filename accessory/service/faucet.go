//THis File is AUTO-GENERATED
package service

import "gateway/accessory/characteristic"

const TypeFaucet = "D7"

type Faucet struct {
	*Service
	Active *characteristic.Active
}

func NewFaucet() *Faucet {

	svc := Faucet{}
	svc.Service = New(TypeFaucet)
	svc.ServiceName = "Faucet"

	svc.Active = characteristic.NewActive()
	svc.AddCharacteristics(svc.Active.Characteristic)

	return &svc
}
