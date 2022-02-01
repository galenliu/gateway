package effects

import "github.com/galenliu/gateway/pkg/rules_engine/property"

type PropertyEffectDescription struct {
	property.PropertyDescription
}

type PropertyEffect struct {
	*property.Property
}

func NewPropertyEffect(des PropertyEffectDescription, bus property.Bus, things property.ThingsHandler) *PropertyEffect {
	e := &PropertyEffect{property.NewProperty(des.PropertyDescription, things)}
	return e
}

func (e *PropertyEffect) ToDescription() PropertyEffectDescription {
	des := PropertyEffectDescription{}
	des.PropertyDescription = e.Property.ToDescription()
	return des
}
