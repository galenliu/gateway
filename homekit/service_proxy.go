package homekit

import "github.com/brutella/hc/service"

type ServiceProxy interface {
	GetMapOfPropertyProxy()map[string]PropertyProxy
}


type _service struct {
	*service.Service
	propertiesProxy map[string]PropertyProxy
}

func (s *_service) GetMapOfPropertyProxy()map[string]PropertyProxy {
	return s.propertiesProxy
}
