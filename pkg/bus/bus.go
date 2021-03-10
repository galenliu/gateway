package bus

import (
	"fmt"
	"gateway/pkg/log"
	"github.com/asaskevich/EventBus"
	"sync"
)

var instance EventBus.Bus
var once sync.Once

func initBus() {
	once.Do(
		func() {
			instance = EventBus.New()
		},
	)

}

func Subscribe(topic string, fn interface{}) error {

	if instance == nil {
		initBus()
	}
	return instance.Subscribe(topic, fn)
}

func Unsubscribe(topic string, fn interface{}) error {
	if instance == nil {
		initBus()
	}
	return instance.Unsubscribe(topic, fn)
}

func Publish(topic string, args ...interface{}) {
	log.Info("publish topic: " + topic)
	if instance == nil {
		initBus()
	}
	log.Info(fmt.Sprintf(topic+" has callback %v", instance.HasCallback(topic)))
	if !instance.HasCallback(topic) {
		return
	}
	instance.Publish(topic, args...)
}

func HasCallBack(topic string) bool {
	return instance.HasCallback(topic)
}

func SubscribeOnce(topic string, fn interface{}) error {
	if instance == nil {
		initBus()
	}
	return instance.SubscribeOnce(topic, fn)
}

func SubscribeAsync(topic string, fn interface{}) error {
	if instance == nil {
		initBus()
	}
	return instance.SubscribeAsync(topic, fn, false)
}

func WaitAsync() {
	if instance == nil {
		initBus()
	}
	instance.WaitAsync()
}

//type OnPropertyChangeFunc = func(property *addon.Property)
//
//func SubscribePropertyChanged(thingId, propName string, fn OnPropertyChangeFunc) error {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, thingId, propName)
//	if instance == nil {
//		initBus()
//	}
//	return instance.Subscribe(topic, fn)
//}
//
//func SubscribeOncePropertyChanged(thingId, propName string, fn OnPropertyChangeFunc) error {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, thingId, propName)
//
//	if instance == nil {
//		initBus()
//	}
//	return instance.SubscribeOnce(topic, fn)
//}
//
//func UnsubscribePropertyChanged(thingId, propName string, fn OnPropertyChangeFunc) error {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, thingId, propName)
//	if instance == nil {
//		initBus()
//	}
//	return instance.Unsubscribe(topic, fn)
//}
//
//func PublishPropertyChanged(prop *addon.Property) {
//	topic := fmt.Sprintf("%s.%s.%s", util.PropertyChanged, prop.DeviceId, prop.Name)
//	if instance == nil {
//		initBus()
//	}
//	log.Info(fmt.Sprintf(topic+" has callback %v", instance.HasCallback(topic)))
//	if !instance.HasCallback(topic) {
//		return
//	}
//	instance.Publish(topic, prop)
//}
