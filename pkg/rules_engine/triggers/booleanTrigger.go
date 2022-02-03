package triggers

import "github.com/galenliu/gateway/api/models/container"

type BooleanTriggerDescription struct {
	PropertyTriggerDescription
	OnValue bool
}

type BooleanTrigger struct {
	*PropertyTrigger
	onValue bool
}

func (b *BooleanTrigger) NewBooleanTrigger(des BooleanTriggerDescription, container container.Container) *BooleanTrigger {
	return &BooleanTrigger{
		onValue:         des.OnValue,
		PropertyTrigger: NewPropertyTrigger(des.PropertyTriggerDescription, container),
	}
}
