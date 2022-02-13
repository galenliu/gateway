package triggers

import (
	things "github.com/galenliu/gateway/api/models/container"
	json "github.com/json-iterator/go"
)

type EventTriggerDescription struct {
	TriggerDescription
	Thing string `json:"thing"`
	Event string `json:"event"`
}

type EventTrigger struct {
	*Trigger
	thing     string
	event     string
	container things.Container
	stopped   bool
}

func NewEventTrigger(desc EventTriggerDescription, container things.Container) *EventTrigger {
	return &EventTrigger{
		Trigger:   NewTrigger(desc.TriggerDescription),
		thing:     desc.Thing,
		event:     desc.Event,
		stopped:   false,
		container: container,
	}
}

func (t *EventTrigger) ToDescription() EventTriggerDescription {
	return EventTriggerDescription{
		TriggerDescription: t.Trigger.ToDescription(),
		Thing:              t.thing,
		Event:              t.event,
	}
}

func (t *EventTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.ToDescription())
}

func (t *EventTrigger) Start() {
	thing := t.container.GetThing(t.thing)
	if thing != nil && !t.stopped {

	}
}

func (t *EventTrigger) Stop() {

}
