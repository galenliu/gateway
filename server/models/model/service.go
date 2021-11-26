package model

import (
	"github.com/galenliu/gateway/pkg/container"
)

// Container  Things
type Container interface {
	GetThing(id string) *container.Thing
	GetThings() []*container.Thing
	GetMapOfThings() map[string]*container.Thing
	CreateThing(data []byte) (*container.Thing, error)
	RemoveThing(id string) error
	UpdateThing(data []byte) error
}
