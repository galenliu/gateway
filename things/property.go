package things

import (
	core2 "github.com/galenliu/gateway/pkg/wot/definitions/core"
	data_schema2 "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	hypermedia_controls2 "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
)

type Property struct {
	core2.PropertyAffordance
	Value   interface{}
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewPropertyFromString(description string) *Property {
	bt := []byte(description)
	var property = Property{}
	property.PropertyAffordance = core2.NewPropertyAffordanceFromString(description)
	if n := json.Get(bt, "name").ToString(); n != "" {
		property.Name = n
	}

	if tid := json.Get(bt, "thingId").ToString(); tid != "" {
		property.Name = tid
	}
	return &property
}


// SetCachedValue 设置本地缓存的值
func (p *Property) SetCachedValue(value interface{}) {
	value = p.convert(value)
	p.Value = p.clamp(value)
}

func (p *Property) ToValue(value interface{}) interface{} {
	newValue := p.convert(value)
	newValue = p.convert(newValue)
	return newValue
}

//确保属性值相应的类型
func (p *Property) convert(v interface{}) interface{} {
	switch p.GetType() {
	case hypermedia_controls2.TypeNumber:
		return to.Float64(v)
	case hypermedia_controls2.TypeInteger:
		return int(to.Uint64(v))
	case hypermedia_controls2.TypeBoolean:
		return to.Bool(v)
	default:
		return v
	}
}

//确保属性值在允许的范围
func (p *Property) clamp(v interface{}) interface{} {
	switch p.DataSchema.(type) {
	case *data_schema2.NumberSchema:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		return d.ClampFloat(to.Float64(v))
	case data_schema2.NumberSchema:
		d := p.DataSchema.(*data_schema2.NumberSchema)
		return d.ClampFloat(to.Float64(v))
	case *data_schema2.IntegerSchema:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		return d.ClampInt(to.Int64(v))
	case data_schema2.IntegerSchema:
		d := p.DataSchema.(*data_schema2.IntegerSchema)
		return d.ClampInt(to.Int64(v))
	default:
		return v
	}

}
