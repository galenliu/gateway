package plugin

import (
	"github.com/galenliu/gateway/pkg/addon"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
)

type Device struct {
	adapter *Adapter
	*addon.Device
}

func newDevice(msg messages.Device) *Device {
	linksFunc := func(links []messages.Link) (links1 []addon.DeviceLink) {
		for _, l := range links {
			links1 = append(links1, addon.DeviceLink{
				Href:      l.Href,
				Rel:       *l.Rel,
				MediaType: *l.MediaType,
			})
		}
		return nil
	}
	pinFunc := func(pin *messages.Pin) (pin1 *addon.DevicePin) {
		return &addon.DevicePin{
			Required: pin.Required,
			Pattern:  *pin.Pattern,
		}
	}

	propertiesFunc := func(properties messages.DeviceProperties) (properties1 map[string]addon.Property) {
		properties1 = make(map[string]addon.Property)
		for n, p := range properties {
			properties1[n] = addon.Property{
				Name:        *p.Name,
				AtType:      *p.AtType,
				Title:       *p.Title,
				Type:        p.Type,
				Unit:        *p.Unit,
				Description: *p.Description,
				Minimum:     *p.Minimum,
				Maximum:     *p.Maximum,
				Enum: func(elems []messages.PropertyEnumElem) []interface{} {
					var enums []interface{}
					for e := range elems {
						enums = append(enums, e)
					}
					return enums
				}(p.Enum),
				ReadOnly:   *p.ReadOnly,
				MultipleOf: *p.MultipleOf,
				Links:      nil,
				Value:      p.Value,
			}
		}
		return
	}

	actionsFunc := func(actions messages.DeviceActions) (actions1 map[string]addon.Action) {
		actions1 = make(map[string]addon.Action)
		for n, a := range actions {
			actions1[n] = addon.Action{
				Type:        *a.Type,
				Title:       *a.Title,
				Description: *a.Description,
				Forms: func(s []messages.ActionFormsElem) (forms []addon.ActionFormsElem) {
					if len(s) == 0 {
						return nil
					}
					for _, a := range s {
						forms = append(forms, addon.ActionFormsElem{
							Op: a,
						})
					}
					return nil
				}(a.Forms),
				Input: nil,
			}
		}
		return
	}

	eventFunc := func(events messages.DeviceEvents) (events1 map[string]addon.Event) {
		events1 = make(map[string]addon.Event)
		for n, e := range events {
			events1[n] = addon.Event{
				AtType:      "",
				Name:        *e.Name,
				Title:       *e.Title,
				Description: *e.Description,
				Links:       nil,
				Type:        *e.Type,
				Unit:        "",
				Minimum:     0,
				Maximum:     0,
				MultipleOf:  0,
				Enum:        nil,
			}
		}
		return
	}

	device := &Device{
		adapter: nil,
		Device: &addon.Device{
			Context:             *msg.Context,
			Type:                msg.Type,
			Id:                  msg.Id,
			Title:               *msg.Title,
			Description:         *msg.Description,
			Links:               linksFunc(msg.Links),
			BaseHref:            *msg.BaseHref,
			Pin:                 pinFunc(msg.Pin),
			Properties:          propertiesFunc(msg.Properties),
			Actions:             actionsFunc(msg.Actions),
			Events:              eventFunc(msg.Events),
			CredentialsRequired: *msg.CredentialsRequired,
		},
	}

	return device
}

func (device *Device) getAdapter() *Adapter {
	return device.adapter
}

func (device *Device) notifyValueChanged(property messages.Property) {
	p, ok := device.GetProperty(*property.Name)
	if !ok {
		return
	}
	if p.ReadOnly {
		return
	}
	valueChanged := p.SetValue(property.Value)
	titleChanged := false
	if property.Title != nil {
		titleChanged = p.SetTitle(*property.Title)
	}
	descriptionChanged := false
	if property.Description != nil {
		descriptionChanged = p.SetDescription(*property.Description)
	}
	if valueChanged || descriptionChanged || titleChanged {
		device.adapter.plugin.pluginServer.manager.bus.PublishPropertyChanged(device.GetId(), p.GetDescriptions())
	}
}

func (device *Device) requestAction(description addon.ActionDescription) {
	var message = messages.DeviceRequestActionRequestJsonData{
		ActionId:   description.Id,
		ActionName: description.Name,
		AdapterId:  device.getAdapter().getId(),
		DeviceId:   device.GetId(),
		Input:      description.Input,
		PluginId:   device.getAdapter().getPlugin().pluginId,
	}
	device.adapter.send(messages.MessageType_DeviceRequestActionRequest, message)
}

func (device *Device) notifyDeviceConnected(connected bool) {
	device.adapter.plugin.pluginServer.manager.bus.PublishConnected(device.GetId(), connected)
}

func (device *Device) notifyAction(actionDescription *addon.ActionDescription) {
	device.adapter.plugin.pluginServer.manager.bus.PublishActionStatus(actionDescription)
}

func (device *Device) setPropertyValue(name string, value interface{}) {
	data := messages.DeviceSetPropertyCommandJsonData{
		AdapterId:     device.getAdapter().getId(),
		DeviceId:      device.GetId(),
		PluginId:      device.getAdapter().getPlugin().getId(),
		PropertyName:  name,
		PropertyValue: value,
	}
	device.getAdapter().send(messages.MessageType_DeviceSetPropertyCommand, data)
}
