package effects

import (
	"github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/rules_engine/property"
	json "github.com/json-iterator/go"
)

type PropertyEffectDescription struct {
	EffectDescription
	Property property.Description `json:"property"`
}

type PropertyEffect struct {
	*Effect
	property *property.Property
}

func (e *PropertyEffect) ToDescription() PropertyEffectDescription {
	des := PropertyEffectDescription{
		EffectDescription: e.Effect.ToDescription(),
		Property:          e.property.ToDescription(),
	}
	return des
}

func (e *PropertyEffect) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToDescription())
}

func NewPropertyEffect(des PropertyEffectDescription, container container.Container) *PropertyEffect {
	e := &PropertyEffect{NewEffect(des.EffectDescription), property.NewProperty(des.Property, container)}
	return e
}
