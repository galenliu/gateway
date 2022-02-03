package triggers

import (
	"github.com/galenliu/gateway/pkg/bus"
)

const TypeBooleanTrigger = "BooleanTrigger"
const TypeEqualityTrigger = "EqualityTrigger"
const TypeEventTrigger = "EventTrigger"
const TypeLevelTrigger = "LevelTrigger"
const TypeMultiTrigger = "MultiTrigger"
const TypePropertyTrigger = "PropertyTrigger"
const TypeTimeTrigger = "TimeTrigger"
const TypeTrigger = "Trigger"

type Entity interface {
	bus.Bus
	Start()
	Stop()
}

type State struct {
	On    bool
	Value any
}

type TriggerDescription struct {
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
}

type Trigger struct {
	*bus.Controller
	t     string
	label string
}

func NewTrigger(des TriggerDescription) *Trigger {
	return &Trigger{
		Controller: bus.NewBusController(),
		t:          des.Type,
		label:      des.Label,
	}
}

func (t *Trigger) ToDescription() *TriggerDescription {
	return &TriggerDescription{
		Type:  t.t,
		Label: t.label,
	}
}

func FromDescription(des TriggerDescription) Entity {
	switch des.Type {
	case TypeBooleanTrigger:
		return nil
	case TypeEqualityTrigger:
		return nil
	case TypeEventTrigger:
		return nil
	case TypeLevelTrigger:
		return nil
	case TypeMultiTrigger:
		return nil
	case TypePropertyTrigger:
		return nil
	case TypeTimeTrigger:
		return nil
	case TypeTrigger:
		return nil
	default:
		return nil
	}
}
