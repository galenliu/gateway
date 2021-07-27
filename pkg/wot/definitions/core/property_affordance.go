package core

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core/property_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

type PropertyAffordance interface {
	schema.DataSchema
}

func NewPropertyAffordanceFromString(description string) PropertyAffordance {
	data := []byte(description)
	switch controls.JSONGetString(data, "type", "") {
	case controls.TypeBoolean:
		return property_affordance.NewBooleanPropertyAffordanceFormString(description)
	case controls.TypeInteger:
		return property_affordance.NewIntegerPropertyAffordanceFormString(description)
	case controls.TypeNumber:
		return property_affordance.NewNumberPropertyAffordanceFormString(description)
	case controls.TypeArray:
		return property_affordance.NewArrayPropertyAffordanceFormString(description)
	case controls.TypeString:
		return property_affordance.NewStringPropertyAffordanceFormString(description)
	case controls.TypeNull:
		return property_affordance.NewNullPropertyAffordanceFormString(description)
	case controls.TypeObject:
		return property_affordance.NewObjectPropertyAffordanceFormString(description)
	default:
		return nil
	}
}
