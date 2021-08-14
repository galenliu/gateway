package rpc_server

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/plugin"
	json "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"net"
)

type PluginHandler interface {
	MessageHandler(mt rpc.MessageType, data []byte) error
}

type RPCServer struct {
	port     string
	logger   logging.Logger
	doneChan chan struct{}
	rpc.UnimplementedPluginServerServer
	pluginSever *plugin.PluginsServer
	userProfile []byte
	preferences []byte
}

func NewRPCServer(server *plugin.PluginsServer, port string, userProfile []byte, preferences []byte, log logging.Logger) *RPCServer {
	s := &RPCServer{}
	s.pluginSever = server
	s.port = port
	s.userProfile = userProfile
	s.preferences = preferences
	s.doneChan = make(chan struct{})
	s.logger = log
	return s
}

func (s *RPCServer) PluginHandler(p rpc.PluginServer_PluginHandlerServer) error {
	r, err := p.Recv()
	if err != nil {
		return err
	}
	message := rpc.PluginRegisterRequestMessage{
		MessageType: 0,
		Data:        &rpc.PluginRegisterRequestMessageDataTemp{PluginId: json.Get(r.Data, "pluginId").ToString()},
	}

	if message.Data.PluginId == "" {
		return fmt.Errorf("plugin id faild")
	}
	err = p.SendMsg(rpc.PluginRegisterResponseMessage{
		MessageType: 0,
		Data: &rpc.PluginRegisterResponseMessageGatewayConfig{
			PluginId:       "",
			GatewayVersion: "",
			UserProfile:    s.userProfile,
			Preferences:    s.preferences,
		},
	})
	if err != nil {
		return err
	}

	clint := NewClint(message.Data.PluginId, p)
	var pluginHandler PluginHandler
	pluginHandler = s.pluginSever.RegisterPlugin(message.Data.PluginId, clint)

	for {
		baseMessage, err := clint.Read()
		if err != nil {
			return err
		}
		err = pluginHandler.MessageHandler(baseMessage.MessageType, baseMessage.Data)
		if err != nil {
			return err
		}
		select {
		case <-s.doneChan:
			return fmt.Errorf("rpc server stopped")
		}

	}

}

func (s *RPCServer) Start() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	sev := grpc.NewServer()
	rpc.RegisterPluginServerServer(sev, s)
	err = sev.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *RPCServer) Stop() error {
	s.doneChan <- struct{}{}
	return nil
}
