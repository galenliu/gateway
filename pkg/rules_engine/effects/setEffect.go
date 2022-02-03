package effects

import (
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
)

type SetEffectDescription struct {
	PropertyEffectDescription
	Value any
}

type SetEffect struct {
	*PropertyEffect
	on    bool
	value any
}

func NewSetEffect(des SetEffectDescription, container container.Container) *SetEffect {
	e := &SetEffect{}
	e.PropertyEffect = NewPropertyEffect(des.PropertyEffectDescription, container)
	return e
}

func (s *SetEffect) ToDescription() SetEffectDescription {
	return SetEffectDescription{
		PropertyEffectDescription: s.PropertyEffect.ToDescription(),
		Value:                     s.value,
	}
}

func (s *SetEffect) SetState(state triggers.State) {
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
