package addon

type Event struct {
	AtType      string        `protobuf:"bytes,1,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
	Name        string        `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Title       string        `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Description string        `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Links       []*Link       `protobuf:"bytes,5,rep,name=links,proto3" json:"links,omitempty"`
	Type        string        `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	Unit        string        `protobuf:"bytes,7,opt,name=unit,proto3" json:"unit,omitempty"`
	Minimum     float32       `protobuf:"fixed32,8,opt,name=minimum,proto3" json:"minimum,omitempty"`
	Maximum     float32       `protobuf:"fixed32,9,opt,name=maximum,proto3" json:"maximum,omitempty"`
	MultipleOf  float32       `protobuf:"fixed32,10,opt,name=multipleOf,proto3" json:"multipleOf,omitempty"`
	Enum        []interface{} `protobuf:"bytes,11,rep,name=enum,proto3" json:"enum,omitempty"`
}

func (e Event) GetAtType() string {
	return e.AtType
}

func (e Event) GetName() string {
	return e.Name
}

func (e Event) GetTitle() string {
	return e.Title
}

func (e Event) GetDescription() string {
	return e.Description
}

func (e Event) GetLinks() []*Link {
	return e.Links
}

func (e Event) GetType() string {
	return e.Type
}

func (e Event) GetEnum() []interface{} {
	return e.Enum
}

func (e Event) GetMinimum() float32 {
	return e.Minimum
}

func (e Event) GetMaximum() float32 {
	return e.Maximum
}
