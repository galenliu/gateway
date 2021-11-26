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
	"github.com/xiam/to"
	"time"
)

func AsWebOfThing(device *addon.Device) Thing {
	thing := Thing{
		Thing: &wot.Thing{
			Context:             controls.URI(device.GetAtContext()),
			Title:               device.GetTitle(),
			Id:                  controls.URI(device.GetId()),
			Type:                device.GetType(),
			Description:         device.GetDescription(),
			Support:             "刘桂林",
			Base:                "",
			Version:             nil,
			Created:             &controls.DataTime{Time: time.Now()},
			Modified:            &controls.DataTime{Time: time.Now()},
			Properties:          mapOfWotProperties(device.GetId(), device.Properties),
			Actions:             mapOfWotActions(device.GetId(), device.Actions),
			Events:              nil,
			Links:               nil,
			Forms:               nil,
			Security:            "",
			SecurityDefinitions: nil,
			Profile:             nil,
			SchemaDefinitions:   nil,
		},
		Pin:                 nil,
		CredentialsRequired: device.CredentialsRequired,
		SelectedCapability:  "",
		Connected:           true,
		GroupId:             "",
	}
	if device.Pin != nil {
		thing.Pin = device.GetPin()
	}
	return thing
}

func mapOfWotProperties(deviceId string, props addon.DeviceProperties) (mapOfProperty map[string]wot.PropertyAffordance) {
	mapOfProperty = make(map[string]wot.PropertyAffordance)
	for name, p := range props {
		if propertyAffordance := asWotProperty(deviceId, p); propertyAffordance != nil {
			mapOfProperty[name] = propertyAffordance
		}
	}
	return
}

func mapOfWotActions(deviceId string, actions addon.DeviceActions) (mapOfProperty wot.ThingActions) {
	mapOfProperty = make(wot.ThingActions)
	for name, a := range actions {
		if actionAffordance := asWotAction(deviceId, name, a); &actionAffordance != nil {
			mapOfProperty[name] = actionAffordance
		}
	}
	return
}

func asWotAction(deviceId, actionName string, a addon.Action) wot.ActionAffordance {
	var aa = wot.ActionAffordance{}
	var i = &ia.InteractionAffordance{
		AtType:       "",
		Title:        "",
		Titles:       map[string]string{constant.ZhCN: a.GetTitle()},
		Description:  a.Description,
		Descriptions: map[string]string{constant.ZhCN: a.GetDescription()},
		Forms: []controls.Form{{
			Href:                controls.URI(fmt.Sprintf("/actions/%s/%s", deviceId, actionName)),
			ContentType:         "",
			ContentCoding:       "",
			Security:            "",
			Scopes:              "",
			Response:            nil,
			AdditionalResponses: nil,
			Subprotocol:         "",
			Op:                  "",
		}},
		UriVariables: nil,
	}
	aa.InteractionAffordance = i
	return aa
}

func asWotProperty(deviceId string, p addon.Property) wot.PropertyAffordance {
	var wp wot.PropertyAffordance
	var i = &ia.InteractionAffordance{
		AtType:       p.GetAtType(),
		Title:        p.GetTitle(),
		Titles:       map[string]string{constant.ZhCN: p.GetTitle()},
		Description:  p.Description,
		Descriptions: map[string]string{constant.ZhCN: p.GetDescription()},
		Forms: []controls.Form{{
			Href:                controls.URI(fmt.Sprintf("/things/%s/%s", deviceId, p.Name)),
			ContentType:         "",
			ContentCoding:       "",
			Security:            "",
			Scopes:              "",
			Response:            nil,
			AdditionalResponses: nil,
			Subprotocol:         "",
			Op:                  "",
		}},
		UriVariables: nil,
	}
	var dataSchema = &schema.DataSchema{
		AtType:       p.AtType,
		Title:        p.Title,
		Titles:       nil,
		Description:  p.GetDescription(),
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
		wp = pa.IntegerPropertyAffordance{
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
		wp = pa.NumberPropertyAffordance{
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
			Value:      controls.Double(to.Float64(p.GetValue())),
		}
	case controls.TypeBoolean:
		wp = pa.BooleanPropertyAffordance{
			InteractionAffordance: i,
			BooleanSchema:         &schema.BooleanSchema{DataSchema: dataSchema},
			Observable:            false,
			Value:                 controls.ToBool(p.Value),
		}
	case controls.TypeString:
		wp = pa.StringPropertyAffordance{
			InteractionAffordance: i,
			StringSchema: &schema.StringSchema{
				DataSchema:       dataSchema,
				MinLength:        0,
				MaxLength:        0,
				Pattern:          "",
				ContentEncoding:  "",
				ContentMediaType: "",
			},
			Observable: false,
			Value:      controls.ToString(p.Value),
		}
	case controls.TypeObject:
		return nil

	default:
		return nil
	}
	return wp
}

func asWotEvent(deivceId string, a addon.Action) wot.EventAffordance {
	ea := wot.EventAffordance{}
	return ea
}
