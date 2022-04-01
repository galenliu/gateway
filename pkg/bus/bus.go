package bus

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"reflect"
)

type ThingsBus interface {
	Publisher
	Subscriber
}

type Publisher interface {
	Publish(topic2 topic.Topic, args ...any)
}

type Subscriber interface {
	Subscribe(topic2 topic.Topic, f any) error
	Unsubscribe(topic2 topic.Topic, f any)
}

type EventBus struct {
	bus Bus
}

func NewBus() *EventBus {
	return &EventBus{bus: New()}
}

func (t *EventBus) Subscribe(topic topic.Topic, fn any) error {
	fmt.Printf("subscribe Topic: %s Pointer: %v Type: %v\n", topic, reflect.ValueOf(fn).Pointer(), reflect.ValueOf(fn).Type())
	top := string(topic)
	err := t.bus.Subscribe(top, fn)
	if err != nil {
		return err
	}
	return nil
}

func (t *EventBus) Publish(topic topic.Topic, args ...any) {
	top := string(topic)
	t.bus.Publish(top, args...)
}

func (t *EventBus) Unsubscribe(topic topic.Topic, f any) {
	fmt.Printf("unsubscribe Topic: %s Pointer: %v Type: %v\n", topic, reflect.ValueOf(f).Pointer(), reflect.ValueOf(f).Type())
	top := string(topic)
	err := t.bus.Unsubscribe(top, f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
