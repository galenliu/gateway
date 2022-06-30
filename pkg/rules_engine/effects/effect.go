package effects

import (
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	"github.com/tidwall/gjson"
)

type Entity interface {
	SetState(state2 state.State)
	GetType() string
	GetLabel() string
}

type EffectDescription struct {
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
}

type Effect struct {
	t string
	l string
}

func NewEffect(des EffectDescription) *Effect {
	e := &Effect{
		t: des.Type,
		l: des.Label,
	}
	return e
}

func (e *Effect) ToDescription() EffectDescription {
	des := EffectDescription{
		Type:  e.t,
		Label: e.l,
	}
	return des
}

func (e *Effect) GetType() string {
	return e.t
}

func (e *Effect) GetLabel() string {
	return e.l
}

func (e *Effect) SetState(s state.State) {
}

func FromDescription(a any, container container.Container) Entity {
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("marshal error: %s\n", err.Error())
	}
	t := gjson.GetBytes(data, "type").String()
	switch t {
	case TypeSetEffect:
		var desc SetEffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("marshal error: %s\n", err.Error())
			return nil
		}
		return NewSetEffect(desc, container)
	case TypeMultiEffect:
		var desc MultiEffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("marshal error: %s\n", err.Error())
			return nil
		}
		return NewMultiEffect(desc, container)
	case TypePulseEffect:
		var desc PulseEffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("marshal error: %s\n", err.Error())
			return nil
		}
		return NewPulseEffect(desc, container)
	case TypeEffect:
		var desc EffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Printf("marshal error: %s\n", err.Error())
			return nil
		}
		return NewEffect(desc)
	default:
		fmt.Println("type error")
		return nil
	}
}

const TypeEffect = "Effect"
const TypeActionEffect = "ActionEffect"
const TypeMultiEffect = "MultiEffect"
const TypeNotificationEffect = "NotificationEffect"
const TypeNotifierOutletEffect = "NotifierOutletEffect"
const TypeSetEffect = "SetEffect"
const TypePulseEffect = "PulseEffect"
