package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/adapter"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/bus/topic"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/xiam/to"
	"sync"
)

type device struct {
	adapter *Adapter
	*devices.Device
	logger               logging.Logger
	requestActionTask    sync.Map
	setPropertyValueTask sync.Map
}

func newDevice(adapter *Adapter, msg messages.Device) *device {

	getString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	linksFunc := func(links []messages.Link) (ls []devices.DeviceLink) {
		if len(links) == 0 {
			return nil
		}
		for _, l := range links {
			ls = append(ls, devices.DeviceLink{
				Href:      l.Href,
				Rel:       getString(l.Rel),
				MediaType: getString(l.MediaType),
			})
		}
		return ls
	}

	pinFunc := func(pin *messages.Pin) *devices.DevicePin {
		if pin == nil {
			return nil
		}
		return &devices.DevicePin{
			Required: pin.Required,
			Pattern:  getString(pin.Pattern),
		}
	}

	getPropertyDescription := func(p messages.Property) properties.PropertyDescription {
		return properties.PropertyDescription{
			Name:        p.Name,
			AtType:      p.AtType,
			Title:       p.Title,
			Type:        p.Type,
			Unit:        p.Unit,
			Description: p.Description,
			Minimum:     p.Minimum,
			Maximum:     p.Maximum,
			Enum:        p.Enum,
			ReadOnly:    p.ReadOnly,
			MultipleOf:  p.MultipleOf,
			Links:       nil,
			Value:       p.Value,
		}
	}

	actionsFunc := func(as messages.DeviceActions) (oActions map[string]actions.Action) {
		oActions = make(map[string]actions.Action)
		for n, a := range as {
			oActions[n] = actions.Action{
				Type:        getString(a.Type),
				Title:       getString(a.Title),
				Description: getString(a.Description),
				Forms: func(s []messages.ActionFormsElem) (forms []actions.ActionFormsElem) {
					if len(s) == 0 {
						return nil
					}
					for _, a := range s {
						forms = append(forms, actions.ActionFormsElem{
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

	eventsFunc := func(es messages.DeviceEvents) (events1 map[string]events.Event) {
		events1 = make(map[string]events.Event)
		for n, e := range es {
			events1[n] = events.Event{
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

	device := device{
		adapter: adapter,
		Device: &devices.Device{
			Context: getString(msg.Context),
			AtType:  msg.Type,
			Id:      msg.Id,
			Title: func() string {
				t := getString(msg.Title)
				if t != "" {
					return t
				}
				return msg.Id
			}(),
			Description: "",
			Links:       linksFunc(msg.Links),
			BaseHref:    *msg.BaseHref,
			Pin:         pinFunc(msg.Pin),
			Actions:     actionsFunc(msg.Actions),
			Events:      eventsFunc(msg.Events),
			CredentialsRequired: func() bool {
				if msg.CredentialsRequired == nil {
					return false
				}
				return *msg.CredentialsRequired
			}(),
		},
	}
	for n, p := range msg.Properties {
		prop := properties.NewProperty(&device, getPropertyDescription(p))
		if prop == nil {
			continue
		}
		device.AddProperty(n, prop)
	}

	device.logger = adapter.logger
	return &device
}

func (device *device) NotifyPropertyChanged(property properties.PropertyDescription) {
	t, ok := device.setPropertyValueTask.Load(*property.Name)
	if ok {
		task := t.(chan any)
		select {
		case task <- property.Value:
		}
	}
	p := device.GetProperty(*property.Name)
	if p == nil {
		return
	}
	if p.IsReadOnly() {
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
		device.adapter.plugin.manager.bus.Pub(topic.DevicePropertyChanged, device.GetId(), p.ToDescription())
	}
}

func (device *device) notifyDeviceConnected(connected bool) {
	device.adapter.plugin.manager.bus.Pub(topic.DeviceConnected, device.GetId(), connected)
}

func (device *device) notifyAction(actionDescription *actions.ActionDescription) {
	device.adapter.plugin.manager.bus.Pub(topic.DeviceActionStatus, device.GetId(), actionDescription)
}

func (device *device) requestAction(ctx context.Context, id, name string, input map[string]any) error {

	t, ok := device.requestActionTask.LoadOrStore(id, make(chan bool))
	if !ok {
		var message = messages.DeviceRequestActionRequestJsonData{
			ActionId:   id,
			ActionName: name,
			AdapterId:  device.GetAdapter().GetId(),
			DeviceId:   device.GetId(),
			Input:      input,
			PluginId:   device.adapter.plugin.pluginId,
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
			return fmt.Errorf("actions task failed")
		}
	}
}

func (device *device) setPropertyValue(ctx context.Context, name string, value any) (any, error) {

	p := device.GetProperty(name)
	if p == nil {
		return nil, errors.New("property not found")
	}
	if p.GetType() == TypeBoolean {
		value = to.Bool(value)
	}
	if p.GetType() == TypeNumber {
		value = to.Float64(value)
	}
	if p.GetType() == TypeInteger {
		value = to.Int64(value)
	}
	if p.GetType() == TypeString {
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
			AdapterId:     device.GetAdapter().GetId(),
			DeviceId:      device.GetId(),
			PluginId:      device.adapter.plugin.pluginId,
			PropertyName:  name,
			PropertyValue: value,
		}
		device.adapter.send(messages.MessageType_DeviceSetPropertyCommand, data)
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

func (device *device) GetAdapter() adapter.AdapterProxy {
	return device.adapter
}

func (device *device) getAdapter() *Adapter {
	return device.adapter
}
