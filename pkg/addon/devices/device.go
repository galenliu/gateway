package devices

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/schemas"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type DeviceProperties map[string]properties.Entity
type DeviceActions map[string]actions.Action
type DeviceEvents map[string]events.Event

type AdapterHandler interface {
	GetId() string
	SendPropertyChangedNotification(deviceId string, p properties.PropertyDescription)
}

type Entity interface {
	SetHandler(h AdapterHandler)
	GetAdapter() AdapterHandler
	GetId() string
	GetAtContext() string
	GetProperty(id string) properties.Entity
	GetAtType() []string
	ToMessage() *messages.Device
	NotifyPropertyChanged(p properties.PropertyDescription)
	SetCredentials(username, password string) error
	SetPin(pin string) error
	GetDevice() *Device
}

type DeviceDescription struct {
	Id          string
	AtType      []Capability
	Title       string
	Description string
}

type Device struct {
	handler             AdapterHandler
	Context             string           `json:"@context,omitempty"`
	AtType              []Capability     `json:"@type,omitempty"`
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

func NewDevice(description DeviceDescription, opts ...Option) *Device {
	device := &Device{
		handler:             nil,
		Context:             schemas.Context,
		AtType:              description.AtType,
		Id:                  description.Id,
		Title:               description.Title,
		Description:         description.Description,
		Links:               nil,
		Forms:               nil,
		BaseHref:            "",
		Pin:                 nil,
		Properties:          make(map[string]properties.Entity, 0),
		Actions:             make(map[string]actions.Action, 0),
		Events:              make(map[string]events.Event, 0),
		CredentialsRequired: false,
	}
	for _, o := range opts {
		o(device)
	}
	if device.Title == "" {
		device.Title = device.AtType[0] + device.Id
	}
	return device
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

func (d *Device) AddProperty(p properties.Entity) {
	if d.Properties == nil {
		d.Properties = make(DeviceProperties, 1)
	}
	p.SetHandler(d)
	d.Properties[p.GetName()] = p
}

func (d *Device) GetId() string {
	return d.Id
}

func (d *Device) GetDevice() *Device {
	return d
}

func (d *Device) GetAtContext() string {
	return d.Context
}

func (d *Device) GetAtType() []string {
	return d.AtType
}

func (d *Device) GetTitle() string {
	return d.Title
}

func (d *Device) GetDescription() string {
	return d.Description
}

func (d *Device) GetLink() []DeviceLink {
	return d.Links
}

func (d *Device) GetBaseHref() string {
	return d.BaseHref
}

func (d *Device) GetProperties() map[string]properties.Entity {
	return d.Properties
}

func (d *Device) GetProperty(id string) properties.Entity {
	p, ok := d.Properties[id]
	if ok {
		return p
	}
	return nil
}

func (d *Device) GetActions() map[string]actions.Action {
	return d.Actions
}

func (d *Device) GetAction(id string) actions.Action {
	return d.Actions[id]
}

func (d *Device) GetPin() *DevicePin {
	return d.Pin
}

func (d *Device) GetEvents() map[string]events.Event {
	return d.Events
}

func (d *Device) GetEvent(id string) events.Event {
	return d.Events[id]
}

func (d *Device) GetCredentialsRequired() bool {
	return d.CredentialsRequired
}

func (d *Device) SetHandler(h AdapterHandler) {
	d.handler = h
}

func (d *Device) GetAdapter() AdapterHandler {
	return d.handler
}

func (d *Device) NotifyPropertyChanged(p properties.PropertyDescription) {
	d.handler.SendPropertyChangedNotification(d.GetId(), p)
}

func (d *Device) ToMessage() *messages.Device {
	return &messages.Device{
		Context: &d.Context,
		Type:    d.AtType,
		BaseHref: func() *string {
			if d.BaseHref == "" {
				return nil
			}
			s := d.BaseHref
			return &s
		}(),
		CredentialsRequired: nil,
		Description:         nil,
		Id:                  d.GetId(),
		Links:               nil,
		Pin:                 nil,
		Properties: func(props DeviceProperties) map[string]messages.Property {
			var mmp = make(map[string]messages.Property)
			for n, p := range props {
				mmp[n] = p.ToMessage()
			}
			return mmp
		}(d.Properties),
		Events: func(es DeviceEvents) map[string]messages.Event {
			var mme = make(map[string]messages.Event)
			for n, e := range es {
				mme[n] = e.ToMessage()
			}
			return mme
		}(d.Events),
		Actions: func(as DeviceActions) map[string]messages.Action {
			var mma = make(map[string]messages.Action)
			for n, e := range as {
				mma[n] = e.ToMessage()
			}
			return mma
		}(d.Actions),
		Title: &d.Title,
	}
}

func (d *Device) SetCredentials(username, password string) error {
	return fmt.Errorf("device:%s SetCredentials not implemented", d.GetId())
}

func (d *Device) SetPin(pin string) error {
	return fmt.Errorf("device:%s SetPin not implemented", d.GetId())
}

type Option func(device *Device)

func WithTitle(title string) Option {
	return func(d *Device) {
		d.Title = title
	}
}
func WithDescription(des string) Option {
	return func(d *Device) {
		d.Description = des
	}
}

func WithCredentialsRequired() Option {
	return func(d *Device) {
		d.CredentialsRequired = true
	}
}

func WithPin(pin string) Option {
	return func(d *Device) {
		d.Pin = &DevicePin{
			Required: true,
			Pattern:  pin,
		}
	}
}
