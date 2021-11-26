package properties

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type Property struct {
	Name        string        `json:"name"`
	Title       string        `json:"title,omitempty"`
	Type        string        `json:"type"`
	AtType      string        `json:"@type,omitempty"`
	Unit        string        `json:"unit,omitempty"`
	Description string        `json:"description,omitempty"`
	Minimum     interface{}   `json:"minimum,omitempty"`
	Maximum     interface{}   `json:"maximum,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	ReadOnly    bool          `json:"readOnly"`
	MultipleOf  interface{}   `json:"multipleOf,omitempty"`
	//Links       []*rpc.Link   `json:"links"`
	Value interface{} `json:"value"`
}

func NewPropertyFormMessage(p *messages.Property) *Property {
	property := &Property{
		Name:        *p.Name,
		Title:       *p.Title,
		Type:        p.Type,
		AtType:      *p.AtType,
		Unit:        *p.Unit,
		Description: *p.Description,
		Minimum:     p.Minimum,
		Maximum:     p.Maximum,
		Enum:        nil,
		ReadOnly:    *p.ReadOnly,
		MultipleOf:  p.MultipleOf,
		//Links:       p.Links,
		Value: p.Value,
	}
	return property
}

//
//func NerPropertyFormString(s string) *Property {
//	data := []byte(s)
//	p := Property{}
//	p.Name = json.Get(data, "name").ToString()
//	p.Type = json.Get(data, "type").ToString()
//	p.AtType = json.Get(data, "@type").ToString()
//	p.Unit = json.Get(data, "unit").ToString()
//	p.Description = json.Get(data, "description").ToString()
//	p.Minimum = json.Get(data, "minimum").GetInterface()
//	p.Maximum = json.Get(data, "maximum").GetInterface()
//	p.ReadOnly = json.Get(data, "readOnly").ToBool()
//	var e []interface{}
//	json.Get(data, "enum").ToVal(&e)
//	p.MultipleOf = json.Get(data, "multipleOf").GetInterface()
//	var f []*rpc.Link
//	json.Get(data, "forms").ToVal(&f)
//	p.Links = f
//	if p.Name == "" || p.Type == "" {
//		return nil
//	}
//	_ = json.UnmarshalFromString(s, &p)
//	return &p
//
//}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetTitle() string {
	return p.Title
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

func (p *Property) GetMinimum() interface{} {
	return p.Minimum
}
func (p *Property) GetMaximum() interface{} {
	return p.Maximum
}

func (p *Property) IsReadOnly() bool {
	return p.ReadOnly
}

func (p *Property) GetMultipleOf() interface{} {
	return p.MultipleOf
}

//func (p *Property) GetForms() []*messages.Link {
//	return p.Links
//}

func (p *Property) GetValue() interface{} {
	return p.Value
}
