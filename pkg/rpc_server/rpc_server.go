package rpc_server

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"google.golang.org/grpc"
	"io"
	"net"
)

type RPCServer struct {
	port     string
	logger   logging.Logger
	sendChan chan *rpc.BaseMessage
	doneChan chan struct{}
	rpcClint rpc.PluginServer_RegisterServer
	rpc.UnimplementedPluginServerServer
}

func NewRPCServer(log logging.Logger) *RPCServer {
	s := &RPCServer{}
	s.sendChan = make(chan *rpc.BaseMessage)
	s.doneChan = make(chan struct{})
	s.logger = log
	return s
}

func (s *RPCServer) Register(p rpc.PluginServer_RegisterServer) error {
	s.rpcClint = p
	message, err := p.Recv()
	if err != nil {
		return err
	}
	if message.MessageType != rpc.MessageType_PluginRegisterRequest {
		return fmt.Errorf("RegisterRequest message type err")
	}
	err = p.SendMsg(rpc.PluginRegisterResponseMessage{
		MessageType: 0,
		Data: &rpc.PluginRegisterResponseMessageGatewayConfig{
			PluginId:       "",
			GatewayVersion: "",
			UserProfile:    nil,
			Preferences:    nil,
		},
	})
	if err != nil {
		return err
	}
	go s.Read()

	for {
		select {
		case m := <-s.sendChan:
			err := s.rpcClint.Send(m)
			if err != nil {
				return err
			}
		case <-s.doneChan:
			return nil
		}
	}
}

func (s *RPCServer) Read() {
	var revMsg interface{}
	err := s.rpcClint.RecvMsg(&revMsg)
	if err == io.EOF {
		s.logger.Info("plugin is closed")
		return
	}
	if err != nil {
		return
	}
}

func (s *RPCServer) Send(message *rpc.BaseMessage) {
	s.sendChan <- message
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
