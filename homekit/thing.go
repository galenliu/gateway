package homekit

import (
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/server/models"
)

type homekitService struct {
	id string
}

func (s *homekitService) GetID() string {
	return s.id
}

type LightBulb struct {
	*homekitService
	*service.Lightbulb
}

func (l *LightBulb) GetService() *service.Service {
	return l.Service
}

func NewLightBulb(t *models.Thing) *LightBulb {
	light := &LightBulb{}
	light.Lightbulb = service.NewLightbulb()
	return nil
}

type Thing struct {
	id string
	*service.Service
}

func NewHomekitService(thing *models.Thing) HService {

	if thing.SelectedCapability == "" {
		return nil
	}
	switch thing.SelectedCapability {
	case Light:
		return NewLightBulb(thing)
	}
	return nil
}
