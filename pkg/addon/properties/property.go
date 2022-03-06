package properties

import (
	"encoding/json"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/xiam/to"
	"strings"
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
	GetProperty() *Property
	GetType() Type
	GetTitle() string
	GetName() string
	GetUnit() string
	GetEnum() []any
	GetAtType() PropertyType
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
	Name        string       `json:"name"`
	Title       string       `json:"title,omitempty"`
	Type        Type         `json:"type"`
	AtType      PropertyType `json:"@type,omitempty"`
	Unit        Unit         `json:"unit,omitempty"`
	Description string       `json:"description,omitempty"`
	Minimum     any          `json:"minimum"`
	Maximum     any          `json:"maximum,omitempty"`
	Enum        []any        `json:"enum,omitempty"`
	ReadOnly    bool         `json:"readOnly"`
	MultipleOf  any          `json:"multipleOf,omitempty"`
	Value       any          `json:"value"`
}

func NewProperty(des PropertyDescription, opts ...Option) *Property {

	if des.Type != TypeString && des.Type != TypeInteger && des.Type != TypeNumber &&
		des.Type != TypeBoolean {
		return nil
	}
	p := &Property{
		handler:     nil,
		Name:        des.Name,
		Title:       des.Title,
		Type:        des.Type,
		AtType:      des.AtType,
		Unit:        des.Unit,
		Description: des.Description,
		Minimum:     des.Minimum,
		Maximum:     des.Maximum,
		Enum:        des.Enum,
		ReadOnly:    des.ReadOnly,
		MultipleOf:  des.MultipleOf,
		Value:       des.Value,
	}
	for _, o := range opts {
		o(p)
	}
	if p.Name == "" {
		if p.AtType != "" {
			p.Name = strings.ToLower(p.AtType)
		} else {
			p.Name = strings.ToLower(p.Type)
		}
	}
	return p
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

func (p *Property) GetAtType() PropertyType {
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

func (p *Property) GetProperty() *Property {
	return p
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
				if nb, ok := v.(Number); ok {
					f := float64(nb)
					return &f
				}
				if i, ok := v.(Integer); ok {
					f := float64(i)
					return &f
				}
			}
			return nil
		}(),
		Minimum: func() *float64 {
			if v := p.GetMinimum(); v != nil {
				if nb, ok := v.(Number); ok {
					f := float64(nb)
					return &f
				}
				if i, ok := v.(Integer); ok {
					f := float64(i)
					return &f
				}
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

type Option func(p *Property)

func WithTitle(title string) Option {
	return func(p *Property) {
		p.Title = title
	}
}

func WithDescription(description string) Option {
	return func(p *Property) {
		p.Description = description
	}
}

func WithUnit(unit Unit) Option {
	return func(p *Property) {
		p.Unit = unit
	}
}
