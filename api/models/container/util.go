package container

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/events"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/constant"
	wot "github.com/galenliu/gateway/pkg/wot/definitions/core"
	ia "github.com/galenliu/gateway/pkg/wot/definitions/core/interaction_affordance"
	pa "github.com/galenliu/gateway/pkg/wot/definitions/core/property_affordance"
	schema "github.com/galenliu/gateway/pkg/wot/definitions/data_schema"
	controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"
)

func AsWebOfThing(device devices.Device) *Thing {

	thing := Thing{
		Thing: &wot.Thing{
			AtContext:    controls.URI(device.GetAtContext()),
			Title:        device.GetTitle(),
			Titles:       map[string]string{},
			Id:           controls.NewURI(device.GetId()),
			AtType:       device.GetAtType(),
			Description:  device.GetDescription(),
			Descriptions: map[string]string{},
			Support:      constant.Support,
			Version: &wot.VersionInfo{
				Model: "1.0.0",
			},
			Created:           nil,
			Modified:          nil,
			Links:             nil,
			Security:          nil,
			Profile:           nil,
			SchemaDefinitions: nil,
		},
		CredentialsRequired: device.CredentialsRequired,
		SelectedCapability:  "",
		Connected:           true,
		GroupId:             "",
	}

	thing.Properties = mapOfWotProperties(thing.Thing, device.Properties)
	thing.Actions = mapOfWotActions(thing.Thing, device.Actions)
	thing.Events = mapOfWotEvent(thing.Thing, device.Events)
	thing.Forms = arrayOfThingFrom(thing.Thing)

	if device.Pin != nil {
		thing.Pin = &ThingPin{
			Required: device.Pin.Required,
			Pattern:  device.Pin.Pattern,
		}
	}
	return &thing
}

