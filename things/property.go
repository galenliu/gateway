package things

import (
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	json "github.com/json-iterator/go"
)

type Property struct {
	wot.PropertyAffordance
	Value   interface{}
	Name    string `json:"name"`
	ThingId string `json:"thingId"`
}

func NewPropertyFromString(description string) *Property {
	bt := []byte(description)
	var property = Property{}
	property.PropertyAffordance = wot.NewPropertyAffordanceFromString(description)
	if n := json.Get(bt, "name").ToString(); n != "" {
		property.Name = n
	}

	if tid := json.Get(bt, "thingId").ToString(); tid != "" {
		property.Name = tid
	}
	return &property
}

//// SetCachedValue 设置本地缓存的值
//func (p *Property) SetCachedValue(value interface{}) {
//	value = p.convert(value)
//	p.Value = p.clamp(value)
//}
//
//func (p *Property) ToValue(value interface{}) interface{} {
//	newValue := p.convert(value)
//	newValue = p.convert(newValue)
//	return newValue
//}



