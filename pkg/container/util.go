package container

import (
	"github.com/galenliu/gateway/pkg/addon"
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
			AtContext:           controls.URI(device.GetAtContext()),
			Title:               device.GetTitle(),
			Id:                  controls.URI(device.GetId()),
			AtType:              device.GetAtType(),
			Description:         device.GetDescription(),
			Support:             "刘桂林",
			Base:                "",
			Version:             nil,
			Created:             &controls.DataTime{Time: time.Now()},
			Modified:            &controls.DataTime{Time: time.Now()},
			Properties:          MapOfWotProperty(device.Properties),
			Actions:             nil,
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

func MapOfWotProperty(props map[string]*addon.Property) (mapOfProperty map[string]wot.PropertyAffordance) {
	mapOfProperty = make(map[string]wot.PropertyAffordance)
	for name, p := range props {
		if propertyAffordance := AsWotProperty(p); propertyAffordance != nil {
			mapOfProperty[name] = propertyAffordance
		}
	}
	return
}

func AsWotProperty(p *addon.Property) *wot.PropertyAffordance {
	var wp wot.PropertyAffordance
	var i = &ia.InteractionAffordance{
		AtType:       p.AtType,
		Title:        p.Title,
		Titles:       nil,
		Description:  p.Description,
		Descriptions: nil,
		Forms:        nil,
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
				DataSchema:       dataSchema,
				Minimum:          controls.ToInteger(p.Minimum),
				ExclusiveMinimum: 0,
				Maximum:          controls.ToInteger(p.Maximum),
				ExclusiveMaximum: 0,
				MultipleOf:       controls.ToInteger(p.MultipleOf),
			},
			Observable: false,
		}
	case controls.TypeNumber:
		wp = pa.NumberPropertyAffordance{
			InteractionAffordance: i,
			NumberSchema: &schema.NumberSchema{
				DataSchema:       dataSchema,
				Minimum:          controls.Double(to.Float64(p.Minimum)),
				ExclusiveMinimum: 0,
				Maximum:          controls.Double(to.Float64(p.Maximum)),
				ExclusiveMaximum: 0,
				MultipleOf:       controls.ToDouble(p.MultipleOf),
			},
			Observable: false,
			Value:      controls.Double(to.Float64(p.Value)),
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
	return &wp
}
