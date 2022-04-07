package property_affordance

import (
	"fmt"
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type PropertySchema interface {
	GetType() controls.DataSchemaType
	IsReadOnly() bool
}

type PropertyAffordance struct {
	*ia.InteractionAffordance
	PropertySchema
	Observable bool `json:"observable,omitempty"`
	Value      any  `json:"value,omitempty" wot:"optional"`
}

func (p *PropertyAffordance) UnmarshalJSON(data []byte) error {
	var prop ArrayPropertyDescription
	err := json.Unmarshal(data, &prop)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &ia.InteractionAffordance{
		AtType:       prop.AtType,
		Title:        prop.Title,
		Titles:       prop.Titles,
		Description:  prop.Description,
		Descriptions: prop.Descriptions,
		Forms:        prop.Forms,
		UriVariables: prop.UriVariables,
	}
	dataType := json.Get(data, "type").ToString()
	switch dataType {
	case controls.TypeInteger:
		var dataSchema schema.IntegerSchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	case controls.TypeNumber:
		var dataSchema schema.NumberSchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	case controls.TypeBoolean:
		var dataSchema schema.BooleanSchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	case controls.TypeArray:
		var dataSchema schema.ArraySchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	case controls.TypeObject:
		var dataSchema schema.ObjectSchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	case controls.TypeNull:
		var dataSchema schema.NullSchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	case controls.TypeString:
		var dataSchema schema.StringSchema
		err := json.Unmarshal(data, &dataSchema)
		if err != nil {
			return err
		}
		p.PropertySchema = dataSchema
	default:
		return fmt.Errorf("unsupported type: %s", dataType)

	}
	p.Observable = prop.Observable
	return nil
}

func (p PropertyAffordance) MarshalJSON() ([]byte, error) {

	switch p.PropertySchema.(type) {
	case schema.NumberSchema:
		dataSchema, ok := p.PropertySchema.(schema.NumberSchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.NumberSchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	case schema.IntegerSchema:
		dataSchema, ok := p.PropertySchema.(schema.IntegerSchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.IntegerSchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	case schema.StringSchema:
		dataSchema, ok := p.PropertySchema.(schema.StringSchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.StringSchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	case schema.BooleanSchema:
		dataSchema, ok := p.PropertySchema.(schema.BooleanSchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.BooleanSchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	case schema.ObjectSchema:
		dataSchema, ok := p.PropertySchema.(schema.ObjectSchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.ObjectSchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	case schema.ArraySchema:
		dataSchema, ok := p.PropertySchema.(schema.ArraySchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.ArraySchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	case schema.NullSchema:
		dataSchema, ok := p.PropertySchema.(schema.NullSchema)
		if !ok {
			return nil, fmt.Errorf("type error")
		}
		return json.Marshal(struct {
			*ia.InteractionAffordance
			schema.NullSchema
			Observable bool `json:"observable,omitempty"`
		}{p.InteractionAffordance, dataSchema, p.Observable})

	default:
		return nil, fmt.Errorf("type error")

	}

}
