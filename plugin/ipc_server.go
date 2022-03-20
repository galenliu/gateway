package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/websocket"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8172
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	registered bool
	conn       *websocket.Conn
}

var clientChan chan *client
var errChan chan string

func NewIpcServer(ctx context.Context, addr string) (chan *client, chan string) {

	clientChan = make(chan *client, 64)
	errChan = make(chan string)
	http.HandleFunc("/", serveWs)
	srv := http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}
	var wg sync.WaitGroup
	go func() {
		<-ctx.Done()
		wg.Add(1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("error shutting down:%s", err.Error())
		}
		wg.Done()
		close(errChan)
		close(clientChan)
	}()
	go func() {
		log.Println("listening at " + addr)
		err := srv.ListenAndServe()
		fmt.Println("waiting for the remaining connections to finish...")
		wg.Wait()
		if err != nil && err != http.ErrServerClosed {
			close(clientChan)
			select {
			case errChan <- err.Error():
			}
		}
		log.Println("gracefully shutdown the http server...")
	}()
	return clientChan, errChan
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	select {
	case clientChan <- &client{
		registered: false,
		conn:       ws,
	}:
	}
}

type message struct {
	MessageType messages.MessageType `json:"messageType"`
	Data        any                  `json:"data"`
}

func (c *client) sendMsg(mt messages.MessageType, data any) error {
	m := message{
		MessageType: mt,
		Data:        data,
	}
	byt, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return c.write(byt)
}

func (c *client) read() ([]byte, error) {

	//err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//if err != nil {
	//	return nil, err
	//}
	_, data, err := c.conn.ReadMessage()
	data = bytes.Trim(data, "\n")
	data = bytes.Trim(data, " ")
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("websocket close error: %v", err)
			return nil, err
		}
		return nil, err
	}
	return data, nil
}

func (c *client) write(data []byte) error {
	data = bytes.Trim(data, "\n")
	data = bytes.Trim(data, " ")
	err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}
