package effects

import (
	"fmt"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
	json "github.com/json-iterator/go"
)

type SetEffectDescription struct {
	PropertyEffectDescription
	Value any `json:"value"`
}

type SetEffect struct {
	*PropertyEffect
	on    bool
	value any
}

func NewSetEffect(des SetEffectDescription, container container.Container) *SetEffect {
	e := &SetEffect{}
	e.value = des.Value
	e.PropertyEffect = NewPropertyEffect(des.PropertyEffectDescription, container)
	return e
}

func (s *SetEffect) ToDescription() SetEffectDescription {
	return SetEffectDescription{
		PropertyEffectDescription: s.PropertyEffect.ToDescription(),
		Value:                     s.value,
	}
}

func (s *SetEffect) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToDescription())
}

func (s *SetEffect) SetState(state state.State) {
	if !s.on && state.On {
		s.on = true
		_, err := s.property.Set(s.value)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}
	if s.on && !state.On {
		s.on = false
	}
}
