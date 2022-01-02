package devices

import (
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Device struct {
	ID          string   `json:"id"`
	Context     string   `json:"@context"`
	AtType      []string `json:"@type"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`

	Links               []*Link                        `json:"links,omitempty"`
	PinRequired         bool                           `json:"pinRequired"`
	CredentialsRequired bool                           `json:"credentialsRequired"`
	BaseHref            string                         `json:"baseHref"`
	Pin                 *Pin                           `json:"pin,omitempty"`
	Properties          map[string]properties.Property `json:"properties,omitempty"`
	Actions             map[string]*addon.Action       `json:"action,omitempty"`
	Events              map[string]*addon.Event        `json:"events,omitempty"`
}

type Pin struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
}

type Link struct {
	Href      string `protobuf:"bytes,1,opt,name=href,proto3" json:"href,omitempty"`
	Rel       string `protobuf:"bytes,2,opt,name=rel,proto3" json:"rel,omitempty"`
	MediaType string `protobuf:"bytes,3,opt,name=mediaType,proto3" json:"mediaType,omitempty"`
}

func NewDeviceFormMessage(dev *addon.Device) *Device {
	device := &Device{
		ID:                  dev.Id,
		Context:             dev.GetAtContext(),
		AtType:              dev.GetAtType(),
		Title:               dev.Title,
		Description:         dev.Description,
		PinRequired:         dev.Pin.Required,
		BaseHref:            dev.BaseHref,
		CredentialsRequired: dev.CredentialsRequired,
	}
	if len(dev.Links) > 0 {
		device.Links = make([]*Link, 2)
		for _, l := range dev.Links {
			device.Links = append(device.Links, &Link{
				Href:      l.Href,
				Rel:       l.Rel,
				MediaType: l.MediaType,
			})
		}
	}
	if dev.Pin != nil {
		device.Pin = &Pin{
			Required: dev.Pin.Required,
			Pattern:  dev.Pin.Pattern,
		}
	}
	if len(dev.Properties) > 0 {
		device.Properties = make(map[string]properties.Property)
		//for name, property := range dev.Properties {
		//	device.Properties[name] = properties.NewPropertyFormMessage(property)
		//}
	}

	if len(dev.Events) > 0 {
		device.Events = make(map[string]*addon.Event)
		//for name, event := range dev.Events {
		//device.Events[name] = events.NewEventFormMessage(event)
		//}
	}
	if len(dev.Actions) > 0 {
		device.Actions = make(map[string]*addon.Action)
		//for name, action := range dev.actions {
		//	device.actions[name] = actions.NewActionFormMessage(action)
		//}
	}
	return device
}

func NewDeviceFormString(des string) *Device {
	var device Device
	//err := json.UnmarshalFromString(des, &device)
	//if err != nil {
	//	return nil
	//}
	return &device
}

func (device *Device) GetID() string {
	return device.ID
}

func (device *Device) GetAtContext() string {
	return device.Context
}

func (device *Device) GetTitle() string {
	return device.Title
}

func (device *Device) GetDescription() string {
	return device.Description
}

func (device *Device) AddProperty(property properties.Property) {
	if device.Properties == nil {
		device.Properties = make(map[string]properties.Property)
	}
	device.Properties[property.GetName()] = property
}

func (device *Device) AddAction(action *addon.Action) {
	if device.Actions == nil {
		device.Actions = make(map[string]*addon.Action)
	}
	//device.actions[action.Name] = action
}

func (device *Device) AddEvent(event *addon.Event) {
	if device.Events == nil {
		device.Events = make(map[string]*addon.Event)
	}
	device.Events[event.Name] = event
}
