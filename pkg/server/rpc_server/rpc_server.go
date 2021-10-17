package rpc_server

import (
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/server"
	json "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	port     string
	logger   logging.Logger
	doneChan chan struct{}
	rpc.UnimplementedPluginServerServer
	pluginSever server.PluginServer
	userProfile *rpc.UsrProfile
}

func NewRPCServer(pluginServer server.PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *RPCServer {
	s := &RPCServer{}
	s.pluginSever = pluginServer
	s.port = port
	s.userProfile = userProfile
	s.doneChan = make(chan struct{})
	s.logger = log
	return s
}

func (s *RPCServer) PluginHandler(p rpc.PluginServer_PluginHandlerServer) error {
	r, err := p.Recv()
	if err != nil {
		return err
	}
	pluginId := json.Get(r.Data, "pluginId").ToString()
	if r.MessageType != rpc.MessageType_PluginRegisterRequest || pluginId == "" {
		return err
	}

	err = p.SendMsg(rpc.PluginRegisterResponseMessage{
		MessageType: 0,
		Data: &rpc.PluginRegisterResponseMessage_Data{
			PluginId:       pluginId,
			GatewayVersion: constant.Version,
			UserProfile:    s.userProfile,
			Preferences:    s.pluginSever.GetPreferences(),
		},
	})
	if err != nil {
		return err
	}
	clint := NewClint(pluginId, p)
	var pluginHandler server.PluginHandler
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

func (s *RPCServer) Run() error {
	err := func() error {
		lis, err := net.Listen("tcp", s.port)
		s.logger.Infof("RPC server run addr: %s", s.port)
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
	}()
	if err != nil {
		s.logger.Errorf("RPC Run err: %s", err)
	}
	return nil
}

func (s *RPCServer) Stop() error {
	s.doneChan <- struct{}{}
	return nil
}
