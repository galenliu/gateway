package effects

import (
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
	json "github.com/json-iterator/go"
)

type Entity interface {
	SetState(state2 triggers.State)
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

func (e *Effect) ToDescription() *EffectDescription {
	des := &EffectDescription{
		Type:  e.t,
		Label: e.l,
	}
	return des
}

func FromDescription(data []byte, container container.Container) Entity {
	t := json.Get(data, "type").ToString()
	switch t {
	case TypeSetEffect:
		var desc SetEffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return NewSetEffect(desc, container)
	case TypeMultiEffect:
		var desc MultiEffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return NewMultiEffect(desc, container)
	case TypePulseEffect:
		var desc PulseEffectDescription
		err := json.Unmarshal(data, &desc)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return NewPulseEffect(desc, container)

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
