package triggers

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
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
	}
}

func (t *LevelTrigger) onValueChanged(v controls.Number) {
	on := true
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
	t.Publish(rules_engine.StateChanged, t.PropertyTrigger.property.ToDescription().Thing, t.property.ToDescription().Id, State{On: on, Value: v})
}
