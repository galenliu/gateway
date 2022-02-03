package effects

import (
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
)

type PulseEffectDescription struct {
	PropertyEffectDescription
	Value any
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
	return e
}

func (s *PulseEffect) ToDescription() SetEffectDescription {
	return SetEffectDescription{
		PropertyEffectDescription: s.PropertyEffect.ToDescription(),
		Value:                     s.value,
	}
}

func (s *PulseEffect) SetState(state triggers.State) {
	if !s.on && state.On {
		s.on = true
		_, err := s.Property.Set(state.Value)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}
	if s.on && !state.On {
		s.on = false
	}
}
