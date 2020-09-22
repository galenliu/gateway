package accessory

import (
	"gateway/logger"
)

type Container struct {
	Accessories   []*Accessory `json:"accessories"`
	AccessoryMaps map[uint64]*Accessory
}

func NewContainer() *Container {
	return &Container{
		Accessories:   make([]*Accessory, 0),
		AccessoryMaps: map[uint64]*Accessory{},
	}

}

func (c *Container) AddAccessory(a *Accessory) error {

	if c.AccessoryMaps[a.ID] != nil {
		logger.Warning.Printf("Accessory Exist,ID: &v Type: &v", a.ID, a.Type)
		return nil
	}
	c.Accessories = append(c.Accessories, a)
	c.AccessoryMaps[a.ID] = a
	logger.Info.Printf("Add Accessory,ID: &v Type: &v", a.ID, a.Type)
	return nil
}

//RemoveAccessory removes accessory form accessories
func (c *Container) RemoveAccessory(a *Accessory) {
	for i, accessory := range c.Accessories {
		if accessory == a {
			c.Accessories = append(c.Accessories[0:i], c.Accessories[i+1:]...)
			logger.Info.Printf("Remove Accessory,ID: &v Type: &v", a.ID, a.Type)
		}
	}

}
