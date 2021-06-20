package TD

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/xiam/to"
)

type PropertyAffordance struct {
	*InteractionAffordance
	IDataSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value,omitempty"`
}

func NewPropertyAffordanceFromString(description string) *PropertyAffordance {
	var p = PropertyAffordance{}
	p.InteractionAffordance = NewInteractionAffordanceFromString(description)
	p.IDataSchema = NewDataSchemaFromString(description)
	if gjson.Get(description, "observable").Exists() {
		p.Observable = gjson.Get(description, "observable").Bool()
	}
	return &p
}

// SetCachedValue 设置本地缓存的值
func (p *PropertyAffordance) SetCachedValue(value interface{}) {
	value = p.convert(value)
	p.Value = p.clamp(value)
}

func (p *PropertyAffordance) ToValue(value interface{}) interface{} {
	newValue := p.convert(value)
	newValue = p.convert(newValue)
	return newValue
}

//确保属性值相应的类型
func (p *PropertyAffordance) convert(v interface{}) interface{} {
	switch p.GetType() {
	case Number:
		return to.Float64(v)
	case Integer:
		return int(to.Uint64(v))
	case Boolean:
		return to.Bool(v)
	default:
		return v
	}
}

//确保属性值在允许的范围
func (p *PropertyAffordance) clamp(v interface{}) interface{} {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		return d.ClampFloat(to.Float64(v))
	case NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		return d.ClampFloat(to.Float64(v))
	case *IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		return d.ClampInt(to.Int64(v))
	case IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		return d.ClampInt(to.Int64(v))
	default:
		return v
	}

}

func (p PropertyAffordance) MarshalJSON() ([]byte, error) {

	if p.IDataSchema == nil {
		return nil, fmt.Errorf("dataschema err")
	}
	switch p.IDataSchema.(type) {
	case *ArraySchema:
		d := p.IDataSchema.(*ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case ArraySchema:
		d := p.IDataSchema.(ArraySchema)
		var pa = ArrayPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *BooleanSchema:
		d := p.IDataSchema.(*BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Bool(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case BooleanSchema:
		d := p.IDataSchema.(BooleanSchema)
		var pa = BooleanPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Bool(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Float64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case NumberSchema:
		d := p.IDataSchema.(NumberSchema)
		var pa = NumberPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Float64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.Int64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case IntegerSchema:
		d := p.IDataSchema.(IntegerSchema)
		var pa = IntegerPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.Int64(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *ObjectSchema:
		d := p.IDataSchema.(*ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case ObjectSchema:
		d := p.IDataSchema.(ObjectSchema)
		var pa = ObjectPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *StringSchema:
		d := p.IDataSchema.(*StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, d, p.Observable, to.String(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case StringSchema:
		d := p.IDataSchema.(StringSchema)
		var pa = StringPropertyAffordance{p.InteractionAffordance, &d, p.Observable, to.String(p.Value)}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case *NullSchema:
		d := p.IDataSchema.(*NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	case NullSchema:
		d := p.IDataSchema.(NullSchema)
		var pa = NullPropertyAffordance{p.InteractionAffordance, &d, p.Observable, p.Value}
		pa.AtType = d.AtType
		return json.MarshalIndent(pa, "", "  ")
	default:
		return nil, fmt.Errorf("property type err")
	}
}

// GetDefaultValue 获取默认的值
func (p *PropertyAffordance) GetDefaultValue() interface{} {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		return d.Minimum
	case NumberSchema:
		d := p.IDataSchema.(NumberSchema)
		return d.Minimum
	case *IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		return d.Minimum
	case IntegerSchema:
		d := p.IDataSchema.(IntegerSchema)
		return d.Minimum

	case *BooleanSchema:
		return false
	case BooleanSchema:
		return false

	case *StringSchema:
		return ""
	case StringSchema:
		return ""

	default:
		return nil
	}
}

// SetMaxValue 设置最大值
func (p *PropertyAffordance) SetMaxValue(v interface{}) {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		d.Maximum = to.Float64(v)
	case NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		d.Maximum = to.Float64(v)
	case *IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		d.Maximum = to.Int64(v)
	case IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		d.Maximum = to.Int64(v)
	default:
		fmt.Print("property type err")
		return
	}
}

// SetMinValue 设置最小值
func (p *PropertyAffordance) SetMinValue(v interface{}) {
	switch p.IDataSchema.(type) {
	case *NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		d.Minimum = to.Float64(v)
	case NumberSchema:
		d := p.IDataSchema.(*NumberSchema)
		d.Minimum = to.Float64(v)
	case *IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		d.Minimum = to.Int64(v)
	case IntegerSchema:
		d := p.IDataSchema.(*IntegerSchema)
		d.Minimum = to.Int64(v)

	default:
		fmt.Print("property type err")
		return
	}
}

type ArrayPropertyAffordance struct {
	*InteractionAffordance
	*ArraySchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type BooleanPropertyAffordance struct {
	*InteractionAffordance
	*BooleanSchema
	Observable bool `json:"observable"`
	Value      bool `json:"value"`
}

type NumberPropertyAffordance struct {
	*InteractionAffordance
	*NumberSchema
	Observable bool    `json:"observable"`
	Value      float64 `json:"value"`
}

type IntegerPropertyAffordance struct {
	*InteractionAffordance
	*IntegerSchema
	Observable bool  `json:"observable"`
	Value      int64 `json:"value"`
}

type ObjectPropertyAffordance struct {
	*InteractionAffordance
	*ObjectSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}

type StringPropertyAffordance struct {
	*InteractionAffordance
	*StringSchema
	Observable bool   `json:"observable"`
	Value      string `json:"value"`
}

type NullPropertyAffordance struct {
	*InteractionAffordance
	*NullSchema
	Observable bool        `json:"observable"`
	Value      interface{} `json:"value"`
}
