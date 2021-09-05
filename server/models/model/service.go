package model

import (
	"github.com/galenliu/gateway/pkg/container"
)

type ThingsManager interface {
	SetPropertyValue(thingId, propertyName string, value interface{}) (interface{}, error)
	GetPropertyValue(thingId, propertyName string) (interface{}, error)
	GetPropertiesValue(thingId string) (map[string]interface{}, error)
}

// Container  Things
type Container interface {
	GetThing(id string) *container.Thing
	GetThings() []*container.Thing
	GetMapThings() map[string]*container.Thing
	CreateThing(data []byte) (*container.Thing, error)
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

type ServiceInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
