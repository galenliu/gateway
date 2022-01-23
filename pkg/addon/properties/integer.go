package properties

import "math"

type Integer int64

type IntegerEntity interface {
	Entity
	GetMinValue() Integer
	GetMaxValue() Integer
	GetValue() Integer
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

func (prop *IntegerProperty) GetMinValue() Integer {
	if v := prop.GetMinimum(); v != nil {
		f, ok := v.(Integer)
		if ok {
			return f
		}
	}
	return math.MinInt64
}

func (prop *IntegerProperty) GetMaxValue() Integer {
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
