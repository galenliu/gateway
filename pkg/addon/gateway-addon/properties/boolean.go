package properties

type BooleanProperty struct {
	*property
}

func NewBooleanProperty(description PropertyDescription) *BooleanProperty {
	p := &BooleanProperty{}
	p.property = NewProperty(description)
	if p.property == nil {
		return nil
	}
	p.Type = TypeBoolean
	return p
}

// SetBooleanValue SetValue sets a value
func (prop *BooleanProperty) SetBooleanValue(value bool) {
}

// GetValue returns the value as bool
func (prop *BooleanProperty) GetValue() bool {
	return prop.property.GetValue().(bool)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *BooleanProperty) OnValueRemoteGet(fn func() bool) {
	//prop.OnValueGet(func() interface{} {
	//	return fn()
	//})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *BooleanProperty) OnValueRemoteUpdate(fn func(bool)) {
	//prop.OnValueUpdate(func(property *addon.Property, newValue, oldValue interface{}) {
	//	fn(newValue.(bool))
	//})
}
