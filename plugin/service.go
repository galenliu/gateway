package plugin

import (
	rpc "github.com/galenliu/gateway-grpc"
)

type Bus interface {
	SubscribePropertyChanged(f func(property *rpc.DevicePropertyChangedNotificationMessage_Data))
	UnsubscribePropertyChanged(f func(property *rpc.DevicePropertyChangedNotificationMessage_Data))
	SubscribeActionStatus(f func(action *rpc.ActionDescription))
	UnsubscribeActionStatus(f func(action *rpc.ActionDescription))
}

type Service struct {
	ID     string
	Name   string
	bus    Bus
	plugin *Plugin
}

func NewService(plugin *Plugin, bus Bus, id string, name string) *Service {
	s := &Service{}
	s.bus = bus
	s.plugin = plugin
	s.ID = id
	s.Name = name
	s.bus = bus
	s.bus.SubscribePropertyChanged(s.handlePropertyChanged)
	s.bus.SubscribeActionStatus(s.handleActionStatus)
	return s
}

func (s *Service) handlePropertyChanged(property *rpc.DevicePropertyChangedNotificationMessage_Data) {
	data := make(map[string]interface{})
	data["thingId"] = property.DeviceId
	data["propertyName"] = property.Property.Name
	data["value"] = property.Property.Value
	s.sendMsg(rpc.MessageType_ServicePropertyChangedNotification, data)
}

func (s *Service) handleActionStatus(action *rpc.ActionDescription) {
	data := make(map[string]interface{})
	s.sendMsg(rpc.MessageType_ServicePropertyChangedNotification, data)
}

func (s *Service) sendMsg(messageType rpc.MessageType, data map[string]interface{}) {
	data["serviceId"] = s.ID
	s.plugin.SendMsg(messageType, data)
}

func (s *Service) unload() {
	s.bus.UnsubscribePropertyChanged(s.handlePropertyChanged)
	s.bus.UnsubscribeActionStatus(s.handleActionStatus)
}
