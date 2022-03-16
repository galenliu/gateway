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
	if len(device.AtType) == 0 {
		return nil
	}
	thing := Thing{
		Thing: &wot.Thing{
			AtContext:         controls.URI(device.GetAtContext()),
			Title:             device.GetTitle(),
			Titles:            map[string]string{},
			Id:                controls.NewURI(device.GetId()),
			AtType:            device.GetAtType(),
			Description:       device.GetDescription(),
			Descriptions:      map[string]string{},
			Support:           constant.Support,
			Created:           nil,
			Modified:          nil,
			Links:             nil,
			Security:          nil,
			Profile:           nil,
			SchemaDefinitions: nil,
		},
		CredentialsRequired: device.CredentialsRequired,
		SelectedCapability:  device.GetAtType()[0],
		GroupId:             "",
	}

	thing.Properties = mapOfWotProperties(device.Properties)
	thing.Actions = mapOfWotActions(device, device.Actions)
	thing.Events = mapOfWotEvent(device, device.Events)
	thing.Forms = arrayOfThingForms(thing.Thing)

	if device.Pin != nil {
		thing.Pin = &ThingPin{
			Required: device.Pin.Required,
			Pattern:  device.Pin.Pattern,
		}
	}
	return &thing
}

func mapOfWotProperties(props devices.DeviceProperties) (mapOfProperty map[string]wot.PropertyAffordance) {

	asWotProperty := func(p properties.Entity) wot.PropertyAffordance {
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
			Titles:       map[string]string{},
			Description:  p.GetDescription(),
			Descriptions: map[string]string{},
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
			Href:        controls.NewURI(fmt.Sprintf("things/%s/properties/%s", p.GetDevice().GetId(), p.GetName())),
			ContentType: controls.JSON,
		}
		if !dataSchema.ReadOnly && !dataSchema.WriteOnly {
			form.Op = controls.NewOpArray(controls.Readproperty, controls.Writeproperty)
		} else if dataSchema.ReadOnly {
			form.Op = controls.NewOpArray(controls.Readproperty)
		} else if dataSchema.WriteOnly {
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
				Value:      p.GetCachedValue(),
			}
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
				Value:      p.GetCachedValue(),
			}
		case controls.TypeBoolean:
			wp = &pa.BooleanPropertyAffordance{
				InteractionAffordance: i,
				BooleanSchema:         &schema.BooleanSchema{DataSchema: dataSchema},
				Observable:            false,
				Value:                 p.GetCachedValue(),
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
				Value:      p.GetCachedValue(),
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
		if propertyAffordance := asWotProperty(p); propertyAffordance != nil {
			mapOfProperty[name] = propertyAffordance
		}
	}
	return
}

func mapOfWotActions(device devices.Device, as devices.DeviceActions) (mapOfProperty wot.ThingActions) {
	asWotAction := func(actionName string, a actions.Action) wot.ActionAffordance {
		var aa = wot.ActionAffordance{}
		var i = &ia.InteractionAffordance{
			AtType:       a.AtType,
			Title:        a.Title,
			Titles:       map[string]string{},
			Description:  a.Description,
			Descriptions: map[string]string{},
			Forms: []controls.Form{{
				Href:        controls.NewURI(fmt.Sprintf("things/%s/actions/%s", device.GetId(), actionName)),
				ContentType: controls.JSON,
				Op:          controls.NewOpArray(controls.Invokeaction),
			}},
			UriVariables: nil,
		}
		aa.InteractionAffordance = i
		return aa
	}
	mapOfProperty = make(wot.ThingActions)
	for name, a := range as {
		if actionAffordance := asWotAction(name, a); &actionAffordance != nil {
			mapOfProperty[name] = &actionAffordance
		}
	}
	return
}

func mapOfWotEvent(device devices.Device, des devices.DeviceEvents) (es wot.ThingEvents) {

	asWotEvent := func(a events.Event) *wot.EventAffordance {
		ea := &wot.EventAffordance{
			InteractionAffordance: ia.InteractionAffordance{
				AtType:       a.Type,
				Title:        a.Title,
				Titles:       map[string]string{},
				Description:  a.Description,
				Descriptions: map[string]string{},
				Forms: []controls.Form{
					{
						Href:        controls.NewURI(fmt.Sprintf("things/%s/events/%s", device.GetId(), a.Name)),
						ContentType: controls.JSON,
						Op:          controls.NewOpArray(controls.SubscribeEvent),
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
		e := asWotEvent(e)
		if e != nil {
			es[n] = e
		}
	}
	if len(es) == 0 {
		return nil
	}
	return es
}

func arrayOfThingForms(t *wot.Thing) []controls.Form {
	var fs = make([]controls.Form, 0)
	if t.Properties != nil && len(t.Properties) > 0 {
		fs = append(fs, controls.Form{
			Href:        controls.NewURI(fmt.Sprintf("/things/%s/properties", t.GetId())),
			ContentType: controls.JSON,
			Op:          controls.NewArrayOrString(controls.Readallproperties),
		})
	}
	if t.Events != nil && len(t.Events) > 0 {
		fs = append(fs, controls.Form{
			Href:        controls.NewURI(fmt.Sprintf("/things/%s/events", t.GetId())),
			ContentType: controls.JSON,
			Op:          controls.NewArrayOrString(controls.Subscribeallevents),
		})
	}

	if t.Actions != nil && len(t.Actions) > 0 {
		fs = append(fs, controls.Form{
			Href:        controls.NewURI(fmt.Sprintf("/things/%s/actions", t.GetId())),
			ContentType: controls.JSON,
			Op:          controls.NewArrayOrString(controls.Queryallactions),
		})
	}
	if len(fs) == 0 {
		return nil
	}
	return fs
}

type Valuer interface {
	~string | ~float64 | ~bool
}

func GetPointerFormValue[T Valuer](value T) *T {
	switch value.(type) {
	case string:
		str := value
		return &str

	}
	return nil
}
