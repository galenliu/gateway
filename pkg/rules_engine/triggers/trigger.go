package triggers

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/plugin"
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
	plugin.Bus
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
	bus                     EventBus.Bus
	t                       string
	label                   string
	onValueChangedCallbacks map[string]func()
}

func NewTrigger(des TriggerDescription) *Trigger {
	return &Trigger{
		bus:                     EventBus.New(),
		t:                       des.Type,
		label:                   des.Label,
		onValueChangedCallbacks: make(map[string]func(), 1),
	}
}

func (t *Trigger) ToDescription() *TriggerDescription {
	return &TriggerDescription{
		Type:  t.t,
		Label: t.label,
	}
}

func (t *Trigger) Subscribe(topic topic.Topic, fn any) func() {
	err := t.bus.Subscribe(string(topic), fn)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return func() {
		err := t.bus.Unsubscribe(string(topic), fn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func (t *Trigger) Publish(topic2 topic.Topic, args ...any) {
	t.bus.Publish(string(topic2), args...)
}

func (t *Trigger) Unsubscribe(topic2 topic.Topic, f any) {
	err := t.bus.Unsubscribe(string(topic2), f)
	if err != nil {
		fmt.Println(err.Error())
		return
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
