package triggers

import things "github.com/galenliu/gateway/api/models/container"

type EventTriggerDescription struct {
}

type EventTrigger struct {
	*PropertyTrigger
}

func NewEventTrigger(desc EventTriggerDescription, container things.Container) *EventTrigger {
	return nil
}
