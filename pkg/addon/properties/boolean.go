package properties

import (
	"fmt"
	"github.com/xiam/to"
)

type BooleanEntity interface {
	Entity
	CheckValue(v any) bool
}

type BooleanProperty struct {
	*Property
}

func NewBooleanProperty(description PropertyDescription) *BooleanProperty {
	p := &BooleanProperty{}
	description.Type = TypeBoolean
	p.Property = NewProperty(description)
	return p
}

func (prop *BooleanProperty) Turn(b bool) {
	fmt.Printf("property: %s Turn func not implemented", prop.GetName())
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

func (prop *BooleanProperty) GetValue() bool {
	v := prop.Value.(bool)
	return v
}

func (prop *BooleanProperty) CheckValue(v any) bool {
	return to.Bool(v)
}
