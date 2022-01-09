package properties

import (
	_type "github.com/galenliu/gateway/pkg/addon"
)

type BooleanProperty struct {
	*Property
}

func NewBooleanProperty(device _type.DeviceProxy, description _type.PropertyDescription) *BooleanProperty {
	p := &BooleanProperty{}
	p.Property = NewProperty(device, description)
	if p.Property == nil {
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
	//prop.OnValueUpdate(func(Property *addon.PropertyProxy, newValue, oldValue interface{}) {
	//	fn(newValue.(bool))
	//})
}
