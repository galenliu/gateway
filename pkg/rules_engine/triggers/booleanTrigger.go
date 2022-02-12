package triggers

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	json "github.com/json-iterator/go"
)

type BooleanTriggerDescription struct {
	PropertyTriggerDescription
	OnValue bool `json:"onValue"`
}

type BooleanTrigger struct {
	*PropertyTrigger
	onValue bool
}

func NewBooleanTrigger(des BooleanTriggerDescription, container container.Container) *BooleanTrigger {
	tri := &BooleanTrigger{
		onValue:         des.OnValue,
		PropertyTrigger: NewPropertyTrigger(des.PropertyTriggerDescription, container),
	}
	return tri
}

func (b *BooleanTrigger) Start() {
	b.PropertyTrigger.Start()
	b.property.Subscribe(topic.StateChanged, b.OnValueChanged)
}

func (b *BooleanTrigger) Stop() {
	b.PropertyTrigger.Stop()
	b.property.Unsubscribe(topic.StateChanged, b.OnValueChanged)
}

func (b *BooleanTrigger) ToDescription() BooleanTriggerDescription {
	desc := BooleanTriggerDescription{
		PropertyTriggerDescription: b.PropertyTrigger.ToDescription(),
		OnValue:                    b.onValue,
	}
	return desc
}

func (b *BooleanTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.ToDescription())
}

func (b *BooleanTrigger) OnValueChanged(propValue any) {
	v, ok := propValue.(bool)
	if ok {
		if v == b.onValue {
			b.Publish(topic.StateChanged, state.State{On: true, Value: v})
		} else {
			b.Publish(topic.StateChanged, state.State{On: false, Value: v})
		}
	}
}
