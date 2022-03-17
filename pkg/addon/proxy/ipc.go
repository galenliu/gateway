package proxy

import (
	"context"
	"fmt"
	"github.com/fasthttp/websocket"
	json "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type handler interface {
	OnMessage([]byte)
}

type IpcClient struct {
	handler  handler
	url      string
	sendLock sync.Mutex
	origin   string
	verbose  bool

	readMsgChan chan []byte
	sendMsgChan chan []byte
}

// NewClient 新建一个Client，注册消息Handler
func NewClient(ctx context.Context, handler handler, path string) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + path, Path: "/"}
	client := &IpcClient{}
	client.readMsgChan = make(chan []byte, 10)
	client.sendMsgChan = make(chan []byte, 10)
	client.handler = handler
	client.url = u.String()
	go client.connection(ctx)
	go func(ctx context.Context) {
		for {
			select {
			case data := <-client.readMsgChan:
				client.handler.OnMessage(data)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
	return client
}

func (client *IpcClient) Send(message any) {
	client.sendLock.Lock()
	defer client.sendLock.Unlock()
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf(err.Error())
	}
	select {
	case client.sendMsgChan <- data:
	default:
		fmt.Printf("send channel is full")
	}
}

func (client *IpcClient) connection(ctx context.Context) {
	for {
		var done = make(chan struct{})
		var err error
		var resp *http.Response
		var ws *websocket.Conn
		ws, resp, err = websocket.DefaultDialer.Dial(client.url, nil)
		if !resp.Close || err == nil {
			defer ws.Close()
			go client.sendLoop(ws, done)
			client.readLoop(ws, ctx)
		}
		if err != nil {
			fmt.Printf("websocket connection err:%v\n", err.Error())
			break
		}
		done <- struct{}{}
		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (client *IpcClient) readLoop(conn *websocket.Conn, ctx context.Context) {

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("websocket err: %s", err.Error())
			break
		}
		select {
		case <-ctx.Done():
			return
		case client.readMsgChan <- message:
		default:
			fmt.Println("read channel is full")
		}
	}
}

func (client *IpcClient) sendLoop(conn *websocket.Conn, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case msg := <-client.sendMsgChan:
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
