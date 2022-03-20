package property_affordance

import (
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	"github.com/goccy/go-json"
)

type BooleanPropertyAffordance struct {
	*ia.InteractionAffordance
	*schema.BooleanSchema
	Observable bool `json:"observable,omitempty"`
	Value      any  `json:"value,omitempty" wot:"optional"`
}

type BooleanPropertyDescription struct {
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

	Observable bool `json:"observable,omitempty"`
}

func (p *BooleanPropertyAffordance) UnmarshalJSON(data []byte) error {
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
	p.BooleanSchema = &schema.BooleanSchema{DataSchema: &schema.DataSchema{
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
	}}
	p.Observable = prop.Observable
	return nil
}

func (p *BooleanPropertyAffordance) MarshalJSON() ([]byte, error) {

	prop := BooleanPropertyDescription{
		AtType:       p.AtType,
		Title:        p.Title,
		Titles:       p.Titles,
		Description:  p.Description,
		Descriptions: p.Descriptions,
		Forms:        p.Forms,
		UriVariables: p.UriVariables,
		Const:        p.Const,
		Default:      p.Default,
		Unit:         p.Unit,
		OneOf:        p.OneOf,
		Enum:         p.Enum,
		ReadOnly:     p.ReadOnly,
		WriteOnly:    p.WriteOnly,
		Format:       p.Format,
		Type:         p.Type,
		Observable:   p.Observable,
	}
	return json.Marshal(prop)
}
