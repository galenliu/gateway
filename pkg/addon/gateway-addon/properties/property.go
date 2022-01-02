package properties

import "encoding/json"

type Property interface {
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
	SetValue(any2 any) bool
	SetTitle(string2 string) bool
	SetDescription(description string) bool
	ToDescription() PropertyDescription
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

type property struct {
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

func NewProperty(description PropertyDescription) *property {
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
	return &property{
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

func (p *property) MarshalJSON() ([]byte, error) {
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

func (p *property) GetName() string {
	return p.Name
}

func (p *property) GetEnum() []any {
	return p.Enum
}

func (p *property) GetTitle() string {
	return p.Title
}

func (p *property) SetTitle(t string) bool {
	if t == p.Title {
		return false
	}
	p.Title = t
	return true
}

func (p *property) SetDescription(d string) bool {
	if d == p.Description {
		return false
	}
	p.Description = d
	return true
}

func (p *property) GetType() string {
	return p.Type
}

func (p *property) GetAtType() string {
	return p.AtType
}

func (p *property) GetUnit() string {
	return p.Unit
}

func (p *property) GetDescription() string {
	return p.Description
}

func (p *property) GetMinimum() any {
	return p.Minimum
}
func (p *property) GetMaximum() any {
	return p.Maximum
}

func (p *property) IsReadOnly() bool {
	return p.ReadOnly
}

func (p *property) GetMultipleOf() any {
	return p.MultipleOf
}

func (p *property) GetValue() any {
	return p.Value
}

func (p *property) SetValue(value any) bool {
	if p.Value == value {
		return false
	}
	p.Value = value
	return true
}

func (p *property) ToDescription() PropertyDescription {
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
