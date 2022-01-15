package properties

type StringProperty struct {
	*Property
}

func NewStringProperty(description PropertyDescription) *StringProperty {
	p := &StringProperty{}
	description.Type = TypeString
	p.Property = NewProperty(description)
	return p
}

// SetValue sets a value
func (prop *StringProperty) SetValue(value string) {
	//	prop.UpdateValue(value)
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *StringProperty) OnValueRemoteGet(fn func() string) {
	//prop.OnValueGet(func() interface{} {
	//	return fn()
	//})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *StringProperty) OnValueRemoteUpdate(fn func(string)) {
	//prop.OnValueUpdate(func(Property *addon.PropertyProxy, newValue, oldValue interface{}) {
	//	fn(newValue.(string))
	//)
}
