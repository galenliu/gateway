package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type IntegerPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.IntegerSchema
	Observable bool `json:"observable,omitempty" wot:"withDefault"` //with default
}

type IntegerPropertyDescription struct {
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

	Minimum          *controls.Integer `json:"minimum,,omitempty"`
	ExclusiveMinimum *controls.Integer `json:"exclusiveMinimum,omitempty"`
	Maximum          *controls.Integer `json:"maximum,omitempty"`
	ExclusiveMaximum *controls.Integer `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *controls.Integer `json:"multipleOf,omitempty"`

	Observable bool `json:"observable,omitempty"`
}

func (p *IntegerPropertyAffordance) UnmarshalJSON(data []byte) error {

	var prop IntegerPropertyDescription
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
	var s schema.IntegerSchema
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	p.IntegerSchema = &schema.IntegerSchema{
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

func (p *IntegerPropertyAffordance) MarshalJSON() ([]byte, error) {
	prop := IntegerPropertyDescription{
		AtType:           p.AtType,
		Title:            p.Title,
		Titles:           p.Titles,
		Description:      p.Description,
		Descriptions:     p.Descriptions,
		Forms:            p.Forms,
		UriVariables:     p.UriVariables,
		Minimum:          p.Minimum,
		ExclusiveMinimum: p.ExclusiveMinimum,
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
