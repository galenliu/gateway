package bus

import (
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
