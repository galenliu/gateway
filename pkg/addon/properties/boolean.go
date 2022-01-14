package properties

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

type BooleanProperty struct {
	*proxy.Boolean
}

func (prop *BooleanProperty) SetValue(b bool) {
	//TODO implement me
	fmt.Print()
}

func NewBooleanProperty(description PropertyDescription) *BooleanProperty {
	p := &BooleanProperty{}
	p.Property = NewProperty(description)
	if p.Property == nil {
		return nil
	}
	p.Type = TypeBoolean
	return p
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
