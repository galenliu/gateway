package properties

import "github.com/xiam/to"

type StringEntity interface {
	Entity
	CheckValue(v any) string
}

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

func (prop *StringProperty) CheckValue(v any) string {
	s := to.String(v)
	return s
}

func (prop *StringProperty) GetValue() string {
	v := prop.Value.(string)
	return v
}
