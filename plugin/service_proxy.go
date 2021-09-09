package plugin

import "github.com/galenliu/gateway/pkg/rpc"

type Service struct {
	ID   string
	Name string

	plugin *Plugin
}

func NewServiceProxy(plugin *Plugin, id string, name string) *Service {
	s := &Service{}
	s.plugin = plugin
	s.ID = id
	s.Name = name
	return s
}

func (s *Service) handlePropertyChanged() {

}

func (s *Service) sendMsg(messageType rpc.MessageType, data map[string]interface{}) {
	data["serviceId"] = s.ID
	s.plugin.SendMsg(messageType, data)
}
