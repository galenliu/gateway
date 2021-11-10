package addon

type Property struct {
	Name        string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	AtType      string        `protobuf:"bytes,2,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
	Title       string        `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Type        string        `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Unit        string        `protobuf:"bytes,5,opt,name=unit,proto3" json:"unit,omitempty"`
	Description string        `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Minimum     float32       `protobuf:"fixed32,7,opt,name=minimum,proto3" json:"minimum,omitempty"`
	Maximum     float32       `protobuf:"fixed32,8,opt,name=maximum,proto3" json:"maximum,omitempty"`
	Enum        []interface{} `protobuf:"bytes,9,rep,name=enum,proto3" json:"enum,omitempty"`
	ReadOnly    bool          `protobuf:"varint,10,opt,name=readOnly,proto3" json:"readOnly,omitempty"`
	MultipleOf  float32       `protobuf:"fixed32,11,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
	Links       []*Link       `protobuf:"bytes,103,rep,name=links,proto3" json:"links,omitempty"`
	Value       interface{}   `protobuf:"bytes,12,opt,name=value,proto3" json:"value,omitempty"`
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

func (p Property) GetLinks() []*Link {
	return p.Links
}

func (p Property) GetType() string {
	return p.Type
}

func (p Property) GetMinimum() float32 {
	return p.Minimum
}

func (e Property) GetMaximum() float32 {
	return e.Maximum
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

func (p Property) GetMultipleOf() float32 {
	return p.MultipleOf
}
