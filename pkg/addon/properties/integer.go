package properties

import (
	"fmt"
	"github.com/xiam/to"
	"math"
)

type Integer int64

type IntegerPropertyDescription struct {
	Name        string             `json:"name,omitempty"`
	AtType      string             `json:"@type,omitempty"`
	Title       string             `json:"title,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Description string             `json:"description,omitempty"`
	Minimum     Integer            `json:"minimum,omitempty"`
	Maximum     Integer            `json:"maximum,omitempty"`
	Enum        []Integer          `json:"enum,omitempty"`
	ReadOnly    bool               `json:"readOnly,omitempty"`
	MultipleOf  any                `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       Integer            `json:"value,omitempty"`
}

type IntegerEntity interface {
	Entity
	GetValue() Integer
	CheckValue(any) Integer
	SetValue(v Integer) error
}

type IntegerProperty struct {
	*Property
}

func NewIntegerProperty(desc IntegerPropertyDescription) *IntegerProperty {
	i := &IntegerProperty{}
	i.Property = NewProperty(PropertyDescription{
		Name:        desc.Name,
		AtType:      desc.AtType,
		Title:       desc.Title,
		Type:        TypeInteger,
		Unit:        desc.Unit,
		Description: desc.Description,
		Minimum:     desc.Minimum,
		Maximum:     desc.Maximum,
		Enum: func() []any {
			enum := make([]any, 0)
			for _, e := range desc.Enum {
				enum = append(enum, e)
			}
			return enum
		}(),
		ReadOnly:   desc.ReadOnly,
		MultipleOf: desc.MultipleOf,
		Links:      desc.Links,
		Value:      desc.Value,
	})
	return i
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *IntegerProperty) OnValueRemoteGet(fn func() int) {

}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *IntegerProperty) OnValueRemoteUpdate(fn func(int)) {
	//prop.OnValueUpdate(func(Property *addon.PropertyProxy, newValue, oldValue interface{}) {
	//	fn(newValue.(int))
	//})
}

func (prop *IntegerProperty) CheckValue(v any) Integer {
	value := to.Int64(v)
	return prop.clamp(Integer(value))
}

func (prop *IntegerProperty) getMinValue() Integer {
	if v := prop.GetMinimum(); v != nil {
		f, ok := v.(Integer)
		if ok {
			return f
		}
	}
	return math.MinInt64
}

func (prop *IntegerProperty) getMaxValue() Integer {
	if v := prop.GetMaximum(); v != nil {
		f, ok := v.(Integer)
		if ok {
			return f
		}
	}
	return math.MaxInt64
}

func (prop *IntegerProperty) GetValue() Integer {
	v := prop.Value.(Integer)
	return v
}

func (prop *IntegerProperty) clamp(v Integer) Integer {
	value := v
	if min := prop.GetMinimum(); min != nil {
		minValue, ok := min.(Integer)
		if ok {
			if minValue < value {
				value = minValue
			}
		}
	}
	if max := prop.GetMaximum(); max != nil {
		maxValue, ok := max.(Integer)
		if ok {
			if maxValue < value {
				value = maxValue
			}
		}
	}
	return value
}

func (prop *IntegerProperty) SetValue(v Integer) error {
	return fmt.Errorf("device:%s property:%s set value:%v not implemented ", prop.GetDevice().GetId(), prop.GetName(), v)
}
