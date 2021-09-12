package rpc_server

import (
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	json "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"net"
)

type PluginServer interface {
	RegisterPlugin(pluginId string, clint rpc.Clint) rpc.PluginHandler
}

type RPCServer struct {
	port     string
	logger   logging.Logger
	doneChan chan struct{}
	gateway_grpc.UnimplementedPluginServerServer
	pluginSever PluginServer
	userProfile *gateway_grpc.PluginRegisterResponseMessage_Data_UsrProfile
	preferences *gateway_grpc.PluginRegisterResponseMessage_Data_Preferences
}

func NewRPCServer(pluginServer PluginServer, port string, userProfile *gateway_grpc.PluginRegisterResponseMessage_Data_UsrProfile, preferences *gateway_grpc.PluginRegisterResponseMessage_Data_Preferences, log logging.Logger) *RPCServer {
	s := &RPCServer{}
	s.pluginSever = pluginServer
	s.port = port
	s.userProfile = userProfile
	s.preferences = preferences
	s.doneChan = make(chan struct{})
	s.logger = log
	return s
}

func (s *RPCServer) PluginHandler(p gateway_grpc.PluginServer_PluginHandlerServer) error {
	r, err := p.Recv()
	if err != nil {
		return err
	}
	pluginId := json.Get(r.Data, "pluginId").ToString()
	if r.MessageType != gateway_grpc.MessageType_PluginRegisterRequest || pluginId == "" {
		return err
	}

	err = p.SendMsg(gateway_grpc.PluginRegisterResponseMessage{
		MessageType: 0,
		Data: &gateway_grpc.PluginRegisterResponseMessage_Data{
			PluginId:       pluginId,
			GatewayVersion: constant.Version,
			UserProfile:    s.userProfile,
			Preferences:    s.preferences,
		},
	})
	if err != nil {
		return err
	}
	clint := NewClint(pluginId, p)
	var pluginHandler rpc.PluginHandler
	pluginHandler = s.pluginSever.RegisterPlugin(pluginId, clint)

	for {
		baseMessage, err := clint.Read()
		if err != nil {
			return err
		}
		err = pluginHandler.OnMsg(baseMessage.MessageType, baseMessage.Data)
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
	go func() {
		err := func() error {
			lis, err := net.Listen("tcp", s.port)
			s.logger.Infof("RPC server run addr: %s", s.port)
			if err != nil {
				return err
			}
			sev := grpc.NewServer()
			gateway_grpc.RegisterPluginServerServer(sev, s)
			err = sev.Serve(lis)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf("RPC Start err: %s", err)
		}
	}()
	return nil
}

func (s *RPCServer) Stop() error {
	s.doneChan <- struct{}{}
	return nil
}
