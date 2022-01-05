package properties

import (
	"encoding/json"
	"github.com/galenliu/gateway/pkg/addon/adapter"
)

type PropertyProxy interface {
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
	GetValue() any
	IsReadOnly() bool
	SetValue(a any) bool
	SetTitle(s string) bool
	SetDescription(description string) bool
	ToDescription() PropertyDescription
}

type DeviceProxy interface {
	GetId() string
	GetAdapter() adapter.AdapterProxy
	GetProperty(id string) PropertyProxy
	NotifyPropertyChanged(p PropertyDescription)
}

type PropertyDescription struct {
	Name        *string            `json:"name,omitempty"`
	AtType      *string            `json:"@type,omitempty"`
	Title       *string            `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Unit        *string            `json:"unit,omitempty"`
	Description *string            `json:"description,omitempty"`
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	Enum        []any              `json:"enum,omitempty"`
	ReadOnly    *bool              `json:"readOnly,omitempty"`
	MultipleOf  *float64           `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       any                `json:"value,omitempty"`
}

type PropertyLinkElem struct {
}

type Property struct {
	device      DeviceProxy
	Name        string   `json:"name"`
	Title       string   `json:"title,omitempty"`
	Type        string   `json:"type"`
	AtType      string   `json:"@type,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Description string   `json:"description,omitempty"`
	Minimum     any      `json:"minimum,omitempty"`
	Maximum     any      `json:"maximum,omitempty"`
	Enum        []any    `json:"enum,omitempty"`
	ReadOnly    bool     `json:"readOnly"`
	MultipleOf  *float64 `json:"multipleOf,omitempty"`
	//Links       []*rpc.Link   `json:"links"`
	Value any `json:"value"`
}

func NewProperty(device DeviceProxy, description PropertyDescription) *Property {
	getString := func(s *string) string {
		if s != nil {
			return *s
		}
		return ""
	}
	if description.Type != TypeString && description.Type != TypeInteger && description.Type != TypeNumber &&
		description.Type != TypeBoolean {
		return nil
	}
	return &Property{
		device:      device,
		Name:        getString(description.Name),
		Title:       getString(description.Title),
		Type:        description.Type,
		AtType:      getString(description.Description),
		Unit:        getString(description.Unit),
		Description: getString(description.Description),
		Minimum:     description.Minimum,
		Maximum:     description.Maximum,
		Enum:        description.Enum,
		ReadOnly: func(b *bool) bool {
			if b == nil {
				return false
			}
			return *b
		}(description.ReadOnly),
		MultipleOf: description.MultipleOf,
		Value:      nil,
	}
}

func (p *Property) MarshalJSON() ([]byte, error) {
	propertyDescription := PropertyDescription{
		Name:        &p.Name,
		AtType:      &p.AtType,
		Title:       &p.Title,
		Type:        p.Title,
		Unit:        &p.Unit,
		Description: &p.Description,
		Minimum: func(a any) *float64 {
			switch a.(type) {
			case float64:
				return a.(*float64)
			default:
				return nil
			}
		}(p.Minimum),
		Maximum: func(a any) *float64 {
			switch a.(type) {
			case float64:
				return a.(*float64)
			default:
				return nil
			}
		}(p.Maximum),
		Enum:       nil,
		ReadOnly:   nil,
		MultipleOf: nil,
		Links:      nil,
		Value:      nil,
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

func (p *Property) GetValue() any {
	return p.Value
}

func (p *Property) SetValue(value any) bool {
	if p.Value == value {
		return false
	}
	p.Value = value
	return true
}

func (p *Property) ToDescription() PropertyDescription {
	get := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}
	getFloat := func(s any) *float64 {
		if s == nil {
			return nil
		}
		switch s.(type) {
		case float64:
			return s.(*float64)
		}
		return nil
	}
	return PropertyDescription{
		Name:        get(p.Name),
		AtType:      get(p.AtType),
		Title:       get(p.Title),
		Type:        p.Type,
		Unit:        get(p.Unit),
		Description: get(p.Description),
		Minimum:     getFloat(p.Minimum),
		Maximum:     getFloat(p.Maximum),
		Enum:        p.Enum,
		ReadOnly: func(b bool) *bool {
			return &b
		}(p.IsReadOnly()),
		MultipleOf: p.MultipleOf,
		Links:      nil,
		Value:      p.Value,
	}
}

func (p *Property) SetCachedValueAndNotify(value any) {
	oldValue := p.GetValue()
	p.SetCachedValue(value)
	hasChanged := oldValue != p.GetValue()
	if hasChanged {
		p.device.NotifyPropertyChanged(p.ToDescription())
	}
}

func (p *Property) SetCachedValue(value any) {
	p.Value = value
}
