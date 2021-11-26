package addon

type Property struct {
	Name        string              `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	AtType      string              `protobuf:"bytes,2,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
	Title       string              `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Type        string              `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Unit        string              `protobuf:"bytes,5,opt,name=unit,proto3" json:"unit,omitempty"`
	Description string              `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Minimum     *float64            `protobuf:"fixed32,7,opt,name=minimum,proto3" json:"minimum,omitempty"`
	Maximum     *float64            `protobuf:"fixed32,8,opt,name=maximum,proto3" json:"maximum,omitempty"`
	Enum        []interface{}       `protobuf:"bytes,9,rep,name=enum,proto3" json:"enum,omitempty"`
	ReadOnly    bool                `protobuf:"varint,10,opt,name=readOnly,proto3" json:"readOnly,omitempty"`
	MultipleOf  *float64            `protobuf:"fixed32,11,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
	Links       []*PropertyLinkElem `protobuf:"bytes,103,rep,name=links,proto3" json:"links,omitempty"`
	Value       interface{}         `protobuf:"bytes,12,opt,name=value,proto3" json:"value,omitempty"`
}

type PropertyDescription struct {
	Name        string              `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	AtType      string              `protobuf:"bytes,2,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
	Title       string              `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Type        string              `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Unit        string              `protobuf:"bytes,5,opt,name=unit,proto3" json:"unit,omitempty"`
	Description string              `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Minimum     *float64            `protobuf:"fixed32,7,opt,name=minimum,proto3" json:"minimum,omitempty"`
	Maximum     *float64            `protobuf:"fixed32,8,opt,name=maximum,proto3" json:"maximum,omitempty"`
	Enum        []interface{}       `protobuf:"bytes,9,rep,name=enum,proto3" json:"enum,omitempty"`
	ReadOnly    bool                `protobuf:"varint,10,opt,name=readOnly,proto3" json:"readOnly,omitempty"`
	MultipleOf  *float64            `protobuf:"fixed32,11,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
	Links       []*PropertyLinkElem `protobuf:"bytes,103,rep,name=links,proto3" json:"links,omitempty"`
	Value       interface{}         `protobuf:"bytes,12,opt,name=value,proto3" json:"value,omitempty"`
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

func (p Property) GetLinks() []*PropertyLinkElem {
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
