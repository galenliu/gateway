package properties

import (
	"encoding/json"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/xiam/to"
)

type PropertyLinkElem struct {
}

type PropertyDescription struct {
	Name        string             `json:"name,omitempty"`
	AtType      string             `json:"@type,omitempty"`
	Title       string             `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Description string             `json:"description,omitempty"`
	Minimum     any                `json:"minimum,omitempty"`
	Maximum     any                `json:"maximum,omitempty"`
	Enum        []any              `json:"enum,omitempty"`
	ReadOnly    bool               `json:"readOnly,omitempty"`
	MultipleOf  any                `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       any                `json:"value,omitempty"`
}

type DeviceHandler interface {
	GetId() string
	NotifyPropertyChanged(property PropertyDescription)
}

type Entity interface {
	MarshalJSON() ([]byte, error)
	GetType() string
	GetTitle() string
	GetName() string
	GetUnit() string
	GetEnum() []any
	GetAtType() string
	GetDescription() string
	GetMinimum() any
	GetMaximum() any
	GetMultipleOf() any
	IsReadOnly() bool
	ToMessage() messages.Property
	ToDescription() PropertyDescription
	SetCachedValue(value any) bool
	SetCachedValueAndNotify(value any)
	SetTitle(s string) bool
	SetDescription(description string) bool
	SetHandler(d DeviceHandler)
	GetDevice() DeviceHandler
	NotifyChanged()
	GetCachedValue() any
}

type Property struct {
	handler     DeviceHandler
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type"`
	AtType      string `json:"@type,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
	Minimum     any    `json:"minimum"`
	Maximum     any    `json:"maximum,omitempty"`
	Enum        []any  `json:"enum,omitempty"`
	ReadOnly    bool   `json:"readOnly"`
	MultipleOf  any    `json:"multipleOf,omitempty"`
	Value       any    `json:"value"`
}

func NewProperty(description PropertyDescription) *Property {

	if description.Type != TypeString && description.Type != TypeInteger && description.Type != TypeNumber &&
		description.Type != TypeBoolean {
		return nil
	}
	return &Property{
		handler:     nil,
		Name:        description.Name,
		Title:       description.Title,
		Type:        description.Type,
		AtType:      description.AtType,
		Unit:        description.Unit,
		Description: description.Description,
		Minimum:     description.Minimum,
		Maximum:     description.Maximum,
		Enum:        description.Enum,
		ReadOnly:    description.ReadOnly,
		MultipleOf:  description.MultipleOf,
		Value:       description.Value,
	}
}

func (p Property) MarshalJSON() ([]byte, error) {
	propertyDescription := PropertyDescription{
		Name:        p.Name,
		AtType:      p.AtType,
		Title:       p.Title,
		Type:        p.Title,
		Unit:        p.Unit,
		Description: p.Description,
		Minimum:     p.Minimum,
		Maximum:     p.Maximum,
		Enum:        nil,
		ReadOnly:    p.ReadOnly,
		MultipleOf:  p.MultipleOf,
		Links:       nil,
		Value:       p.Value,
	}
	return json.Marshal(propertyDescription)
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetEnum() []any {
	return p.Enum
}

func (p *Property) GetTitle() string {
	return p.Title
}

func (p *Property) SetTitle(t string) bool {
	if t == p.Title {
		return false
	}
	p.Title = t
	return true
}

func (p *Property) SetDescription(d string) bool {
	if d == p.Description {
		return false
	}
	p.Description = d
	return true
}

func (p *Property) GetType() string {
	return p.Type
}

func (p *Property) GetAtType() string {
	return p.AtType
}

func (p *Property) GetUnit() string {
	return p.Unit
}

func (p *Property) GetDescription() string {
	return p.Description
}

func (p *Property) GetMinimum() any {
	return p.Minimum
}
func (p *Property) GetMaximum() any {
	return p.Maximum
}

func (p *Property) IsReadOnly() bool {
	return p.ReadOnly
}

func (p *Property) GetMultipleOf() any {
	return p.MultipleOf
}

func (p *Property) GetCachedValue() any {
	return p.Value
}

func (p *Property) SetHandler(h DeviceHandler) {
	p.handler = h
}

func (p *Property) GetDevice() DeviceHandler {
	return p.handler
}

func (p *Property) ToDescription() PropertyDescription {

	return PropertyDescription{
		Name:        p.Name,
		AtType:      p.AtType,
		Title:       p.Title,
		Type:        p.Type,
		Unit:        p.Unit,
		Description: p.Description,
		Minimum:     p.Minimum,
		Maximum:     p.Maximum,
		Enum:        p.Enum,
		ReadOnly:    p.IsReadOnly(),
		MultipleOf:  p.MultipleOf,
		Links:       nil,
		Value:       p.Value,
	}
}

func (p *Property) ToMessage() messages.Property {
	get := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}
	return messages.Property{
		Type:        p.Type,
		AtType:      get(p.AtType),
		Description: get(p.Description),
		Enum:        p.Enum,
		Maximum: func() *float64 {
			if v := p.GetMaximum(); v != nil {
				f := to.Float64(v)
				return &f
			}
			return nil
		}(),
		Minimum: func() *float64 {
			if v := p.GetMinimum(); v != nil {
				f := to.Float64(v)
				return &f
			}
			return nil
		}(),
		MultipleOf: func() *float64 {
			if v := p.GetMultipleOf(); v != nil {
				f := to.Float64(v)
				return &f
			}
			return nil
		}(),
		Name: get(p.Name),
		ReadOnly: func(b bool) *bool {
			if b {
				return &b
			}
			return nil
		}(p.ReadOnly),
		Title: get(p.Title),
		Unit:  get(p.Unit),
		Value: p.Value,
	}
}

func (p *Property) SetCachedValueAndNotify(value any) {
	oldValue := p.GetCachedValue()
	p.SetCachedValue(value)
	hasChanged := oldValue != p.GetCachedValue()
	if hasChanged {
		if p.handler != nil {
			p.handler.NotifyPropertyChanged(p.ToDescription())
		}
	}
}

func (p *Property) NotifyChanged() {
	if p.handler != nil {
		p.handler.NotifyPropertyChanged(p.ToDescription())
	}
}

func (p *Property) SetCachedValue(value any) bool {
	if p.Value == value {
		return false
	}
	p.Value = value
	return true
}
