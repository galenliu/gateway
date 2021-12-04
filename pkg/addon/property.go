package addon

type Property struct {
	Name        string             `json:"name,omitempty"`
	AtType      string             `json:"@type,omitempty"`
	Title       string             `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Description string             `json:"description,omitempty"`
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty"`
	ReadOnly    bool               `json:"readOnly,omitempty"`
	MultipleOf  *float64           `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       interface{}        `json:"value,omitempty"`
}

type PropertyDescription struct {
	Name        string             `json:"name,omitempty"`
	AtType      string             `json:"atType,omitempty"`
	Title       string             `json:"title,omitempty"`
	Type        string             `json:"type,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Description string             `json:"description,omitempty"`
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty"`
	ReadOnly    bool               `json:"readOnly,omitempty"`
	MultipleOf  *float64           `json:"multipleOf,omitempty"`
	Links       []PropertyLinkElem `json:"links,omitempty"`
	Value       interface{}        `json:"value,omitempty"`
}

type PropertyLinkElem struct {
}

func (p Property) GetUnit() string {
	return p.Unit
}

func (p Property) GetAtType() string {
	return p.AtType
}

func (p Property) GetName() string {
	return p.Name
}

func (p Property) GetTitle() string {
	return p.Title
}

func (p Property) GetDescription() string {
	return p.Description
}

func (p Property) GetLinks() []PropertyLinkElem {
	return p.Links
}

func (p Property) GetType() string {
	return p.Type
}

func (p Property) GetMinimum() *float64 {
	return p.Minimum
}

func (p Property) GetMaximum() *float64 {
	return p.Maximum
}

func (p Property) GetEnum() []interface{} {
	return p.Enum
}

func (p Property) GetValue() interface{} {
	return p.Value
}

func (p Property) GetReadOnly() bool {
	return p.ReadOnly
}

func (p Property) GetMultipleOf() *float64 {
	if p.MultipleOf != nil {
		return p.MultipleOf
	}
	return nil
}

func (p Property) SetTitle(s string) bool {
	if p.Title == s {
		return false
	}
	p.Title = s
	return true
}

func (p Property) SetValue(value interface{}) bool {
	if p.Value == value {
		return false
	}
	p.Value = value
	return true
}

func (p Property) SetDescription(s string) bool {
	if p.Description == s {
		return false
	}
	p.Description = s
	return true
}

func (p Property) GetDescriptions() *PropertyDescription {
	return &PropertyDescription{
		Name:        p.Name,
		AtType:      p.AtType,
		Title:       p.Title,
		Type:        p.Type,
		Unit:        p.Unit,
		Description: p.Description,
		Minimum:     p.Minimum,
		Maximum:     p.Maximum,
		Enum:        p.Enum,
		ReadOnly:    p.ReadOnly,
		MultipleOf:  p.MultipleOf,
		Links:       p.Links,
		Value:       p.Value,
	}
}
