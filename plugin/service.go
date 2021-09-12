package plugin

import (
	"github.com/galenliu/gateway-grpc"
)

type Bus interface {
	SubscribePropertyChanged(f func(property *gateway_grpc.DevicePropertyChangedNotificationMessage_Data))
	UnsubscribePropertyChanged(f func(property *gateway_grpc.DevicePropertyChangedNotificationMessage_Data))
	SubscribeActionStatus(f func(action *gateway_grpc.ActionDescription))
	UnsubscribeActionStatus(f func(action *gateway_grpc.ActionDescription))
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

func (s *Service) handlePropertyChanged(property *gateway_grpc.DevicePropertyChangedNotificationMessage_Data) {
	data := make(map[string]interface{})
	data["thingId"] = property.DeviceId
	data["propertyName"] = property.Property.Name
	data["value"] = property.Property.Value
	s.sendMsg(gateway_grpc.MessageType_ServicePropertyChangedNotification, data)
}

func (s *Service) handleActionStatus(action *gateway_grpc.ActionDescription) {
	data := make(map[string]interface{})
	s.sendMsg(gateway_grpc.MessageType_ServicePropertyChangedNotification, data)
}

func (s *Service) sendMsg(messageType gateway_grpc.MessageType, data map[string]interface{}) {
	data["serviceId"] = s.ID
	s.plugin.SendMsg(messageType, data)
}

func (s *Service) unload() {
	s.bus.UnsubscribePropertyChanged(s.handlePropertyChanged)
	s.bus.UnsubscribeActionStatus(s.handleActionStatus)
}
