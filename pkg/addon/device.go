package addon

import "github.com/galenliu/gateway/pkg/addon/gateway-addon/properties"

type DeviceProperties map[string]properties.Property
type DeviceActions map[string]Action
type DeviceEvents map[string]Event

type Device struct {
	Context             string           `json:"@context,omitempty"`
	AtType              []string         `json:"@type,omitempty"`
	Id                  string           `json:"id,omitempty"`
	Title               string           `json:"title,omitempty"`
	Description         string           `json:"description,omitempty"`
	Links               []DeviceLink     `json:"links,omitempty"`
	Forms               []DeviceForm     `json:"forms,omitempty"`
	BaseHref            string           `json:"baseHref,omitempty"`
	Pin                 *DevicePin       `json:"pin,omitempty"`
	Properties          DeviceProperties `json:"properties,omitempty"`
	Actions             DeviceActions    `json:"actions,omitempty"`
	Events              DeviceEvents     `json:"events,omitempty"`
	CredentialsRequired bool             `json:"credentialsRequired,omitempty"`
}

type DeviceLink struct {
	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

type DeviceForm struct {
}

type DevicePin struct {
	Required bool   `json:"required,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

func (d Device) GetId() string {
	return d.Id
}

func (d Device) GetAtContext() string {
	return d.Context
}

func (d Device) GetAtType() []string {
	return d.AtType
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

func (d Device) GetProperties() map[string]properties.Property {
	return d.Properties
}

func (d Device) GetProperty(id string) (p properties.Property, ok bool) {
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
