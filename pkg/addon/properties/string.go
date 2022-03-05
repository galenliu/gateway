package properties

import (
	"fmt"
	"github.com/xiam/to"
)

type StringPropertyDescription struct {
	Name        string             `json:"name,omitempty"`
	AtType      PropertyType       `json:"@type,omitempty"`
	Title       string             `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Description string             `json:"description,omitempty"`
	Enum        []string           `json:"enum,omitempty"`
	ReadOnly    bool               `json:"readOnly,omitempty"`
	MultipleOf  string             `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       string             `json:"value,omitempty"`
}

type StringEntity interface {
	Entity
	CheckValue(v any) string
	SetValue(v string) error
}

type StringProperty struct {
	*Property
}

func NewStringProperty(desc StringPropertyDescription) *StringProperty {
	s := &StringProperty{}
	s.Property = NewProperty(PropertyDescription{
		Name:        desc.Name,
		AtType:      desc.AtType,
		Title:       desc.Title,
		Type:        TypeString,
		Unit:        desc.Unit,
		Description: desc.Description,
		Enum: func() []any {
			enum := make([]any, 0)
			for _, e := range desc.Enum {
				enum = append(enum, e)
			}
			return enum
		}(),
		ReadOnly:   desc.ReadOnly,
		MultipleOf: desc.MultipleOf,
		Value:      desc.Value,
	})
	return s
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
