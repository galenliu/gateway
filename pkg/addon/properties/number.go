package properties

import (
	"fmt"
	"github.com/xiam/to"
	"math"
)

type NumberPropertyDescription struct {
	Name       string       `json:"name,omitempty"`
	AtType     PropertyType `json:"@type,omitempty"`
	Minimum    Number       `json:"minimum,omitempty"`
	Maximum    Number       `json:"maximum,omitempty"`
	Enum       []Number     `json:"enum,omitempty"`
	ReadOnly   bool         `json:"readOnly,omitempty"`
	MultipleOf Number       `json:"multipleOf,omitempty"`
	Value      Number       `json:"value,omitempty"`
}

type Number float64

type NumberEntity interface {
	Entity
	CheckValue(v any) Number
	SetValue(Number) error
}

type NumberProperty struct {
	*Property
}

func NewNumberProperty(desc NumberPropertyDescription, opts ...Option) *NumberProperty {
	n := &NumberProperty{}
	n.Property = NewProperty(PropertyDescription{
		Name:    desc.Name,
		AtType:  desc.AtType,
		Type:    TypeNumber,
		Minimum: desc.Minimum,
		Maximum: desc.Maximum,
		Enum: func() []any {
			enum := make([]any, 0)
			for _, e := range desc.Enum {
				enum = append(enum, e)
			}
			return enum
		}(),
		ReadOnly:   desc.ReadOnly,
		MultipleOf: desc.MultipleOf,
		Value:      desc.Value,
	}, opts...)
	return n
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *NumberProperty) OnValueRemoteGet(fn func() float64) {
	//prop.OnValueGet(func() interface{} {
	//	return fn()
	//})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *NumberProperty) OnValueRemoteUpdate(fn func(float64)) {
	//prop.OnValueUpdate(func(Property *addon.PropertyProxy, newValue, oldValue interface{}) {
	//	fn(newValue.(float64))
	//})
}

func (prop *NumberProperty) GetValue() Number {
	v := prop.Value.(Number)
	return v
}

func (prop *NumberProperty) CheckValue(v any) Number {
	f := to.Float64(v)
	return prop.clamp(Number(f))
}

func (prop *NumberProperty) clamp(v Number) Number {
	value := v
	if max := prop.getMaxValue(); max < value {
		value = max
	}
	if min := prop.getMinValue(); min > value {
		value = min
	}
	return value
}

func (prop *NumberProperty) getMaxValue() Number {
	if v := prop.GetMaximum(); v != nil {
		f, ok := v.(Number)
		if ok {
			return f
		}
	}
	return math.MaxFloat64
}

func (prop *NumberProperty) getMinValue() Number {
	if v := prop.GetMinimum(); v != nil {
		n, ok := v.(Number)
		if ok {
			return n
		}
	}
	return math.MinInt64
}

func (prop *NumberProperty) SetValue(v Number) error {
	return fmt.Errorf("device:%s property:%s set value:%v not implemented ", prop.GetDevice().GetId(), prop.GetName(), v)
}
