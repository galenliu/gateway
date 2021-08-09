package rpc_server

import (
	"github.com/galenliu/gateway/pkg/rpc"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	port string
}

func (s *RPCServer) Start() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	sev := grpc.NewServer()
	rpc.BaseMessage{}
	rpc.BaseMessage{}
	return nil
}
