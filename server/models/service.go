package models

type ThingsManager interface {
	SetPropertyValue(thingId, propertyName string, value interface{}) (interface{}, error)
	GetPropertyValue(thingId, propertyName string) (interface{}, error)
	GetPropertiesValue(thingId string) (map[string]interface{}, error)
}

// Container  Things
type Container interface {
	GetThing(id string) *Thing
	GetThings() []*Thing
	GetMapThings() map[string]*Thing
	CreateThing(data []byte) (*Thing, error)
	RemoveThing(id string) error
	UpdateThing(data []byte) error
}

type Service interface {
	GetID() string
	OnNewThingAdded([]byte)
	OnPropertyChanged([]byte)
	OnAction([]byte)
	SetPropertyValue(v interface{})

	NewService(id string, manager ThingsManager, c Container) Service
}
