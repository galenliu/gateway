package ipc

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/logging"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	port   string
	logger logging.Logger
	ctx    context.Context
	rpc.UnimplementedPluginServerServer
	pluginSever PluginServer
	userProfile *rpc.UsrProfile
}

func NewRPCServer(ctx context.Context, pluginServer PluginServer, port string, userProfile *rpc.UsrProfile, log logging.Logger) *RPCServer {
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

	clint := &rpcClint{p}

	pluginHandler, err := s.pluginSever.RegisterPlugin(clint)

	if err != nil {
		return err
	}

	for {
		baseMessage, err := clint.ReadMessage()
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

type rpcClint struct {
	rpc.PluginServer_PluginHandlerServer
}

func (r *rpcClint) WriteMessage(message *rpc.BaseMessage) error {
	return r.PluginServer_PluginHandlerServer.Send(message)
}

func (r *rpcClint) ReadMessage() (*rpc.BaseMessage, error) {
	return r.Recv()
}
