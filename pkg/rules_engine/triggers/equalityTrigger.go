package triggers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type EqualityTriggerDescription struct {
	PropertyTriggerDescription
	OnValue controls.Number `json:"onValue"`
}

type EqualityTrigger struct {
	*PropertyTrigger
	onValue controls.Number
}

func NewEqualityTrigger(desc EqualityTriggerDescription, container things.Container) *EqualityTrigger {
	return &EqualityTrigger{
		PropertyTrigger: NewPropertyTrigger(desc.PropertyTriggerDescription, container),
		onValue:         desc.OnValue,
	}
}

func (e *EqualityTrigger) Start() {
	e.PropertyTrigger.Start()
	e.property.Subscribe(topic.ValueChanged, e.OnValueChanged)
}

func (e *EqualityTrigger) Stop() {
	e.PropertyTrigger.Stop()
	e.property.Unsubscribe(topic.ValueChanged, e.OnValueChanged)
}

func (e *EqualityTrigger) OnValueChanged(v any) {
	value, ok := v.(controls.Number)
	if !ok {
		return
	}
	on := value == e.onValue
	e.Publish(topic.StateChanged, state.State{
		On:    on,
		Value: value,
	})
}

func (e *EqualityTrigger) ToDescription() EqualityTriggerDescription {
	return EqualityTriggerDescription{
		PropertyTriggerDescription: e.PropertyTrigger.ToDescription(),
		OnValue:                    e.onValue,
	}
}

func (e *EqualityTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToDescription())
}
