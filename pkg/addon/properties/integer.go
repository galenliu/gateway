package properties

import (
	"fmt"
	"github.com/xiam/to"
	"math"
)

type Integer int64

type IntegerEntity interface {
	Entity
	GetValue() Integer
	CheckValue(any) Integer
	SetValue(v Integer) error
}

type IntegerProperty struct {
	*Property
}

func NewIntegerProperty(description PropertyDescription) *IntegerProperty {
	i := &IntegerProperty{}
	description.Type = TypeInteger
	i.Property = NewProperty(description)
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
