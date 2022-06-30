package effects

import (
	"encoding/json"
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
)

type Description = any

type MultiEffectDescription struct {
	Effects []Description `json:"effects"`
	EffectDescription
}

type MultiEffect struct {
	*Effect
	effects []Entity
}

func (e *MultiEffect) ToDescription() MultiEffectDescription {
	mulEffect := MultiEffectDescription{}
	mulEffect.EffectDescription = EffectDescription{
		Type:  e.t,
		Label: e.l,
	}
	if len(e.effects) > 0 {
		mulEffect.Effects = make([]Description, 0)
		for _, ef := range e.effects {
			mulEffect.Effects = append(mulEffect.Effects, ef)
		}
	}
	return mulEffect
}

func (e *MultiEffect) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToDescription())
}

func (e *MultiEffect) SetState(state state.State) {
	for _, e := range e.effects {
		e.SetState(state)
	}
}

func NewMultiEffect(desc MultiEffectDescription, container container.Container) *MultiEffect {
	me := &MultiEffect{
		Effect: NewEffect(desc.EffectDescription),
	}
	for _, dsc := range desc.Effects {
		e := FromDescription(dsc, container)
		if e != nil {
			me.effects = append(me.effects, e)
		}
	}
	return me
}
