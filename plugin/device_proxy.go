package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/bus/topic"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/xiam/to"
	"sync"
)

type Device struct {
	adapter *Adapter
	*addon.Device
	logger               logging.Logger
	requestActionTask    sync.Map
	setPropertyValueTask sync.Map
}

func newDevice(adapter *Adapter, msg messages.Device) *Device {

	getString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	linksFunc := func(links []messages.Link) (ls []addon.DeviceLink) {
		if len(links) == 0 {
			return nil
		}
		for _, l := range links {
			ls = append(ls, addon.DeviceLink{
				Href:      l.Href,
				Rel:       getString(l.Rel),
				MediaType: getString(l.MediaType),
			})
		}
		return ls
	}

	pinFunc := func(pin *messages.Pin) *addon.DevicePin {
		if pin == nil {
			return nil
		}
		return &addon.DevicePin{
			Required: pin.Required,
			Pattern:  getString(pin.Pattern),
		}
	}

	propertiesFunc := func(properties messages.DeviceProperties) (properties1 map[string]addon.Property) {
		properties1 = make(map[string]addon.Property)
		for n, p := range properties {
			properties1[n] = addon.Property{
				Name:        getString(p.Name),
				AtType:      getString(p.AtType),
				Title:       getString(p.Title),
				Type:        p.Type,
				Unit:        getString(p.Unit),
				Description: getString(p.Description),
				Minimum:     p.Minimum,
				Maximum:     p.Maximum,
				Enum: func(elems []messages.PropertyEnumElem) []any {
					var enums []any
					for e := range elems {
						enums = append(enums, e)
					}
					return enums
				}(p.Enum),
				ReadOnly: func() bool {
					if p.ReadOnly == nil {
						return false
					}
					return *p.ReadOnly
				}(),
				MultipleOf: p.MultipleOf,
				Links:      nil,
				Value:      p.Value,
			}
		}
		return
	}

	actionsFunc := func(actions messages.DeviceActions) (oActions map[string]addon.Action) {
		oActions = make(map[string]addon.Action)
		for n, a := range actions {
			oActions[n] = addon.Action{
				Type:        getString(a.Type),
				Title:       getString(a.Title),
				Description: getString(a.Description),
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

	eventsFunc := func(events messages.DeviceEvents) (events1 map[string]addon.Event) {
		events1 = make(map[string]addon.Event)
		for n, e := range events {
			events1[n] = addon.Event{
				AtType:      "",
				Name:        getString(e.Name),
				Title:       getString(e.Name),
				Description: getString(e.Description),
				Links:       nil,
				Type:        getString(e.Type),
				Unit:        getString(e.Unit),
				Minimum:     0,
				Maximum:     0,
				MultipleOf:  0,
				Enum:        nil,
			}
		}
		return
	}

	device := &Device{
		adapter: adapter,
		Device: &addon.Device{
			Context: getString(msg.Context),
			Type:    msg.Type,
			Id:      msg.Id,
			Title: func() string {
				t := getString(msg.Title)
				if t != "" {
					return t
				}
				return msg.Id
			}(),
			Description: func() string {
				d := getString(msg.Description)
				if d != "" {
					return d
				}
				return *msg.Description
			}(),
			Links:      linksFunc(msg.Links),
			BaseHref:   *msg.BaseHref,
			Pin:        pinFunc(msg.Pin),
			Properties: propertiesFunc(msg.Properties),
			Actions:    actionsFunc(msg.Actions),
			Events:     eventsFunc(msg.Events),
			CredentialsRequired: func() bool {
				if msg.CredentialsRequired == nil {
					return false
				}
				return *msg.CredentialsRequired
			}(),
		},
	}
	device.logger = adapter.logger
	return device
}

func (device *Device) notifyValueChanged(property messages.Property) {
	t, ok := device.setPropertyValueTask.Load(*property.Name)
	if ok {
		task := t.(chan any)
		select {
		case task <- property.Value:
		}
	}
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
		device.adapter.plugin.manager.bus.Pub(topic.DevicePropertyChanged, device.GetId(), p.GetDescriptions())
	}
}

func (device *Device) notifyDeviceConnected(connected bool) {
	device.adapter.plugin.manager.bus.Pub(topic.DeviceConnected, device.GetId(), connected)
}

func (device *Device) notifyAction(actionDescription *addon.ActionDescription) {
	device.adapter.plugin.manager.bus.Pub(topic.DeviceActionStatus, device.GetId(), actionDescription)
}

func (device *Device) requestAction(ctx context.Context, id, name string, input map[string]any) error {

	t, ok := device.requestActionTask.LoadOrStore(id, make(chan bool))
	if !ok {
		var message = messages.DeviceRequestActionRequestJsonData{
			ActionId:   id,
			ActionName: name,
			AdapterId:  device.getAdapter().getId(),
			DeviceId:   device.GetId(),
			Input:      input,
			PluginId:   device.getAdapter().getPlugin().pluginId,
		}
		device.adapter.send(messages.MessageType_DeviceRequestActionRequest, message)
	}
	task := t.(chan bool)
	defer func() {
		device.requestActionTask.Delete(id)
	}()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("requestActionTask timeout")
		case b := <-task:
			if b {
				return nil
			}
			return fmt.Errorf("action task failed")
		}
	}
}

func (device *Device) setPropertyValue(ctx context.Context, name string, value any) (any, error) {

	p, _ := device.GetProperty(name)
	if p.Type == TypeBoolean {
		value = to.Bool(value)
	}
	if p.Type == TypeNumber {
		value = to.Float64(value)
	}
	if p.Type == TypeInteger {
		value = to.Int64(value)
	}
	if p.Type == TypeString {
		value = to.String(value)
	}

	t, ok := device.setPropertyValueTask.LoadOrStore(name, make(chan any))
	task := t.(chan any)

	defer func() {
		device.setPropertyValueTask.Delete(name)
	}()

	defer func() {
		device.setPropertyValueTask.Delete(name)
	}()
	if !ok {
		data := messages.DeviceSetPropertyCommandJsonData{
			AdapterId:     device.getAdapter().getId(),
			DeviceId:      device.GetId(),
			PluginId:      device.getAdapter().getPlugin().getId(),
			PropertyName:  name,
			PropertyValue: value,
		}
		device.getAdapter().send(messages.MessageType_DeviceSetPropertyCommand, data)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout for setPropertyValue")
		case v := <-task:
			return v, nil
		}
	}
}

func (device *Device) getAdapter() *Adapter {
	return device.adapter
}
