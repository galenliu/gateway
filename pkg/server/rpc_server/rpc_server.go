package rpc_server

import (
	"context"
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
	port   string
	logger logging.Logger
	ctx    context.Context
	rpc.UnimplementedPluginServerServer
	pluginSever server.PluginServer
	userProfile *rpc.UsrProfile
}

func NewRPCServer(ctx context.Context, pluginServer server.PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *RPCServer {
	s := &RPCServer{}
	s.pluginSever = pluginServer
	s.port = port
	s.userProfile = userProfile
	s.ctx = ctx
	s.logger = log
	go s.Run()
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
		case <-s.ctx.Done():

			return fmt.Errorf("rpc server stopped")
		}

	}

}

func (s *RPCServer) Run() {

	lis, err := net.Listen("tcp", s.port)
	s.logger.Infof("RPC server run addr: %s", s.port)
	if err != nil {
		s.logger.Error(err.Error())
	}
	sev := grpc.NewServer()
	rpc.RegisterPluginServerServer(sev, s)
	for {
		err = sev.Serve(lis)
		if err != nil {
			s.logger.Error(err.Error())
		}

		if err != nil {
			s.logger.Errorf("RPC Start err: %s", err)
		}
		select {
		case <-s.ctx.Done():
			sev.Stop()
		}

	}

}
