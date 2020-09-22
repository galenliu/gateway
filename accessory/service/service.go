package service

import (
	"gateway/accessory/characteristic"
	"gateway/logger"
)

type Service struct {
	ID              uint64
	Type            string
	Characteristics []*characteristic.Characteristic
	Hidden          bool
	Primary         bool
	Linked          []*Service
	ServiceName     string
}

func New(typ string) *Service {
	s := Service{Type: typ,
		Characteristics: []*characteristic.Characteristic{},
		Linked:          []*Service{},
	}
	return &s
}

func (s *Service) GetCharacteristics() []*characteristic.Characteristic {
	var result []*characteristic.Characteristic
	for _, c := range s.Characteristics {
		result = append(result, c)
	}
	return result
}

func (s *Service) AddCharacteristics(c *characteristic.Characteristic) {
	s.Characteristics = append(s.Characteristics, c)
	logger.Info.Printf("characteristic:&v add service: & ", c.ID, s.ID)
}

func (s *Service) AddLinkService(other *Service) {
	s.Linked = append(s.Linked, other)

}
