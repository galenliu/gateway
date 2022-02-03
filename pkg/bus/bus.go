package bus

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/galenliu/gateway/pkg/bus/topic"
)

type Bus interface {
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

type Controller struct {
	bus EventBus.Bus
}

func NewBusController() *Controller {
	return &Controller{bus: EventBus.New()}
}

func (t *Controller) Subscribe(topic topic.Topic, fn any) func() {
	err := t.bus.Subscribe(string(topic), fn)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return func() {
		err := t.bus.Unsubscribe(string(topic), fn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func (t *Controller) Publish(topic2 topic.Topic, args ...any) {
	t.bus.Publish(string(topic2), args...)
}

func (t *Controller) Unsubscribe(topic2 topic.Topic, f any) {
	err := t.bus.Unsubscribe(string(topic2), f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
