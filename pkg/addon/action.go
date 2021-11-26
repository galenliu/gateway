package addon

type ActionInput map[string]interface{}

type Action struct {
	Type        string            `protobuf:"bytes,1,opt,name=atType,json=@type,proto3" json:"atType,omitempty"`
	AtType      string            `json:"@type,omitempty"`
	Title       string            `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Links       []ActionLinksElem `protobuf:"bytes,4,rep,name=links,proto3" json:"links,omitempty"`
	Forms       []ActionFormsElem `protobuf:"bytes,4,rep,name=forms,proto3" json:"forms,omitempty"`
	Input       ActionInput       `protobuf:"bytes,5,opt,name=input,proto3,oneof" json:"input,omitempty"`
}

func (a Action) GetType() string {
	return a.Type
}

func (a Action) GetTitle() string {
	return a.Title
}

func (a Action) GetDescription() string {
	return a.Description
}

func (a Action) GetInput() map[string]interface{} {
	return a.Input
}

type ActionDescription struct {
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Input         map[string]interface{} `protobuf:"bytes,3,opt,name=input,proto3,oneof" json:"input,omitempty"`
	Status        string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	TimeRequested string                 `protobuf:"bytes,5,opt,name=timeRequested,proto3" json:"timeRequested,omitempty"`
	TimeCompleted string                 `protobuf:"bytes,6,opt,name=timeCompleted,proto3" json:"timeCompleted,omitempty"`
}

type ActionLinksElem struct {
}

type ActionFormsElem struct {
	Op interface{} `json:"op"`
}

func (a ActionDescription) GetName() string {
	return a.Name
}

func (a ActionDescription) GetDescription() interface{} {
	return nil
}
