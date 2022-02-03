package effects

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/triggers"
)

type MultiEffectDescription struct {
	Effects []EffectDescription
	EffectDescription
}

type MultiEffect struct {
	*Effect
}

func (e *MultiEffect) ToDescription() *MultiEffectDescription {
	return nil
}

func (e *MultiEffect) SetState(state triggers.State) {

}

func NewMultiEffect(desc MultiEffectDescription, container container.Container) *MultiEffect {
	return nil
}
