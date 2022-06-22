package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
)

type IntegerPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.IntegerSchema
	Observable bool `json:"observable,omitempty" wot:"withDefault"` //with default
	Value      any  `json:"value,omitempty" wot:"optional"`
}
