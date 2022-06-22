package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
)

type NullPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NullSchema
	Observable bool `json:"observable"`
	Value      any  `json:"value,omitempty" wot:"optional"`
}
