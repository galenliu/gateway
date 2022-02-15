package property

import (
	"fmt"
	things "github.com/galenliu/gateway/api/models/container"
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
		Controller:  bus.NewBusController(),
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

func (p *Property) onPropertyChanged(msg topic.ThingPropertyChangedMessage) {
	if msg.ThingId == p.thing && msg.PropertyName == p.id {
		p.Publish(topic.ValueChanged, msg.Value)
	}
}

func (p *Property) onThingAdded(message topic.ThingAddedMessage) {
	if p.thing == message.ThingId {
		err := p.getInitialValue()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (p *Property) onThingConnected(msg topic.ThingConnectedMessage) {
	if msg.ThingId == p.thing && msg.Connected {
		err := p.getInitialValue()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (p *Property) Start() {
	p.cleanUp = append(p.cleanUp, p.container.Subscribe(topic.ThingPropertyChanged, p.onPropertyChanged))
	err := p.getInitialValue()
	if err != nil {
		p.cleanUp = append(p.cleanUp, p.container.Subscribe(topic.ThingAdded, p.onThingAdded))
		p.cleanUp = append(p.cleanUp, p.container.Subscribe(topic.ThingConnected, p.onThingConnected))
		return
	}
}

func (p *Property) Stop() {
	for _, f := range p.cleanUp {
		f()
	}
}

func (p *Property) getInitialValue() error {
	v, err := p.container.GetThingPropertyValue(p.thing, p.id)
	if err != nil {
		return err
	}
	p.Publish(topic.ValueChanged, v)
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
