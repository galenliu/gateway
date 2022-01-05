package devices

import (
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type DeviceProperties map[string]properties.PropertyProxy
type DeviceActions map[string]actions.Action
type DeviceEvents map[string]events.Event

type Device struct {
	adapter             adapter.AdapterProxy
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

func (d Device) AddProperty(n string, p properties.PropertyProxy) {
	d.Properties[n] = p
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

func (d Device) GetAdapter() adapter.AdapterProxy {
	return d.adapter
}

func (d Device) GetLink() []DeviceLink {
	return d.Links
}

func (d Device) GetBaseHref() string {
	return d.BaseHref
}

func (d Device) GetProperties() map[string]properties.PropertyProxy {
	return d.Properties
}

func (d Device) GetProperty(id string) properties.PropertyProxy {
	p, ok := d.Properties[id]
	if ok {
		return p
	}
	return nil
}

func (d Device) GetActions() map[string]actions.Action {
	return d.Actions
}

func (d Device) GetAction(id string) actions.Action {
	return d.Actions[id]
}

func (d Device) GetPin() *DevicePin {
	return d.Pin
}

func (d Device) GetEvents() map[string]events.Event {
	return d.Events
}

func (d Device) GetEvent(id string) events.Event {
	return d.Events[id]
}

func (d Device) GetCredentialsRequired() bool {
	return d.CredentialsRequired
}
