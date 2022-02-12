package effects

import (
	"encoding/json"
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
)

type PulseEffectDescription struct {
	PropertyEffectDescription
	Value any `json:"value"`
}

type PulseEffect struct {
	*PropertyEffect
	on       bool
	value    any
	oldValue any
}

func NewPulseEffect(des PulseEffectDescription, container container.Container) *PulseEffect {
	e := &PulseEffect{}
	e.PropertyEffect = NewPropertyEffect(des.PropertyEffectDescription, container)
	e.value = des.Value
	return e
}

func (s *PulseEffect) ToDescription() PulseEffectDescription {
	return PulseEffectDescription{
		PropertyEffectDescription: s.PropertyEffect.ToDescription(),
		Value:                     s.value,
	}
}

func (s *PulseEffect) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToDescription())
}

func (s *PulseEffect) SetState(state state.State) {
	if !s.on && state.On {
		s.on = true
		_, err := s.property.Set(state.Value)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}
	if s.on && !state.On {
		s.on = false
	}
}
