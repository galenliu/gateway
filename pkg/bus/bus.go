package bus

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/bus/topic"
)

type ThingsBus interface {
	Publisher
	Subscriber
}

type Publisher interface {
	Publish(topic2 topic.Topic, args ...any)
}

type Subscriber interface {
	Subscribe(topic2 topic.Topic, f any) func()
	Unsubscribe(topic2 topic.Topic, f any)
}

type EventBus struct {
	bus Bus
}

func NewBus() *EventBus {
	return &EventBus{bus: New()}
}

func (t *EventBus) Subscribe(topic topic.Topic, fn any) func() {
	top := string(topic)
	err := t.bus.Subscribe(top, fn)
	if err != nil {
		fmt.Printf("bus subscription error: %s \t\n", err.Error())
		return func() {
			fmt.Printf("bus unsubscribe error: %s \t\n", err.Error())
		}
	}
	return func() {
		if t.bus.HasCallback(string(topic)) {
			err := t.bus.Unsubscribe(top, fn)
			if err != nil {
				fmt.Printf("bus unsubscribe error: %s \t\n", err.Error())
				return
			}
		}
	}
}

func (t *EventBus) Publish(topic topic.Topic, args ...any) {
	top := string(topic)
	go t.bus.Publish(top, args...)
}

func (t *EventBus) Unsubscribe(topic topic.Topic, f any) {
	top := string(topic)
	err := t.bus.Unsubscribe(top, f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
