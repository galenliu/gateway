package core

import (
	"fmt"
	"github.com/galenliu/gateway/wot/definitions/data_schema"
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
)

type PropertyAffordance interface {
	MarshalJSON() ([]byte, error)
}

type propertyAffordance struct {
	*InteractionAffordance
	data_schema.DataSchema
	Observable bool        `json:"observable,omitempty"`
	Value      interface{} `json:"value,omitempty"`
}

func NewPropertyAffordanceFromString(description string) PropertyAffordance {
	data := []byte(description)
	var p = propertyAffordance{}
	p.InteractionAffordance = NewInteractionAffordanceFromString(description)
	typ := json.Get(data, "type").ToString()
	switch typ {
	case hypermedia_controls.Boolean:
		p.DataSchema = data_schema.NewBooleanSchemaFromString(description)
	case hypermedia_controls.Integer:
		p.DataSchema = data_schema.NewIntegerSchemaFromString(description)
	case hypermedia_controls.Number:
		p.DataSchema = data_schema.NewNumberSchemaFromString(description)
	case hypermedia_controls.Array:
		p.DataSchema = data_schema.NewArraySchemaFromString(description)
	case hypermedia_controls.String:
		p.DataSchema = data_schema.NewStringSchemaFromString(description)
	case hypermedia_controls.Null:
		p.DataSchema = data_schema.NewNullSchemaFromString(description)
	case hypermedia_controls.Object:
		p.DataSchema = data_schema.NewObjectSchemaFromString(description)
	}
	vt := json.Get(data, "observable").ValueType()
	if vt == json.BoolValue {
		p.Observable = json.Get(data, "observable").ToBool()
	}
	p.Value = json.Get(data, "value").GetInterface()

	return p
}

func (p *propertyAffordance) GetDescription() string {
	return p.InteractionAffordance.Description
}

func (p propertyAffordance) MarshalJSON() ([]byte, error) {

	if p.DataSchema == nil {
		return nil, fmt.Errorf("dataschema nil")
	}
	return json.MarshalIndent(&p, "", " ")

	//switch p.DataSchema.(type) {
	//case *data_schema.ArraySchema:
	//	d := p.DataSchema.(*data_schema.ArraySchema)
	//	var pa = ArrayPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.ArraySchema:
	//	d := p.DataSchema.(data_schema.ArraySchema)
	//	var pa = ArrayPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case *data_schema.BooleanSchema:
	//	d := p.DataSchema.(*data_schema.BooleanSchema)
	//	var pa = BooleanPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Bool(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.BooleanSchema:
	//	d := p.DataSchema.(data_schema.BooleanSchema)
	//	var pa = BooleanPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Bool(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case *data_schema.NumberSchema:
	//	d := p.DataSchema.(*data_schema.NumberSchema)
	//	var pa = NumberPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Float64(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.NumberSchema:
	//	d := p.DataSchema.(data_schema.NumberSchema)
	//	var pa = NumberPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Float64(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case *data_schema.IntegerSchema:
	//	d := p.DataSchema.(*data_schema.IntegerSchema)
	//	var pa = IntegerPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Int64(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.IntegerSchema:
	//	d := p.DataSchema.(data_schema.IntegerSchema)
	//	var pa = IntegerPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Int64(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case *data_schema.ObjectSchema:
	//	d := p.DataSchema.(*data_schema.ObjectSchema)
	//	var pa = ObjectPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.ObjectSchema:
	//	d := p.DataSchema.(data_schema.ObjectSchema)
	//	var pa = ObjectPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case *data_schema.StringSchema:
	//	d := p.DataSchema.(*data_schema.StringSchema)
	//	var pa = StringPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.String(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.StringSchema:
	//	d := p.DataSchema.(data_schema.StringSchema)
	//	var pa = StringPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.String(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case *data_schema.NullSchema:
	//	d := p.DataSchema.(*data_schema.NullSchema)
	//	var pa = NullPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.NullSchema:
	//	d := p.DataSchema.(data_schema.NullSchema)
	//	var pa = NullPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//default:
	//	return nil, fmt.Errorf("property type err")
	//}
}

// GetDefaultValue 获取默认的值
func (p *propertyAffordance) GetDefaultValue() interface{} {
	switch p.DataSchema.GetType() {
	case hypermedia_controls.Number:
		d := p.DataSchema.(*data_schema.NumberSchema)
		return d.Minimum
	case hypermedia_controls.Integer:
		d := p.DataSchema.(*data_schema.IntegerSchema)
		return d.Minimum
	case hypermedia_controls.Boolean:
		return false
	case hypermedia_controls.String:
		return ""
	case hypermedia_controls.Array:
		return []interface{}{}
	case hypermedia_controls.Null:
		return nil
	default:
		return nil
	}
}

// SetMaxValue 设置最大值
func (p *propertyAffordance) SetMaxValue(v interface{}) {
	switch p.DataSchema.(type) {
	case *data_schema.NumberSchema:
		d := p.DataSchema.(*data_schema.NumberSchema)
		d.Maximum = to.Float64(v)
	case data_schema.NumberSchema:
		d := p.DataSchema.(*data_schema.NumberSchema)
		d.Maximum = to.Float64(v)
	case *data_schema.IntegerSchema:
		d := p.DataSchema.(*data_schema.IntegerSchema)
		d.Maximum = to.Int64(v)
	case data_schema.IntegerSchema:
		d := p.DataSchema.(*data_schema.IntegerSchema)
		d.Maximum = to.Int64(v)
	default:
		fmt.Print("property type err")
		return
	}
}

// SetMinValue 设置最小值
func (p *propertyAffordance) SetMinValue(v interface{}) {
	switch p.DataSchema.(type) {
	case *data_schema.NumberSchema:
		d := p.DataSchema.(*data_schema.NumberSchema)
		d.Minimum = to.Float64(v)
	case data_schema.NumberSchema:
		d := p.DataSchema.(*data_schema.NumberSchema)
		d.Minimum = to.Float64(v)
	case *data_schema.IntegerSchema:
		d := p.DataSchema.(*data_schema.IntegerSchema)
		d.Minimum = to.Int64(v)
	case data_schema.IntegerSchema:
		d := p.DataSchema.(*data_schema.IntegerSchema)
		d.Minimum = to.Int64(v)

	default:
		fmt.Print("property type err")
		return
	}
}

type ArrayPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.ArraySchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type BooleanPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

type NumberPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

type IntegerPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.IntegerSchema
	Observable bool  `json:"observable"`
	Value      int64 `json:"value"`
}

type ObjectPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type StringPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.StringSchema
	Observable bool   `json:"observable"`
	Value      string `json:"value"`
}

type NullPropertyAffordance struct {
	*InteractionAffordance
	*data_schema.NullSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}
