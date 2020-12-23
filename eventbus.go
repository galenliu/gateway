package gateway

type OnPropertyChangedFunc func(name string, value interface{})
type OnActionStatusFunc func()
type RemoveFunc func()

var eventBus *EventsBus

type EventsBus struct {
	eventId                   int64
	OnPropertyChangedListener map[int64]OnPropertyChangedFunc
}

func NewEventBus() *EventsBus {
	eventBus = &EventsBus{
		eventId:                   0,
		OnPropertyChangedListener: make(map[int64]OnPropertyChangedFunc),
	}
	return eventBus
}

func (e *EventsBus) listenPropertyChanged(f OnPropertyChangedFunc) func() {
	e.eventId++
	e.OnPropertyChangedListener[e.eventId] = f
	return func() {
		delete(e.OnPropertyChangedListener, e.eventId)
	}
}

func (e *EventsBus) onPropertyChanged(thingId, propName string, value interface{}) {
	for _, f := range e.OnPropertyChangedListener {
		f(propName, value)
	}
}

func ListenPropertyChange(f OnPropertyChangedFunc) func() {
	eventBus.eventId++
	eventBus.OnPropertyChangedListener[eventBus.eventId] = f
	return func() {
		delete(eventBus.OnPropertyChangedListener, eventBus.eventId)
	}
}

func OnPropertyChanged(thingId, propName string, value interface{}) {
	for _, f := range eventBus.OnPropertyChangedListener {
		f(propName, value)
	}
}
