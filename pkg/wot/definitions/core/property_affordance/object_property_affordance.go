package property_affordance

import (
	"github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type ObjectPropertyAffordance struct {
	*interaction_affordance.InteractionAffordance
	*schema.ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

func (p ObjectPropertyAffordance) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&p, "", "   ")
}

func (p ObjectPropertyAffordance) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p)
}

func NewObjectPropertyAffordanceFormString(description string) *ObjectPropertyAffordance {
	data := []byte(description)
	var p = ObjectPropertyAffordance{}
	p.InteractionAffordance = interaction_affordance.NewInteractionAffordanceFromString(description)
	p.ObjectSchema = schema.NewObjectSchemaFromString(description)
	if p.ObjectSchema == nil {
		return nil
	}
	p.Observable = controls.JSONGetBool(data, "observable", false)
	return &p
}
