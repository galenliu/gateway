package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
)

type NumberPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NumberSchema
	Observable bool `json:"observable,omitempty"`
	Value      any  `json:"value,omitempty" wot:"optional"`
}
