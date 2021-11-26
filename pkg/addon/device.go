package addon

type DeviceProperties map[string]Property
type DeviceActions map[string]Action
type DeviceEvents map[string]Event

type Device struct {
	Context             string           `protobuf:"bytes,1,opt,name=atContext,json=@context,proto3" json:"@context,omitempty"`
	Type                []string         `protobuf:"bytes,2,opt,name=atType,json=@type,proto3" json:"@type,omitempty"`
	Id                  string           `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Title               string           `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Description         string           `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Links               []DeviceLink     `protobuf:"bytes,6,rep,name=links,proto3" json:"links,omitempty"`
	Forms               []DeviceForm     `protobuf:"bytes,6,rep,name=forms,proto3" json:"forms,omitempty"`
	BaseHref            string           `protobuf:"bytes,7,opt,name=baseHref,proto3" json:"baseHref,omitempty"`
	Pin                 *DevicePin       `protobuf:"bytes,8,opt,name=pin,proto3" json:"pin,omitempty"`
	Properties          DeviceProperties `protobuf:"bytes,9,rep,name=properties,proto3" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Actions             DeviceActions    `protobuf:"bytes,10,rep,name=actions,proto3" json:"actions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Events              DeviceEvents     `protobuf:"bytes,11,rep,name=events,proto3" json:"events,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CredentialsRequired bool             `protobuf:"varint,12,opt,name=credentialsRequired,proto3" json:"credentialsRequired,omitempty"`
}

type DeviceLink struct {
	Href      string `protobuf:"bytes,1,opt,name=href,proto3" json:"href,omitempty"`
	Rel       string `protobuf:"bytes,2,opt,name=rel,proto3" json:"rel,omitempty"`
	MediaType string `protobuf:"bytes,3,opt,name=mediaType,proto3" json:"mediaType,omitempty"`
}

type DeviceForm struct {
}

type DevicePin struct {
	Required bool   `protobuf:"varint,1,opt,name=required,proto3" json:"required"`
	Pattern  string `protobuf:"bytes,2,opt,name=pattern,proto3" json:"pattern"`
}

func (d Device) GetId() string {
	return d.Id
}

func (d Device) GetAtContext() string {
	return d.Context
}

func (d Device) GetType() []string {
	return d.Type
}

func (d Device) GetTitle() string {
	return d.Title
}

func (d Device) GetDescription() string {
	return d.Description
}

func (d Device) GetLink() []DeviceLink {
	return d.Links
}

func (d Device) GetBaseHref() string {
	return d.BaseHref
}

func (d Device) GetProperties() map[string]Property {
	return d.Properties
}

func (d Device) GetProperty(id string) (p Property, ok bool) {
	p, ok = d.Properties[id]
	return
}

func (d Device) GetActions() map[string]Action {
	return d.Actions
}

func (d Device) GetAction(id string) Action {
	return d.Actions[id]
}

func (d Device) GetPin() *DevicePin {
	return d.Pin
}

func (d Device) GetEvents() map[string]Event {
	return d.Events
}

func (d Device) GetEvent(id string) Event {
	return d.Events[id]
}

func (d Device) GetCredentialsRequired() bool {
	return d.CredentialsRequired
}
