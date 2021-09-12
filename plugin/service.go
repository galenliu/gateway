package models

import (
	"github.com/galenliu/gateway/pkg/rpc"
)

type Plugin interface {
	SendMsg(rpc.MessageType, map[string]interface{})
}

type Bus interface {
	SubscribePropertyChanged(f func(property *Property))
	UnsubscribePropertyChanged(f func(property *Property))
}

type Service struct {
	ID     string
	Name   string
	bus    Bus
	plugin Plugin
}

func NewService(plugin Plugin, bus Bus, id string, name string) *Service {
	s := &Service{}
	s.bus = bus
	s.plugin = plugin
	s.ID = id
	s.Name = name
	s.bus = bus
	s.bus.SubscribePropertyChanged(s.handlePropertyChanged)
	return s
}

func (s *Service) handlePropertyChanged(property *Property) {
	data := make(map[string]interface{})
	data["thingId"] = property.Device.GetID()
	data["propertyName"] = property.Name
	data["value"] = property.Value
	s.sendMsg(rpc.MessageType_ServicePropertyChangedNotification, data)
}

func (s *Service) sendMsg(messageType rpc.MessageType, data map[string]interface{}) {
	data["serviceId"] = s.ID
	s.plugin.SendMsg(messageType, data)
}

func (s *Service) unload() {
	s.bus
}
