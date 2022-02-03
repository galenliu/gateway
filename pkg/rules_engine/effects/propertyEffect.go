package effects

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/property"
)

type PropertyEffectDescription struct {
	EffectDescription
	property.PropertyDescription
}

type PropertyEffect struct {
	*property.Property
}

func NewPropertyEffect(des PropertyEffectDescription, container container.Container) *PropertyEffect {
	e := &PropertyEffect{property.NewProperty(des.PropertyDescription, container)}
	return e
}

func (e *PropertyEffect) ToDescription() PropertyEffectDescription {
	des := PropertyEffectDescription{}
	des.PropertyDescription = e.Property.ToDescription()
	return des
}
