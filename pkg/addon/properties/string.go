package properties

import (
	"fmt"
	"github.com/xiam/to"
)

type StringEntity interface {
	Entity
	CheckValue(v any) string
	SetValue(v string) error
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
func (prop *StringProperty) SetValue(v string) error {
	return fmt.Errorf("device:%s property:%s set value:%v not implemented ", prop.GetDevice().GetId(), prop.GetName(), v)
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
