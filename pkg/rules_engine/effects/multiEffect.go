package effects

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/state"
)

type Description interface {
}

type MultiEffectDescription struct {
	Effects []Description
	EffectDescription
}

type MultiEffect struct {
	*Effect
	effects []Entity
}

//func (e *MultiEffect) ToDescription() Description {
//	mulEffect := &MultiEffectDescription{}
//	mulEffect.EffectDescription = EffectDescription{
//		Type:  e.t,
//		Label: e.l,
//	}
//	if len(e.effects) > 0 {
//		mulEffect.Effects = make([]Description, 1)
//		for _, ef := range e.effects {
//			mulEffect.Effects = append(mulEffect.Effects, ef.ToDescription())
//		}
//	}
//	return mulEffect
//}

func (e *MultiEffect) SetState(state state.State) {
	for _, e := range e.effects {
		e.SetState(state)
	}
}

func NewMultiEffect(desc MultiEffectDescription, container container.Container) *MultiEffect {
	return nil
}
