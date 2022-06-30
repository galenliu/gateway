package triggers

import (
	"encoding/json"
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/tidwall/gjson"
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
	bus.ThingsBus
	GetType() string
	Start()
	Stop()
}

type TriggerDescription struct {
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
}

type Trigger struct {
	*bus.EventBus
	t     string
	label string
}

func NewTrigger(des TriggerDescription) *Trigger {
	return &Trigger{
		EventBus: bus.NewBus(),
		t:        des.Type,
		label:    des.Label,
	}
}

func (t *Trigger) GetType() string {
	return t.t
}

func (t *Trigger) GetLabel() string {
	return t.label
}

func (t *Trigger) Start() {

}

func (t *Trigger) Stop() {

}

func (t *Trigger) ToDescription() TriggerDescription {
	return TriggerDescription{
		Type:  t.t,
		Label: t.label,
	}
}

func FromDescription(a any, container things.Container) Entity {
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Error marshalling err: %s", err.Error())
		return nil
	}
	switch gjson.GetBytes(data, "type").String() {
	case TypeBooleanTrigger:
		var desc BooleanTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewBooleanTrigger(desc, container)
	case TypeEqualityTrigger:
		var desc EqualityTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewEqualityTrigger(desc, container)
	case TypeEventTrigger:
		var desc EventTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewEventTrigger(desc, container)
	case TypeLevelTrigger:
		var desc LevelTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewLevelTrigger(desc, container)
	case TypeMultiTrigger:
		var desc MultiTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewMultiTrigger(desc, container)
	case TypePropertyTrigger:
		var desc PropertyTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewPropertyTrigger(desc, container)
	case TypeTimeTrigger:
		var desc TimerTriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewTimerTrigger(desc)
	case TypeTrigger:
		var desc TriggerDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("unmarshal err: %s", err.Error())
			return nil
		}
		return NewTrigger(desc)
	default:
		return nil
	}
}
