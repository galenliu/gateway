package properties

import (
	"fmt"
	"github.com/xiam/to"
)

type BooleanPropertyDescription struct {
	Name        string `json:"name,omitempty"`
	AtType      string `json:"@type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ReadOnly    bool   `json:"readOnly,omitempty"`
	Value       bool   `json:"value,omitempty"`
}

type BooleanEntity interface {
	Entity
	CheckValue(v any) bool
	TurnOn() error
	TurnOff() error
	IsOn() bool
}

type Boolean struct {
	*Property
}

func NewBoolean(description BooleanPropertyDescription, opts ...Option) *Boolean {
	p := &Boolean{}
	p.Property = NewProperty(PropertyDescription{
		Name:        description.Name,
		AtType:      description.AtType,
		Title:       description.Title,
		Type:        TypeBoolean,
		Description: description.Description,
		ReadOnly:    description.ReadOnly,
		Value:       description.Value,
	}, opts...)
	return p
}

// OnValueRemoteGet calls fn when the value was read by a client.
func (prop *Boolean) OnValueRemoteGet(fn func() bool) {
	//prop.OnValueGet(func() interface{} {
	//	return fn()
	//})
}

// OnValueRemoteUpdate calls fn when the value was updated by a client.
func (prop *Boolean) OnValueRemoteUpdate(fn func(bool)) {
	//prop.OnValueUpdate(func(Property *addon.PropertyProxy, newValue, oldValue interface{}) {
	//	fn(newValue.(bool))
	//})
}

func (prop *Boolean) IsOn() bool {
	v := prop.Value.(bool)
	return v
}

func (prop *Boolean) CheckValue(v any) bool {
	return to.Bool(v)
}

func (prop *Boolean) Toggle() {
	fmt.Printf("property: %s Toggle func not implemented", prop.GetName())
}

func (prop *Boolean) TurnOn() error {
	return fmt.Errorf("device:%s property:%s turn on not implemented ", prop.GetDevice().GetId(), prop.GetName())
}

func (prop *Boolean) TurnOff() error {
	return fmt.Errorf("device:%s property:%s turn off not implemented ", prop.GetDevice().GetId(), prop.GetName())
}
