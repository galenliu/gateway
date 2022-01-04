package properties

type IntegerProperty struct {
	*Property
}

func NewIntegerProperty() *IntegerProperty {
	return &IntegerProperty{&Property{
		Type: TypeInteger,
	}}
}

// SetValue sets a value
func (prop *IntegerProperty) SetCachedValueAndNotify(value int) {
	//prop.PropertyProxy.SetCachedValueAndNotify(value)
}

// GetValue returns the value as bool
func (prop *IntegerProperty) GetValue() int {
	return prop.Property.GetValue().(int)
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

func (prop *IntegerProperty) SetMinValue(v int64) {
	//prop.PropertyProxy.SetMinValue(v)
}

func (prop *IntegerProperty) SetMaxValue(v int64) {
	//prop.PropertyProxy.SetMaxValue(v)
}
