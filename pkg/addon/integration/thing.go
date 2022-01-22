package integration

import things "github.com/galenliu/gateway/api/models/container"

type ThingEntity struct {
	*things.Thing
}

func NewThingEntity(thing *things.Thing) *ThingEntity {
	return &ThingEntity{
		thing,
	}
}

func (t *ThingEntity) GetProperty(id string) Property {
	return nil
}