func mapOfWotProperties(thing *wot.Thing, props devices.DeviceProperties) (mapOfProperty map[string]wot.PropertyAffordance) {
	asWotProperty := func(thing *wot.Thing, p properties.Entity) wot.PropertyAffordance {
		var wp wot.PropertyAffordance
		var i = &ia.InteractionAffordance{
			AtType:       p.GetAtType(),
			Title:        p.GetTitle(),
			Titles:       map[string]string{},
			Description:  p.GetDescription(),
			Descriptions: map[string]string{},
			UriVariables: nil,
		}
		var dataSchema = &schema.DataSchema{
			AtType:       p.GetAtType(),
			Title:        p.GetTitle(),
			Titles:       nil,
			Description:  p.GetDescription(),
			Descriptions: nil,
			Const:        nil,
			Default:      nil,
			Unit:         p.GetUnit(),
			OneOf:        nil,
			Enum:         p.GetEnum(),
			ReadOnly:     p.IsReadOnly(),
			WriteOnly:    false,
			Format:       "",
			Type:         p.GetType(),
		}
		form := controls.Form{
			Href:                controls.URI(fmt.Sprintf("%s/properties/%s", thing.Base, p.GetName())),
			ContentType:         controls.JSON,
			ContentCoding:       "",
			Security:            nil,
			Scopes:              nil,
			Response:            nil,
			AdditionalResponses: nil,
			Subprotocol:         "",
		}
		if !dataSchema.ReadOnly && !dataSchema.WriteOnly {
			form.Op = controls.NewOpArray(controls.Readproperty, controls.Writeproperty)
		}
		if dataSchema.ReadOnly {
			form.Op = controls.NewOpArray(controls.Readproperty)
		}
		if dataSchema.WriteOnly {
			form.Op = controls.NewOpArray(controls.Writeproperty)
		}
		i.Forms = make([]controls.Form, 0)
		i.Forms = append(i.Forms, form)

		switch p.GetType() {
		case controls.TypeInteger:
			wp = &pa.IntegerPropertyAffordance{
				InteractionAffordance: i,
				IntegerSchema: &schema.IntegerSchema{
					DataSchema: dataSchema,
					Minimum: func() *controls.Integer {
						if v := p.GetMinimum(); v != nil {
							i := controls.ToInteger(v)
							return &i
						}
						return nil
					}(),
					ExclusiveMinimum: nil,
					Maximum: func() *controls.Integer {
						if v := p.GetMaximum(); v != nil {
							i := controls.ToInteger(v)
							return &i
						}
						return nil
					}(),
					ExclusiveMaximum: nil,
					MultipleOf: func() *controls.Integer {
						if v := p.GetMultipleOf(); v != nil {
							i := controls.ToInteger(v)
							return &i
						}
						return nil
					}(),
				},
				Observable: false,
			}
			fmt.Printf("")
		case controls.TypeNumber:
			wp = &pa.NumberPropertyAffordance{
				InteractionAffordance: i,
				NumberSchema: &schema.NumberSchema{
					DataSchema: dataSchema,
					Minimum: func() *controls.Double {
						if m := p.GetMinimum(); m != nil {
							d := controls.ToDouble(m)
							return &d
						}
						return nil
					}(),
					ExclusiveMinimum: nil,
					Maximum: func() *controls.Double {
						if m := p.GetMaximum(); m != nil {
							d := controls.ToDouble(m)
							return &d
						}
						return nil
					}(),
					ExclusiveMaximum: nil,
					MultipleOf: func() *controls.Double {
						if m := p.GetMultipleOf(); m != nil {
							d := controls.ToDouble(m)
							return &d
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
						if v := p.GetMinimum(); v != nil {
							i := controls.ToUnsignedInt(v)
							return &i
						}
						return nil
					}(),
					MaxLength: func() *controls.UnsignedInt {
						if v := p.GetMaximum(); v != nil {
							i := controls.ToUnsignedInt(v)
							return &i
						}
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
		if propertyAffordance := asWotProperty(thing, p); propertyAffordance != nil {
			mapOfProperty[name] = propertyAffordance
		}
	}
	return
}

func mapOfWotActions(thing *wot.Thing, as devices.DeviceActions) (mapOfProperty wot.ThingActions) {

	asWotAction := func(thing *wot.Thing, actionName string, a actions.Action) wot.ActionAffordance {
		var aa = wot.ActionAffordance{}
		var i = &ia.InteractionAffordance{
			AtType:       a.AtType,
			Title:        a.Title,
			Titles:       map[string]string{},
			Description:  a.Description,
			Descriptions: map[string]string{},
			Forms: []controls.Form{{
				Href:                controls.URI(fmt.Sprintf("%s/actions/%s", thing.Base, actionName)),
				ContentType:         "",
				ContentCoding:       "",
				Security:            nil,
				Scopes:              nil,
				Response:            nil,
				AdditionalResponses: nil,
				Subprotocol:         "",
				Op:                  controls.NewOpArray(controls.Invokeaction),
			}},
			UriVariables: nil,
		}
		aa.InteractionAffordance = i
		return aa
	}

	mapOfProperty = make(wot.ThingActions)
	for name, a := range as {
		if actionAffordance := asWotAction(thing, name, a); &actionAffordance != nil {
			mapOfProperty[name] = &actionAffordance
		}
	}
	return
}

func mapOfWotEvent(thing *wot.Thing, des devices.DeviceEvents) (es wot.ThingEvents) {

	asWotEvent := func(thing *wot.Thing, a events.Event) wot.EventAffordance {
		ea := wot.EventAffordance{
			InteractionAffordance: ia.InteractionAffordance{
				AtType:       a.Type,
				Title:        a.Title,
				Titles:       map[string]string{},
				Description:  a.Description,
				Descriptions: map[string]string{},
				Forms: []controls.Form{
					{
						Href:                controls.NewURI(fmt.Sprintf(thing.Id.ToString()+"/things/%s/events/%s", thing.GetId(), a.Name)),
						ContentType:         controls.JSON,
						ContentCoding:       "",
						Security:            nil,
						Scopes:              nil,
						Response:            nil,
						AdditionalResponses: nil,
						Subprotocol:         "",
						Op:                  controls.NewOpArray(controls.SubscribeEvent),
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
	for n, e := range des {
		es[n] = asWotEvent(thing, e)
	}
	return nil
}

func arrayOfThingFrom(t *wot.Thing) []controls.Form {
	var fs = make([]controls.Form, 0)
	if t.Properties != nil && len(t.Properties) > 0 {
		fs = append(fs, controls.Form{
			Href:                "/properties",
			ContentType:         controls.JSON,
			ContentCoding:       "",
			Security:            nil,
			Scopes:              nil,
			Response:            nil,
			AdditionalResponses: nil,
			Subprotocol:         "",
			Op:                  controls.NewArrayOrString(controls.Readallproperties),
		})
	}
	if t.Events != nil && len(t.Events) > 0 {
		fs = append(fs, controls.Form{
			Href:                "/events",
			ContentType:         controls.JSON,
			ContentCoding:       "",
			Security:            nil,
			Scopes:              nil,
			Response:            nil,
			AdditionalResponses: nil,
			Subprotocol:         "",
			Op:                  controls.NewArrayOrString(controls.Subscribeallevents),
		})
	}
	if len(fs) == 0 {
		return nil
	}
	return fs
}
