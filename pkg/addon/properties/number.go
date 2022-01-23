package properties

import "math"

type Number float64

type NumberEntity interface {
	Entity
	GetMaxValue() Number
	GetMinValue() Number
}

type NumberProperty struct {
	*Property
}

func NewNumberProperty() *NumberProperty {
	p := &NumberProperty{}
	p.Type = TypeNumber
	return p
}

// SetValue sets a value
func (prop *NumberProperty) SetValue(value float64) {
	//prop.UpdateValue(value)
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

func (prop *NumberProperty) GetMaxValue() Number {
	if v := prop.GetMaximum(); v != nil {
		f, ok := v.(Number)
		if ok {
			return f
		}
	}
	return math.MaxFloat64
}

func (prop *NumberProperty) GetMinValue() Number {
	if v := prop.GetMinimum(); v != nil {
		n, ok := v.(Number)
		if ok {
			return n
		}
	}
	return math.MinInt64
}
