package proxy

import (
	"fmt"
	"github.com/fasthttp/websocket"
	"net/url"
	"sync"
)

const IpcDefaultPort = "9500"

type IpcClient struct {
	ws       *websocket.Conn
	manager  *Manager
	url      string
	sendLock sync.Mutex
	pluginId string
	origin   string
	verbose  bool
}

// NewClient 新建一个Client，注册消息Handler
func NewClient(PluginId string, manager *Manager) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + IpcDefaultPort, Path: "/"}
	client := &IpcClient{}
	client.pluginId = PluginId
	client.manager = manager
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

func (client *IpcClient) send(message any) {
	client.sendLock.Lock()
	defer client.sendLock.Unlock()
	err := client.ws.WriteJSON(message)
	if err != nil {
		fmt.Printf("client.send err: %s \r\n", err.Error())
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
		client.manager.onMessage(message)
	}
}

func (client *IpcClient) close() {
	if client.ws != nil {
		_ = client.ws.Close()
	}
}
