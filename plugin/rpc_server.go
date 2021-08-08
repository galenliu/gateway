package plugin

import (
    "fmt"
    "github.com/galenliu/gateway/pkg/logging"
    "google.golang.org/grpc"
    "net"
)


type RpcServer struct {
      port string
      logger logging.Logger
}

func (rpc *RpcServer )NewRpcServer(port string,log logging.Logger) *RpcServer {
    s := &RpcServer{}
    if port == ""{
        return nil
    }
    s.port = port
    s.logger = log
   return s
}

func (rpc *RpcServer) Start() error {
    lis,err := net.Listen("tcp",rpc.port)
    if err != nil{
        return fmt.Errorf("rpc server listen err: %s",err.Error())
    }
   serv:= grpc.NewServer()
    err = serv.Serve(lis)
    if err != nil {
        return err
    }
    rpc.logger.Info("rpc listen at %s",rpc.port)
    return nil
}