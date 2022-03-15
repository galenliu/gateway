package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/bus/topic"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/xiam/to"
	"golang.org/x/sync/singleflight"
	"sync"
)

var gsf singleflight.Group

type device struct {
	*devices.Device
	logger               logging.Logger
	requestActionTask    sync.Map
	removeActionTask     sync.Map
	setPropertyValueTask sync.Map
}

func newDeviceFromMessage(msg messages.Device) *device {

	linksFunc := func(links []messages.Link) (ls []devices.DeviceLink) {
		if len(links) == 0 {
			return nil
		}
		for _, l := range links {
			ls = append(ls, devices.DeviceLink{
				Href:      l.Href,
				Rel:       util.GetFromPointer(l.Rel),
				MediaType: util.GetFromPointer(l.MediaType),
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
			Pattern:  util.GetFromPointer(pin.Pattern),
		}
	}

	asActionMap := func(as messages.DeviceActions) (oActions map[string]actions.Action) {
		oActions = make(map[string]actions.Action)
		for n, a := range as {
			oActions[n] = actions.Action{
				Type:        util.GetFromPointer(a.Type),
				Title:       util.GetFromPointer(a.Title),
				Description: util.GetFromPointer(a.Description),
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

	asEventMap := func(es messages.DeviceEvents) (events1 map[string]events.Event) {
		events1 = make(map[string]events.Event)
		for n, e := range es {
			events1[n] = events.Event{
				AtType:      "",
				Name:        util.GetFromPointer(e.Name),
				Title:       util.GetFromPointer(e.Name),
				Description: util.GetFromPointer(e.Description),
				Links:       nil,
				Type:        util.GetFromPointer(e.Type),
				Unit:        util.GetFromPointer(e.Unit),
				Minimum:     0,
				Maximum:     0,
				MultipleOf:  0,
				Enum:        nil,
			}
		}
		return
	}

	asPropertiesMap := func(device *device, deviceProperties messages.DeviceProperties) devices.DeviceProperties {
		props := make(devices.DeviceProperties, 0)
		for _, p := range msg.Properties {
			prop := properties.NewProperty(properties.PropertyDescription{
				Name:        util.GetFromPointer(p.Name),
				AtType:      util.GetFromPointer(p.AtType),
				Title:       util.GetFromPointer(p.Title),
				Type:        p.Type,
				Unit:        util.GetFromPointer(p.Unit),
				Description: util.GetFromPointer(p.Description),
				Minimum: func() any {
					if p.Minimum == nil {
						return nil
					}
					return *p.Minimum
				}(),
				Maximum: func() any {
					if p.Maximum == nil {
						return nil
					}
					return *p.Maximum
				}(),
				Enum:     p.Enum,
				ReadOnly: util.GetFromPointer(p.ReadOnly),
				MultipleOf: func() any {
					if p.MultipleOf == nil {
						return nil
					}
					return *p.MultipleOf
				}(),
				Links: nil,
				Value: p.Value,
			})
			if prop == nil {
				continue
			}
			device.AddProperty(prop)
		}
		return props
	}

	device := &device{
		Device: &devices.Device{
			Context:             util.GetFromPointer(msg.Context),
			AtType:              msg.Type,
			Id:                  msg.Id,
			Title:               util.GetFromPointer(msg.Title),
			Description:         util.GetFromPointer(msg.Description),
			Links:               linksFunc(msg.Links),
			BaseHref:            util.GetFromPointer[string](msg.BaseHref),
			Pin:                 pinFunc(msg.Pin),
			Actions:             asActionMap(msg.Actions),
			Events:              asEventMap(msg.Events),
			CredentialsRequired: util.GetFromPointer[bool](msg.CredentialsRequired),
		},
	}
	asPropertiesMap(device, msg.Properties)
	return device
}

func (device *device) onPropertyChanged(property properties.PropertyDescription) {
	t, ok := device.setPropertyValueTask.Load(property.Name)
	if ok {
		task := t.(chan any)
		select {
		case task <- property.Value:
		}
	}
	p := device.GetProperty(property.Name)
	if p == nil {
		return
	}
	if p.IsReadOnly() {
		return
	}
	p.SetCachedValueAndNotify(property.Value)
}

func (device *device) notifyDeviceConnected(connected bool) {
	device.getManager().Publish(topic.DeviceConnected, topic.DeviceConnectedMessage{
		DeviceId:  device.GetId(),
		Connected: connected,
	})
}

func (device *device) actionNotify(actionDescription actions.ActionDescription) {
	device.getManager().Publish(topic.DeviceActionStatus, topic.DeviceActionStatusMessage{
		DeviceId: device.GetId(),
		Action:   actionDescription,
	})
}

func (device *device) requestAction(ctx context.Context, id, name string, input map[string]any) error {
	t, ok := device.requestActionTask.LoadOrStore(id, make(chan bool))
	if !ok {
		var message = messages.DeviceRequestActionRequestJsonData{
			ActionId:   id,
			ActionName: name,
			AdapterId:  device.getAdapter().GetId(),
			DeviceId:   device.GetId(),
			Input:      input,
			PluginId:   device.getAdapter().plugin.pluginId,
		}
		device.getAdapter().Send(messages.MessageType_DeviceRequestActionRequest, message)
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
	defer func() {
		device.setPropertyValueTask.Delete(name)
	}()
	task := t.(chan any)
	if !ok {
		data := messages.DeviceSetPropertyCommandJsonData{
			AdapterId:     device.getAdapter().GetId(),
			DeviceId:      device.GetId(),
			PluginId:      device.getAdapter().plugin.pluginId,
			PropertyName:  name,
			PropertyValue: value,
		}
		device.getAdapter().Send(messages.MessageType_DeviceSetPropertyCommand, data)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout")
		case v := <-task:
			return v, nil
		}
	}

}

func (device *device) getAdapter() *Adapter {
	a := device.Device.GetHandler()
	adapter, ok := a.(*Adapter)
	if ok {
		return adapter
	}
	return nil
}

func (device *device) getManager() *Manager {
	return device.getAdapter().plugin.manager
}
