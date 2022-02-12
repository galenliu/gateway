package triggers

import things "github.com/galenliu/gateway/api/models/container"

type EqualityTriggerDescription struct {
	*PropertyTriggerDescription
}

type EqualityTrigger struct {
	*PropertyTrigger
}

func NewEqualityTrigger(desc EqualityTriggerDescription, container things.Container) *EqualityTrigger {
	return nil
}
