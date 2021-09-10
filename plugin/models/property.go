package models

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/xiam/to"
	"math"
)

type device interface {
	NotifyValueChanged(property *Property)
}

type Property struct {
	device      device
	Name        string        `json:"name"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	AtType      string        `json:"@type"`
	Unit        string        `json:"unit"`
	Description string        `json:"description"`
	Minimum     interface{}   `json:"minimum"`
	Maximum     interface{}   `json:"maximum"`
	Enum        []interface{} `json:"enum"`
	ReadOnly    bool          `json:"readOnly"`
	MultipleOf  interface{}   `json:"multipleOf"`
	Forms       []interface{} `json:"forms"`
	Value       interface{}   `json:"value"`
}

func NewPropertyFromString(dev device, propertyDesc string) *Property {
	data := []byte(propertyDesc)
	p := Property{}
	p.device = dev
	p.Name = json.Get(data, "name").ToString()
	p.Type = json.Get(data, "type").ToString()
	p.AtType = json.Get(data, "@type").ToString()
	p.Unit = json.Get(data, "unit").ToString()
	p.Description = json.Get(data, "description").ToString()
	p.Minimum = json.Get(data, "minimum").GetInterface()
	p.Maximum = json.Get(data, "maximum").GetInterface()
	var e []interface{}
	json.Get(data, "enum").ToVal(&e)
	p.MultipleOf = json.Get(data, "multipleOf").GetInterface()
	var f []interface{}
	json.Get(data, "forms").ToVal(&f)
	p.Forms = f
	if p.Name == "" || p.Type == "" {
		return nil
	}
	return &p
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetTitle() string {
	return p.Title
}

func (p *Property) GetType() string {
	return p.Type
}

func (p *Property) GetAtType() string {
	return p.AtType
}

func (p *Property) GetUnit() string {
	return p.Unit
}

func (p *Property) GetDescription() string {
	return p.Description
}

func (p *Property) GetMinimum() interface{} {
	return p.Minimum
}
func (p *Property) GetMaximum() interface{} {
	return p.Maximum
}

func (p *Property) IsReadOnly() bool {
	return p.ReadOnly
}

func (p *Property) GetMultipleOf() interface{} {
	return p.MultipleOf
}

func (p *Property) GetForms() []interface{} {
	return p.Forms
}

func (p *Property) GetValue() interface{} {
	return p.Value
}

func (p *Property) SetValue(value interface{}) error {
	if p.ReadOnly {
		return fmt.Errorf("read-only Prop")
	}
	var numberValue = to.Float64(value)
	if p.Minimum != nil {
		if to.Float64(p.Minimum) > numberValue {
			return fmt.Errorf("value less than minimum: %s", p.Minimum)
		}
	}
	if p.Maximum != nil {
		if to.Float64(p.Maximum) < numberValue {
			return fmt.Errorf("value greater than minimum: %s", p.Maximum)
		}
	}
	if p.MultipleOf != nil {
		if numberValue/to.Float64(p.MultipleOf)-math.Round(numberValue/to.Float64(p.MultipleOf)) != 0 {
			return fmt.Errorf("value is not a multiple of : %s", p.MultipleOf)
		}
	}
	if len(p.Enum) > 0 {
		for e := range p.Enum {
			if e == value {
				break
			}
			return fmt.Errorf("invalid enum value")
		}
	}
	p.setCachedValueAndNotify(value)
	return nil
}

func (p *Property) setCachedValueAndNotify(value interface{}) bool {
	oldValue := p.Value
	p.SetCachedValue(value)
	var hasChanged = oldValue != p.Value
	if hasChanged {
		p.device.NotifyValueChanged(p)
	}
	return hasChanged
}

func (p *Property) SetCachedValue(value interface{}) interface{} {
	if p.Type == TypeBoolean {
		p.Value = !!to.Bool(value)
	} else {
		p.Value = value
	}
	return p.Value
}
