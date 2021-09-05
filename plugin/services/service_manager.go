package services

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models"
	"github.com/galenliu/gateway/server/models/model"
)

type ServiceManager struct {
	serviceInfo *model.ServiceInfo
	logger      logging.Logger
}

func NewServiceManager(s models.ServiceStorage, log logging.Logger) *ServiceManager {
	m := &ServiceManager{}
	m.logger = log
	return m
}

func (m ServiceManager) EnableService(id string) error {
	panic("implement me")
}

func (m ServiceManager) DisableService(id string) error {
	panic("implement me")
}

func (m ServiceManager) LoadService(id string) error {
	panic("implement me")
}

func (m ServiceManager) LoadServices() {

}
