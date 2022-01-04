package properties

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

// GetValue returns the value as bool
func (prop *NumberProperty) GetValue() float64 {
	return prop.Property.GetValue().(float64)
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
