package core

import (
	"fmt"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
)

type PropertyAffordance interface {
	MarshalJSON() ([]byte, error)
}

type propertyAffordance struct {
	*InteractionAffordance
	schema.DataSchema
	Observable bool        `json:"observable,omitempty"`
	Value      interface{} `json:"value,omitempty"`
}

func NewPropertyAffordanceFromString(description string) *propertyAffordance {
	data := []byte(description)
	var p = propertyAffordance{}
	p.InteractionAffordance = NewInteractionAffordanceFromString(description)
	typ := json.Get(data, "type").ToString()
	switch typ {
	case controls.TypeBoolean:
		p.DataSchema = schema.NewBooleanSchemaFromString(description)
	case controls.TypeInteger:
		p.DataSchema = schema.NewIntegerSchemaFromString(description)
	case controls.TypeNumber:
		p.DataSchema = schema.NewNumberSchemaFromString(description)
	case controls.TypeArray:
		p.DataSchema = schema.NewArraySchemaFromString(description)
	case controls.TypeString:
		p.DataSchema = schema.NewStringSchemaFromString(description)
	case controls.TypeNull:
		p.DataSchema = schema.NewNullSchemaFromString(description)
	case controls.TypeObject:
		p.DataSchema = schema.NewObjectSchemaFromString(description)
	}
	vt := json.Get(data, "observable").ValueType()
	if vt == json.BoolValue {
		p.Observable = json.Get(data, "observable").ToBool()
	}
	p.Value = json.Get(data, "value").GetInterface()
	return &p
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
	case controls.TypeNumber:
		d := p.DataSchema.(*schema.NumberSchema)
		return d.Minimum
	case controls.TypeInteger:
		d := p.DataSchema.(*schema.IntegerSchema)
		return d.Minimum
	case controls.TypeBoolean:
		return false
	case controls.TypeString:
		return ""
	case controls.TypeArray:
		return []interface{}{}
	case controls.TypeNull:
		return nil
	default:
		return nil
	}
}

// SetMaxValue 设置最大值
//func (p *propertyAffordance) SetMaxValue(v interface{}) {
//	switch p.DataSchema.(type) {
//	case *schema.NumberSchema:
//		d := p.DataSchema.(*schema.NumberSchema)
//		d.Maximum = to.Float64(v)
//	case schema.NumberSchema:
//		d := p.DataSchema.(*schema.NumberSchema)
//		d.Maximum = to.Float64(v)
//	case *schema.IntegerSchema:
//		d := p.DataSchema.(*schema.IntegerSchema)
//		d.Maximum = to.Int64(v)
//	case schema.IntegerSchema:
//		d := p.DataSchema.(*schema.IntegerSchema)
//		d.Maximum = to.Int64(v)
//	default:
//		fmt.Print("property type err")
//		return
//	}
//}

//// SetMinValue 设置最小值
//func (p *propertyAffordance) SetMinValue(v interface{}) {
//	switch p.DataSchema.(type) {
//	case *schema.NumberSchema:
//		d := p.DataSchema.(*schema.NumberSchema)
//		d.Minimum = to.Float64(v)
//	case schema.NumberSchema:
//		d := p.DataSchema.(*schema.NumberSchema)
//		d.Minimum = to.Float64(v)
//	case *schema.IntegerSchema:
//		d := p.DataSchema.(*schema.IntegerSchema)
//		d.Minimum = to.Int64(v)
//	case schema.IntegerSchema:
//		d := p.DataSchema.(*schema.IntegerSchema)
//		d.Minimum = to.Int64(v)
//
//	default:
//		fmt.Print("property type err")
//		return
//	}
//}

type ArrayPropertyAffordance struct {
	*InteractionAffordance
	*schema.ArraySchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type BooleanPropertyAffordance struct {
	*InteractionAffordance
	*schema.BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

type NumberPropertyAffordance struct {
	*InteractionAffordance
	*schema.NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

type IntegerPropertyAffordance struct {
	*InteractionAffordance
	*schema.IntegerSchema
	Observable bool  `json:"observable"`
	Value      int64 `json:"value"`
}

type ObjectPropertyAffordance struct {
	*InteractionAffordance
	*schema.ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type StringPropertyAffordance struct {
	*InteractionAffordance
	*schema.StringSchema
	Observable bool   `json:"observable"`
	Value      string `json:"value"`
}

type NullPropertyAffordance struct {
	*InteractionAffordance
	*schema.NullSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}
