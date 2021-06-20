package core

import (
	"fmt"
	"github.com/galenliu/gateway/wot/definitions/data_schema"
	"github.com/galenliu/gateway/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/xiam/to"
)

type PropertyAffordance interface {
}

type propertyAffordance struct {
	*InteractionAffordance
	data_schema.DataSchemaInterface
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value,omitempty"`
}

func NewPropertyAffordanceFromString(description string) PropertyAffordance {
	var p = propertyAffordance{}
	p.InteractionAffordance = NewInteractionAffordanceFromString(description)
	p.DataSchemaInterface = data_schema.NewDataSchemaFromString(description)
	if gjson.Get(description, "observable").Exists() {
		p.Observable = gjson.Get(description, "observable").Bool()
	}
	return p
}

// SetCachedValue 设置本地缓存的值
func (p *propertyAffordance) SetCachedValue(value interface{}) {
	value = p.convert(value)
	p.Value = p.clamp(value)
}

func (p *propertyAffordance) ToValue(value interface{}) interface{} {
	newValue := p.convert(value)
	newValue = p.convert(newValue)
	return newValue
}

//确保属性值相应的类型
func (p *propertyAffordance) convert(v interface{}) interface{} {
	switch p.GetType() {
	case hypermedia_controls.Number:
		return to.Float64(v)
	case hypermedia_controls.Integer:
		return int(to.Uint64(v))
	case hypermedia_controls.Boolean:
		return to.Bool(v)
	default:
		return v
	}
}

//确保属性值在允许的范围
func (p *propertyAffordance) clamp(v interface{}) interface{} {
	switch p.DataSchemaInterface.(type) {
	case *data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		return d.ClampFloat(to.Float64(v))
	case data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		return d.ClampFloat(to.Float64(v))
	case *data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		return d.ClampInt(to.Int64(v))
	case data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		return d.ClampInt(to.Int64(v))
	default:
		return v
	}

}

func (p propertyAffordance) MarshalJSON() ([]byte, error) {

	if p.DataSchemaInterface == nil {
		return nil, fmt.Errorf("dataschema err")
	}
	switch p.DataSchemaInterface.(type) {
	case *data_schema.ArraySchema:
		d := p.DataSchemaInterface.(*data_schema.ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.ArraySchema:
		d := p.DataSchemaInterface.(data_schema.ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *data_schema.BooleanSchema:
		d := p.DataSchemaInterface.(*data_schema.BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Bool(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.BooleanSchema:
		d := p.DataSchemaInterface.(data_schema.BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Bool(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Float64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.NumberSchema:
		d := p.DataSchemaInterface.(data_schema.NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Float64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Int64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(data_schema.IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Int64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *data_schema.ObjectSchema:
		d := p.DataSchemaInterface.(*data_schema.ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.ObjectSchema:
		d := p.DataSchemaInterface.(data_schema.ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *data_schema.StringSchema:
		d := p.DataSchemaInterface.(*data_schema.StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.String(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.StringSchema:
		d := p.DataSchemaInterface.(data_schema.StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.String(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *data_schema.NullSchema:
		d := p.DataSchemaInterface.(*data_schema.NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case data_schema.NullSchema:
		d := p.DataSchemaInterface.(data_schema.NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	default:
		return nil, fmt.Errorf("property type err")
	}
}

// GetDefaultValue 获取默认的值
func (p *propertyAffordance) GetDefaultValue() interface{} {
	switch p.DataSchemaInterface.(type) {
	case *data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		return d.Minimum
	case data_schema.NumberSchema:
		d := p.DataSchemaInterface.(data_schema.NumberSchema)
		return d.Minimum
	case *data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		return d.Minimum
	case data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(data_schema.IntegerSchema)
		return d.Minimum

	case *data_schema.BooleanSchema:
		return false
	case data_schema.BooleanSchema:
		return false

	case *data_schema.StringSchema:
		return ""
	case data_schema.StringSchema:
		return ""

	default:
		return nil
	}
}

// SetMaxValue 设置最大值
func (p *propertyAffordance) SetMaxValue(v interface{}) {
	switch p.DataSchemaInterface.(type) {
	case *data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		d.Maximum = to.Float64(v)
	case data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		d.Maximum = to.Float64(v)
	case *data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		d.Maximum = to.Int64(v)
	case data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		d.Maximum = to.Int64(v)
	default:
		fmt.Print("property type err")
		return
	}
}

// SetMinValue 设置最小值
func (p *propertyAffordance) SetMinValue(v interface{}) {
	switch p.DataSchemaInterface.(type) {
	case *data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		d.Minimum = to.Float64(v)
	case data_schema.NumberSchema:
		d := p.DataSchemaInterface.(*data_schema.NumberSchema)
		d.Minimum = to.Float64(v)
	case *data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
		d.Minimum = to.Int64(v)
	case data_schema.IntegerSchema:
		d := p.DataSchemaInterface.(*data_schema.IntegerSchema)
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
