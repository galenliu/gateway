package models

import (
	"github.com/galenliu/gateway/server/models/model"
	json "github.com/json-iterator/go"
)

type ServiceStorage interface {
	LoadServiceConfig(string2 string) (string, error)
	StoreServiceConfig(id, conf string) error
	LoadServiceSetting(key string) (value string, err error)
	StoreServiceSetting(key, value string) error
}

type Services struct {
	services map[string]model.Service
	manager  model.ThingsManager
	storage  ServiceStorage
}

func NewServicesModel(m model.ThingsManager) *Services {
	s := &Services{}
	s.manager = m
	return s
}

func (s *Services) AddService(ser model.Service) {
	s.services[ser.GetID()] = ser
}

func (s *Services) RemoveService(ser model.Service) {
	delete(s.services, ser.GetID())
}

func (s *Services) NewThingAdded(data []byte) {
	for _, ser := range s.GetServices() {
		ser.OnNewThingAdded(data)
	}
}

func (s *Services) PropertyChanged(data []byte) {
	for _, ser := range s.GetServices() {
		ser.OnPropertyChanged(data)
	}
}

func (s Services) NotifyAction(data []byte) {
	for _, ser := range s.GetServices() {
		ser.OnAction(data)
	}
}

func (s *Services) GetServices() []model.Service {
	return nil
}

func (s *Services) LoadConfig(id string) (*model.ServiceInfo, error) {
	conf, err := s.storage.LoadServiceConfig(id)
	if err != nil {
		return nil, err
	}
	var info model.ServiceInfo
	err = json.Unmarshal([]byte(conf), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (s *Services) StoreConfig(id string, config string) (*model.ServiceInfo, error) {
	err := s.storage.StoreServiceConfig(id, config)
	if err != nil {
		return nil, err
	}
	var info model.ServiceInfo
	err = json.Unmarshal([]byte(config), &info)
	if err != nil && &info == nil {
		return nil, err
	}
	return &info, nil
}
