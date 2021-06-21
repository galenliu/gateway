package core

import (
	"fmt"
	data_schema2 "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	hypermedia_controls2 "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
)

type PropertyAffordance interface {
	MarshalJSON() ([]byte, error)
}

type propertyAffordance struct {
	*InteractionAffordance
	data_schema2.DataSchema
	Observable bool        `json:"observable,omitempty"`
	Value      interface{} `json:"value,omitempty"`
}

func NewPropertyAffordanceFromString(description string) PropertyAffordance {
	data := []byte(description)
	var p = propertyAffordance{}
	p.InteractionAffordance = NewInteractionAffordanceFromString(description)
	typ := json.Get(data, "type").ToString()
	switch typ {
	case hypermedia_controls2.TypeBoolean:
		p.DataSchema = data_schema2.NewBooleanSchemaFromString(description)
	case hypermedia_controls2.TypeInteger:
		p.DataSchema = data_schema2.NewIntegerSchemaFromString(description)
	case hypermedia_controls2.TypeNumber:
		p.DataSchema = data_schema2.NewNumberSchemaFromString(description)
	case hypermedia_controls2.TypeArray:
		p.DataSchema = data_schema2.NewArraySchemaFromString(description)
	case hypermedia_controls2.TypeString:
		p.DataSchema = data_schema2.NewStringSchemaFromString(description)
	case hypermedia_controls2.TypeNull:
		p.DataSchema = data_schema2.NewNullSchemaFromString(description)
	case hypermedia_controls2.TypeObject:
		p.DataSchema = data_schema2.NewObjectSchemaFromString(description)
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
	//	var pa = StringPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.TypeString(p.Value)}
	//	pa.AtType = d.AtType
	//	return json.MarshalIndent(pa, "", "  ")
	//case data_schema.StringSchema:
	//	d := p.DataSchema.(data_schema.StringSchema)
	//	var pa = StringPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.TypeString(p.Value)}
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
	case hypermedia_controls2.TypeNumber:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		return d.Minimum
	case hypermedia_controls2.TypeInteger:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		return d.Minimum
	case hypermedia_controls2.TypeBoolean:
		return false
	case hypermedia_controls2.TypeString:
		return ""
	case hypermedia_controls2.TypeArray:
		return []interface{}{}
	case hypermedia_controls2.TypeNull:
		return nil
	default:
		return nil
	}
}

// SetMaxValue 设置最大值
func (p *propertyAffordance) SetMaxValue(v interface{}) {
	switch p.DataSchema.(type) {
	case *data_schema2.NumberSchema:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		d.Maximum = to.Float64(v)
	case data_schema2.NumberSchema:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		d.Maximum = to.Float64(v)
	case *data_schema2.IntegerSchema:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		d.Maximum = to.Int64(v)
	case data_schema2.IntegerSchema:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		d.Maximum = to.Int64(v)
	default:
		fmt.Print("property type err")
		return
	}
}

// SetMinValue 设置最小值
func (p *propertyAffordance) SetMinValue(v interface{}) {
	switch p.DataSchema.(type) {
	case *data_schema2.NumberSchema:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		d.Minimum = to.Float64(v)
	case data_schema2.NumberSchema:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		d.Minimum = to.Float64(v)
	case *data_schema2.IntegerSchema:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		d.Minimum = to.Int64(v)
	case data_schema2.IntegerSchema:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		d.Minimum = to.Int64(v)

	default:
		fmt.Print("property type err")
		return
	}
}

type ArrayPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.ArraySchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type BooleanPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

type NumberPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

type IntegerPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.IntegerSchema
	Observable bool  `json:"observable"`
	Value      int64 `json:"value"`
}

type ObjectPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type StringPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.StringSchema
	Observable bool   `json:"observable"`
	Value      string `json:"value"`
}

type NullPropertyAffordance struct {
	*InteractionAffordance
	*data_schema2.NullSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}
