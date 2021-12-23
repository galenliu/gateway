package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/constant"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	pa "github.com/galenliu/gateway/pkg/wot/definitions/core/property_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

func AsWebOfThing(device *addon.Device) Thing {
	thing := Thing{
		Thing: &wot.Thing{
			AtContext:    controls.URI(device.GetAtContext()),
			Title:        device.GetTitle(),
			Titles:       map[string]string{},
			Id:           controls.URI(device.GetId()),
			AtType:       device.GetAtType(),
			Description:  device.GetDescription(),
			Descriptions: map[string]string{},
			Support:      constant.Support,
			Base:         "",
			Version: &wot.VersionInfo{
				Instance: "1.1",
				Model:    "",
			},
			Created:           nil,
			Modified:          nil,
			Properties:        mapOfWotProperties(device.GetId(), device.Properties),
			Actions:           mapOfWotActions(device.GetId(), device.Actions),
			Events:            mapOfWotEvent(device.GetId(), device.Events),
			Links:             nil,
			Forms:             nil,
			Security:          nil,
			Profile:           nil,
			SchemaDefinitions: nil,
		},
		Pin:                 nil,
		CredentialsRequired: device.CredentialsRequired,
		SelectedCapability:  "",
		Connected:           true,
		GroupId:             "",
	}
	if device.Pin != nil {
		thing.Pin = &ThingPin{
			Required: device.Pin.Required,
			Pattern:  device.Pin.Pattern,
		}
	}
	return thing
}

func mapOfWotProperties(deviceId string, props addon.DeviceProperties) (mapOfProperty map[string]wot.PropertyAffordance) {
	asWotProperty := func(deviceId string, p addon.Property) wot.PropertyAffordance {
		var wp wot.PropertyAffordance
		var i = &ia.InteractionAffordance{
			AtType:       p.AtType,
			Title:        p.Title,
			Titles:       map[string]string{},
			Description:  p.Description,
			Descriptions: map[string]string{},
			Forms: []controls.Form{{
				Href:                controls.URI(fmt.Sprintf("/things/%s/properties/%s", deviceId, p.Name)),
				ContentType:         "",
				ContentCoding:       "",
				Security:            nil,
				Scopes:              nil,
				Response:            nil,
				AdditionalResponses: nil,
				Subprotocol:         "",
				Op:                  nil,
			}},
			UriVariables: nil,
		}
		var dataSchema = &schema.DataSchema{
			AtType:       p.AtType,
			Title:        p.Title,
			Titles:       nil,
			Description:  p.Description,
			Descriptions: nil,
			Const:        nil,
			Default:      nil,
			Unit:         p.Unit,
			OneOf:        nil,
			Enum:         p.GetEnum(),
			ReadOnly:     p.GetReadOnly(),
			WriteOnly:    false,
			Format:       "",
			Type:         p.Type,
		}

		switch p.Type {
		case controls.TypeInteger:
			wp = &pa.IntegerPropertyAffordance{
				InteractionAffordance: i,
				IntegerSchema: &schema.IntegerSchema{
					DataSchema: dataSchema,
					Minimum: func() *controls.Integer {
						var min controls.Integer
						if m := p.GetMinimum(); m != nil {
							min = controls.Integer(*m)
							return &min
						}
						return nil
					}(),
					ExclusiveMinimum: nil,
					Maximum: func() *controls.Integer {
						var max controls.Integer
						if m := p.GetMaximum(); m != nil {
							max = controls.Integer(*m)
							return &max
						}
						return nil
					}(),
					ExclusiveMaximum: nil,
					MultipleOf: func() *controls.Integer {
						var mo controls.Integer
						if m := p.GetMultipleOf(); m != nil {
							mo = controls.Integer(*m)
							return &mo
						}
						return nil
					}(),
				},
				Observable: false,
			}
		case controls.TypeNumber:
			wp = &pa.NumberPropertyAffordance{
				InteractionAffordance: i,
				NumberSchema: &schema.NumberSchema{
					DataSchema: dataSchema,
					Minimum: func() *controls.Double {
						if m := p.GetMinimum(); m != nil {
							return (*controls.Double)(m)
						}
						return nil
					}(),
					ExclusiveMinimum: nil,
					Maximum: func() *controls.Double {
						if m := p.GetMaximum(); m != nil {
							return (*controls.Double)(m)
						}
						return nil
					}(),
					ExclusiveMaximum: nil,
					MultipleOf: func() *controls.Double {
						if m := p.GetMultipleOf(); m != nil {
							return (*controls.Double)(m)
						}
						return nil
					}(),
				},
				Observable: false,
			}
		case controls.TypeBoolean:
			wp = &pa.BooleanPropertyAffordance{
				InteractionAffordance: i,
				BooleanSchema:         &schema.BooleanSchema{DataSchema: dataSchema},
				Observable:            false,
			}
		case controls.TypeString:
			wp = &pa.StringPropertyAffordance{
				InteractionAffordance: i,
				StringSchema: &schema.StringSchema{
					DataSchema: dataSchema,
					MinLength: func() *controls.UnsignedInt {
						if p.Minimum == nil {
							return nil
						}
						var min = controls.UnsignedInt(*p.Minimum)
						return &min
					}(),
					MaxLength: func() *controls.UnsignedInt {
						return nil
					}(),
					Pattern:          "",
					ContentEncoding:  "",
					ContentMediaType: "",
				},
				Observable: false,
			}
		case controls.TypeObject:
			return nil

		default:
			return nil
		}
		return wp
	}
	mapOfProperty = make(map[string]wot.PropertyAffordance)
	for name, p := range props {
		if propertyAffordance := asWotProperty(deviceId, p); propertyAffordance != nil {
			mapOfProperty[name] = propertyAffordance
		}
	}
	return
}

