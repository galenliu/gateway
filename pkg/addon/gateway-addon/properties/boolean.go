package properties

type BooleanProperty struct {
	*Property
}

func NewBooleanProperty(typ string) *BooleanProperty {
	p := &BooleanProperty{}
	p.Type = typ
	p.Type = TypeBoolean
	return p
}

// SetValue sets a value
func (prop *BooleanProperty) SetBooleanValue(value bool) {
}

// GetValue returns the value as bool
func (prop *BooleanProperty) GetValue() bool {
	return prop.Property.GetValue().(bool)
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
