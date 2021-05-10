package homekit

import (
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/server/models"
)

type ServiceProxy struct {
	*models.Thing
	*service.Service
	characteristics map[string]HCharacteristic
}

func NewService(thing *models.Thing) *ServiceProxy {
	s := &ServiceProxy{}
	s.Thing = thing
	s.characteristics = make(map[string]HCharacteristic)
	return nil
}

func (s *ServiceProxy) GetService() *service.Service {
	return s.Service
}

func (s *ServiceProxy) GetHCharacteristic(name string) HCharacteristic {
	c, ok := s.characteristics[name]
	if ok {
		return c
	}
	return nil
}

//func NewLightBulb(t *models.Thing) *homekitService {
//	light := &LightBulb{}
//	light.id = t.ID
//	light.Lightbulb = service.NewLightbulb()
//	return nil
//}

// NewHomekitService return ServiceProxy implement HService
func NewHomekitService(thing *models.Thing) *ServiceProxy {
	s := NewService(thing)
	if thing.SelectedCapability == "" {
		return nil
	}
	switch thing.SelectedCapability {
	case Light:
		light := service.NewLightbulb()
		s.Service = light.Service
		for _, p := range thing.Properties {
			char := NewCharacteristicProxy(p)
			if char == nil {
				continue
			}
			s.characteristics[char.GetName()] = char
			light.AddCharacteristic(char.GetCharacteristic())
			return s
		}
	default:
		return nil

	}
	return nil
}