func mapOfWotActions(deviceId string, actions addon.DeviceActions) (mapOfProperty wot.ThingActions) {

	asWotAction := func(deviceId, actionName string, a addon.Action) wot.ActionAffordance {
		var aa = wot.ActionAffordance{}
		var i = &ia.InteractionAffordance{
			AtType:       a.AtType,
			Title:        a.Title,
			Titles:       map[string]string{},
			Description:  a.Description,
			Descriptions: map[string]string{},
			Forms: []controls.Form{{
				Href:                controls.URI(fmt.Sprintf("things/%s/actions/%s", deviceId, actionName)),
				ContentType:         "",
				ContentCoding:       "",
				Security:            nil,
				Scopes:              nil,
				Response:            nil,
				AdditionalResponses: nil,
				Subprotocol:         "",
				Op:                  []string{controls.Op_observeallProperties},
			}},
			UriVariables: nil,
		}
		aa.InteractionAffordance = i
		return aa
	}

	mapOfProperty = make(wot.ThingActions)
	for name, a := range actions {
		if actionAffordance := asWotAction(deviceId, name, a); &actionAffordance != nil {
			mapOfProperty[name] = actionAffordance
		}
	}
	return
}

func mapOfWotEvent(deviceId string, events addon.DeviceEvents) (es wot.ThingEvents) {

	asWotEvent := func(deviceId string, a addon.Event) wot.EventAffordance {
		ea := wot.EventAffordance{
			InteractionAffordance: ia.InteractionAffordance{
				AtType:       a.Type,
				Title:        a.Title,
				Titles:       map[string]string{},
				Description:  a.Description,
				Descriptions: map[string]string{},
				Forms: []controls.Form{
					{
						Href:                controls.URI(fmt.Sprintf("/things/%s/events/%s", deviceId, a.Name)),
						ContentType:         "",
						ContentCoding:       "",
						Security:            nil,
						Scopes:              nil,
						Response:            nil,
						AdditionalResponses: nil,
						Subprotocol:         "",
						Op:                  []string{controls.Op_subscribeEvent},
					},
				},
				UriVariables: nil,
			},
			Subscription: nil,
			Data:         nil,
			Cancellation: nil,
		}
		return ea
	}

	es = make(wot.ThingEvents)
	for n, e := range events {
		es[n] = asWotEvent(deviceId, e)
	}
	return nil
}
