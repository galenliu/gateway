package virtual_adapter

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type Property struct {
	*properties.Property
}

func NewVirtualProperty(p *properties.Property) *Property {
	return &Property{p}
}

func (p *Property) SetPropertyValue(v any) error {
	if p.GetAtType() == properties.TypeThermostatModeProperty {
		str, ok := v.(string)
		if ok {
			if str != "auto" {
				p.GetDevice().NotifyPropertyChanged(properties.PropertyDescription{
					Name:        "heatingCooling",
					AtType:      properties.TypeHeatingCoolingProperty,
					Title:       "Heating/Cooling",
					Type:        "String",
					Unit:        properties.UnitDegreeCelsius,
					Description: "",
					Minimum:     nil,
					Maximum:     nil,
					Enum:        nil,
					ReadOnly:    false,
					MultipleOf:  nil,
					Links:       nil,
					Value:       str,
				})
			}
		}
	}
	p.SetCachedValue(v)
	p.NotifyChanged()
	fmt.Printf("device: %s set property: %s value: %v \t\n", p.GetDevice().GetId(), p.GetName(), v)
	return nil
}