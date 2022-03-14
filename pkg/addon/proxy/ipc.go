package proxy

import (
	"fmt"
	"github.com/fasthttp/websocket"
	json "github.com/json-iterator/go"
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
	done     chan struct{}
}

// NewClient 新建一个Client，注册消息Handler
func NewClient(handler handler, path string) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + path, Path: "/"}
	client := &IpcClient{}
	client.handler = handler
	client.done = make(chan struct{})
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
	data, _ := json.Marshal(message)
	err := client.ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Printf("client.send err: %s \r\n", err.Error())
		return
	}
}

func (client *IpcClient) readLoop() {

	for {
		select {
		case <-client.done:
			return
		default:
			_, message, err := client.ws.ReadMessage()
			if err != nil {
				fmt.Printf("read messages error: %s", err.Error())
				return
			}
			client.handler.OnMessage(message)
		}
	}
}

func (client *IpcClient) close() {
	select {
	case client.done <- struct{}{}:
	}
	if client.ws != nil {
		_ = client.ws.Close()
	}
}
