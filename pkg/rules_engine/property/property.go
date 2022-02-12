package property

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
)

type Property struct {
	container things.Container
	*bus.Controller
	id          string
	typ         string
	thing       string
	unit        string
	description string
	href        string
	cleanUp     []func()
}

type Description struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Thing       string `json:"thing"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
	Href        string `json:"href,omitempty"`
}

func NewProperty(description Description, container things.Container) *Property {
	p := &Property{
		container:   container,
		id:          description.Id,
		typ:         description.Type,
		thing:       description.Thing,
		unit:        description.Unit,
		description: description.Description,
		href:        description.Href,
		cleanUp:     make([]func(), 1),
	}
	return p
}

func (p *Property) onPropertyChanged(property properties.Entity) {
	if property.GetDevice().GetId() == p.thing && property.GetName() == p.id {
		p.Publish(topic.ValueChanged, property.GetDevice().GetId(), property.GetName(), property.GetCachedValue())
	}
}

func (p *Property) onThingAdded(thingId string) {
	if p.thing == thingId {
		err := p.getInitialValue()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (p *Property) Start() {
	p.cleanUp = append(p.cleanUp, p.container.Subscribe(topic.DevicePropertyChanged, p.onPropertyChanged))
	err := p.getInitialValue()
	if err != nil {
		p.cleanUp = append(p.cleanUp, p.container.Subscribe(topic.ThingAdded, p.onThingAdded))
		return
	}
}

func (p *Property) Stop() {
	for _, f := range p.cleanUp {
		f()
	}
}

func (p *Property) get() (any, error) {
	return p.container.GetThingPropertyValue(p.thing, p.id)
}

func (p *Property) getInitialValue() error {
	v, err := p.get()
	if err != nil {
		return err
	}
	p.Publish(topic.ValueChanged, p.thing, p.id, v)
	return nil
}

func (p *Property) Set(value any) (any, error) {

	v, err := p.container.SetThingPropertyValue(p.thing, p.id, value)
	if err != nil {
		v, err = p.container.SetThingPropertyValue(p.thing, p.id, value)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return v, err
}

func (p *Property) GetThing() string {
	return p.thing
}

func (p *Property) ToDescription() Description {
	return Description{
		Id:          p.id,
		Type:        p.typ,
		Thing:       p.thing,
		Unit:        p.unit,
		Description: p.description,
		Href:        p.href,
	}
}

func (p *Property) GetType() string {
	return p.typ
}
