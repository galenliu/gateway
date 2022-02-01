package triggers

import (
	"github.com/galenliu/gateway/pkg/rules_engine"
	"github.com/galenliu/gateway/pkg/rules_engine/property"
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
	bus property.Bus
	*PropertyTrigger
	value     controls.Number
	levelType string
}

func NewLevelTrigger(des LevelTriggerDescription, bus property.Bus, things property.ThingsHandler) *LevelTrigger {
	return &LevelTrigger{
		bus:             bus,
		PropertyTrigger: NewPropertyTrigger(des.PropertyTriggerDescription, things),
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
	t.bus.Pub(rules_engine.StateChanged, t.Property.ToDescription().Thing, t.Property.ToDescription().Id, State{On: on, Value: v})
}
