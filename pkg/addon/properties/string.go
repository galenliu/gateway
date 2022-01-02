package properties

type StringProperty struct {
	*property
}

func NewStringProperty(description PropertyDescription) *StringProperty {
	p := &StringProperty{}
	p.property = NewProperty(description)
	if p.property == nil {
		return nil
	}
	p.Type = TypeBoolean
	return p
}

func newStringProperty() *StringProperty {
	p := &StringProperty{}
	p.Type = TypeString
	return p
}

// SetValue sets a value
func (prop *StringProperty) SetValue(value string) {
	//	prop.UpdateValue(value)
}

// GetValue returns the value as bool
func (prop *StringProperty) GetValue() string {
	return prop.property.GetValue().(string)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *StringProperty) OnValueRemoteGet(fn func() string) {
	//prop.OnValueGet(func() interface{} {
	//	return fn()
	//})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *StringProperty) OnValueRemoteUpdate(fn func(string)) {
	//prop.OnValueUpdate(func(property *addon.Property, newValue, oldValue interface{}) {
	//	fn(newValue.(string))
	//)
}
