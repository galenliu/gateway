package proxy

import (
	"fmt"
	"github.com/fasthttp/websocket"
	"net/url"
	"sync"
)

type handler interface {
	OnMessage([]byte)
}

type IpcClient struct {
	ws       *websocket.Conn
	handler  handler
	url      string
	sendLock sync.Mutex
	origin   string
	verbose  bool
}

// NewClient 新建一个Client，注册消息Handler
func NewClient(handler handler, path string) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + path, Path: "/"}
	client := &IpcClient{}
	client.handler = handler
	client.url = u.String()
	var err error = nil
	client.ws, _, err = websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		fmt.Printf("dial err: %s \r\n", err.Error())
		return nil
	}
	go client.readLoop()
	return client
}

func (client *IpcClient) Send(message any) {
	client.sendLock.Lock()
	defer client.sendLock.Unlock()
	err := client.ws.WriteJSON(message)
	if err != nil {
		fmt.Printf("client.Send err: %s \r\n", err.Error())
		return
	}
}

func (client *IpcClient) readLoop() {
	for {
		_, message, err := client.ws.ReadMessage()
		if err != nil {
			fmt.Printf("read messages error: %s", err.Error())
			return
		}
		client.handler.OnMessage(message)
	}
}

func (client *IpcClient) close() {
	if client.ws != nil {
		_ = client.ws.Close()
	}
}
