package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type NumberPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.NumberSchema
	Observable bool `json:"observable,omitempty"`
}

type NumberPropertyDescription struct {
	AtType       string                       `json:"@type,omitempty,optional"`
	Title        string                       `json:"title,omitempty,optional"`
	Titles       map[string]string            `json:"titles,omitempty,optional"`
	Description  string                       `json:"description,omitempty,optional"`
	Descriptions map[string]string            `json:"descriptions,omitempty,optional"`
	Forms        []controls.Form              `json:"forms,omitempty,mandatory" wot:"optional"`
	UriVariables map[string]schema.DataSchema `json:"uriVariables,omitempty"`

	Const     any                 `json:"const,omitempty" wot:"optional"`
	Default   any                 `json:"default,omitempty" wot:"optional"`
	Unit      string              `json:"unit,omitempty" wot:"optional"`
	OneOf     []schema.DataSchema `json:"oneOf,,omitempty" wot:"optional"`
	Enum      []any               `json:"enum,omitempty" wot:"optional"`
	ReadOnly  bool                `json:"readOnly,omitempty" wot:"withDefault"`
	WriteOnly bool                `json:"writeOnly,omitempty" wot:"withDefault"`
	Format    string              `json:"format,omitempty" wot:"optional"`
	Type      string              `json:"type,,omitempty" wot:"optional"`

	Minimum          *controls.Double `json:"minimum,omitempty"`
	ExclusiveMinimum *controls.Double `json:"exclusiveMinimum,omitempty"`
	Maximum          *controls.Double `json:"maximum,omitempty"`
	ExclusiveMaximum *controls.Double `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *controls.Double `json:"multipleOf,omitempty"`

	Observable bool `json:"observable,omitempty"`
}

func (p *NumberPropertyAffordance) UnmarshalJSON(data []byte) error {
	var prop NumberPropertyDescription
	err := json.Unmarshal(data, &prop)
	if err != nil {
		return err
	}
	p.InteractionAffordance = &ia.InteractionAffordance{
		Type:         prop.Type,
		Title:        prop.Title,
		Titles:       prop.Titles,
		Description:  prop.Description,
		Descriptions: prop.Descriptions,
		Forms:        prop.Forms,
		UriVariables: prop.UriVariables,
	}
	p.NumberSchema = &schema.NumberSchema{
		DataSchema: &schema.DataSchema{
			AtType:       prop.AtType,
			Title:        prop.Title,
			Titles:       prop.Titles,
			Description:  prop.Description,
			Descriptions: prop.Descriptions,
			Const:        prop.Const,
			Default:      prop.Default,
			Unit:         prop.Unit,
			OneOf:        prop.OneOf,
			Enum:         prop.Enum,
			ReadOnly:     prop.ReadOnly,
			WriteOnly:    prop.WriteOnly,
			Format:       prop.Format,
			Type:         prop.Type,
		},
		Minimum:          prop.Minimum,
		ExclusiveMinimum: prop.ExclusiveMinimum,
		Maximum:          prop.Maximum,
		ExclusiveMaximum: prop.ExclusiveMaximum,
		MultipleOf:       prop.MultipleOf,
	}
	p.Observable = prop.Observable
	return nil
}

func (p *NumberPropertyAffordance) MarshalJSON() ([]byte, error) {
	prop := NumberPropertyDescription{
		AtType:           p.AtType,
		Title:            p.Title,
		Titles:           p.Titles,
		Description:      p.Description,
		Descriptions:     p.Descriptions,
		Forms:            p.Forms,
		UriVariables:     p.UriVariables,
		Minimum:          p.Minimum,
		ExclusiveMinimum: p.Maximum,
		Maximum:          p.Maximum,
		ExclusiveMaximum: p.ExclusiveMaximum,
		MultipleOf:       p.MultipleOf,
		Const:            p.Const,
		Default:          p.Default,
		Unit:             p.Unit,
		OneOf:            p.OneOf,
		Enum:             p.Enum,
		ReadOnly:         p.ReadOnly,
		WriteOnly:        p.WriteOnly,
		Format:           p.Format,
		Type:             p.Type,
		Observable:       p.Observable,
	}
	return json.Marshal(prop)
}
