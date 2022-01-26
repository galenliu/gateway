package property

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/rules_engine"
)

type Bus interface {
	bus.Publisher
	bus.Subscriber
}

type ThingsHandler interface {
	SetThingProperty(thingId, propertyName string, value any) (any, error)
	GetThingProperty(thingId, propertyName string) (any, error)
}

type PropertyDescription struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Thing       string `json:"thing"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
	Href        string `json:"href"`
}

type Property struct {
	bus         Bus
	things      ThingsHandler
	id          string
	typ         string
	thing       string
	unit        string
	description string
	href        string
	cleanUp     []func()
}

func NewProperty(description PropertyDescription, bus Bus, handler ThingsHandler) *Property {
	p := &Property{
		bus:         bus,
		things:      handler,
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
		p.bus.Pub(rules_engine.ValueChanged, property.GetDevice().GetId(), property.GetName(), property.GetCachedValue())
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

func (p *Property) start() {
	p.cleanUp = append(p.cleanUp, p.bus.Sub(topic.DevicePropertyChanged, p.onPropertyChanged))
	err := p.getInitialValue()
	if err != nil {
		p.cleanUp = append(p.cleanUp, p.bus.Sub(topic.ThingAdded, p.onThingAdded))
		return
	}
}

func (p *Property) stop() {
	for _, f := range p.cleanUp {
		f()
	}
}

func (p *Property) get() (any, error) {
	return p.things.GetThingProperty(p.thing, p.id)
}

func (p *Property) getInitialValue() error {
	v, err := p.get()
	if err != nil {

		return err
	}
	p.bus.Pub(rules_engine.ValueChanged, p.thing, p.id, v)
	return nil
}

func (p *Property) Set(value any) (any, error) {
	v, err := p.things.SetThingProperty(p.thing, p.id, value)
	if err != nil {
		v, err = p.things.SetThingProperty(p.thing, p.id, value)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return v, err
}

func (p *Property) ToDescription() PropertyDescription {
	return PropertyDescription{
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
