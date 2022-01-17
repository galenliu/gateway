package properties

import (
	"encoding/json"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
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
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	Enum        []any              `json:"enum,omitempty"`
	ReadOnly    bool               `json:"readOnly,omitempty"`
	MultipleOf  *float64           `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       any                `json:"value,omitempty"`
}

type DeviceHandler interface {
	NotifyPropertyChanged(property PropertyDescription)
}

type Property struct {
	Handler     DeviceHandler
	Name        string   `json:"name"`
	Title       string   `json:"title,omitempty"`
	Type        string   `json:"type"`
	AtType      string   `json:"@type,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Description string   `json:"description,omitempty"`
	Minimum     *float64 `json:"minimum,omitempty"`
	Maximum     *float64 `json:"maximum,omitempty"`
	Enum        []any    `json:"enum,omitempty"`
	ReadOnly    bool     `json:"readOnly"`
	MultipleOf  *float64 `json:"multipleOf,omitempty"`
	Value       any      `json:"value"`
}

func NewProperty(description PropertyDescription) *Property {

	if description.Type != TypeString && description.Type != TypeInteger && description.Type != TypeNumber &&
		description.Type != TypeBoolean {
		return nil
	}
	return &Property{
		Handler:     nil,
		Name:        description.Name,
		Title:       description.Title,
		Type:        description.Type,
		AtType:      description.Description,
		Unit:        description.Unit,
		Description: description.Description,
		Minimum:     description.Minimum,
		Maximum:     description.Maximum,
		Enum:        description.Enum,
		ReadOnly:    description.ReadOnly,
		MultipleOf:  description.MultipleOf,
		Value:       nil,
	}
}

func (p *Property) MarshalJSON() ([]byte, error) {
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

func (p *Property) GetMinimum() *float64 {
	return p.Minimum
}
func (p *Property) GetMaximum() *float64 {
	return p.Maximum
}

func (p *Property) IsReadOnly() bool {
	return p.ReadOnly
}

func (p *Property) GetMultipleOf() *float64 {
	return p.MultipleOf
}

func (p *Property) GetValue() any {
	return p.Value
}

func (p *Property) SetHandler(h DeviceHandler) {
	p.Handler = h
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
		Maximum:     p.Maximum,
		Minimum:     p.Minimum,
		MultipleOf:  p.MultipleOf,
		Name:        get(p.Name),
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
	oldValue := p.GetValue()
	p.SetCachedValue(value)
	hasChanged := oldValue != p.GetValue()
	if hasChanged {
		p.Handler.NotifyPropertyChanged(p.ToDescription())
	}
}

func (p *Property) NotifyChanged() {
	p.Handler.NotifyPropertyChanged(p.ToDescription())
}

func (p *Property) SetCachedValue(value any) bool {
	if p.Value == value {
		return false
	}
	p.Value = value
	return true
}
