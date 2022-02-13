package triggers

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

const (
	LESS    = "LESS"
	EQUAL   = "EQUAL"
	GREATER = "GREATER"
)

type LevelTriggerDescription struct {
	PropertyTriggerDescription
	Value     controls.Number `json:"value"`
	LevelType string          `json:"levelType"`
}

type LevelTrigger struct {
	*PropertyTrigger
	value     controls.Number
	levelType string
}

func NewLevelTrigger(des LevelTriggerDescription, container container.Container) *LevelTrigger {
	return &LevelTrigger{
		PropertyTrigger: NewPropertyTrigger(des.PropertyTriggerDescription, container),
		value:           des.Value,
		levelType:       des.LevelType,
	}
}

func (t *LevelTrigger) OnValueChanged(a any) {
	v, ok := a.(controls.Number)
	if !ok {
		return
	}
	on := false
	switch t.levelType {
	case LESS:
		if v < t.value {
			on = true
		}
		break
	case EQUAL:
		if v == t.value {
			on = true
		}
		break
	case GREATER:
		if v < t.value {
			on = true
		}
		break
	}
	t.Publish(topic.StateChanged, state.State{On: on, Value: v})
}

func (t *LevelTrigger) ToDescription() LevelTriggerDescription {
	return LevelTriggerDescription{
		PropertyTriggerDescription: t.PropertyTrigger.ToDescription(),
		Value:                      t.value,
		LevelType:                  t.levelType,
	}
}

func (t *LevelTrigger) Start() {
	t.PropertyTrigger.Start()
	t.property.Subscribe(topic.ValueChanged, t.OnValueChanged)
}

func (t *LevelTrigger) Stop() {
	t.PropertyTrigger.Stop()
	t.property.Unsubscribe(topic.ValueChanged, t.OnValueChanged)
}

func (t *LevelTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.ToDescription())
}
