package addon

type Action struct {
	AtType      string        `protobuf:"bytes,1,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
	Title       string        `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description string        `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Links       []*Link       `protobuf:"bytes,4,rep,name=links,proto3" json:"links,omitempty"`
	Input       []interface{} `protobuf:"bytes,5,opt,name=input,proto3,oneof" json:"input,omitempty"`
}

type ActionDescription struct {
	Id            string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string        `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Input         []interface{} `protobuf:"bytes,3,opt,name=input,proto3,oneof" json:"input,omitempty"`
	Status        string        `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	TimeRequested string        `protobuf:"bytes,5,opt,name=timeRequested,proto3" json:"timeRequested,omitempty"`
	TimeCompleted string        `protobuf:"bytes,6,opt,name=timeCompleted,proto3" json:"timeCompleted,omitempty"`
}

func (a *Action) GetDescription() string {
	return ""
}

func (a *Action) GetName() string {
	return ""
}
