package ipc_server

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IPCServer struct {
	logger      logging.Logger
	addr        string
	path        string
	Connections chan *Connection
	locker      *sync.Mutex
}

func NewIPCServer(port string, log logging.Logger) *IPCServer {
	ipc := &IPCServer{
		addr:        "localhost:" + port,
		Connections: make(chan *Connection, 5),
	}
	ipc.logger = log
	go func() {
		err := ipc.Start()
		if err != nil {
			ipc.logger.Warningf("IPC Server Err: %s", err.Error())
		}
	}()
	return ipc
}


func (s *IPCServer) Start() error {
	http.HandleFunc("/", s.handle)
	s.logger.Info("IPC server run addr: %s", s.addr)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		s.logger.Error("ipc s fail,err: %s", err.Error())
		return err
	}
	return nil
}

func (s *IPCServer) Close() error {
	close(s.Connections)
	return nil
}

func (s *IPCServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	s.logger.Debug("accept new connection")
	if conn == nil {
		return
	}
	//升级协议时可能发生的错误
	if err != nil {
		s.logger.Error("ipc s upgrade failed,err: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Connections <- NewConn(conn)
}
