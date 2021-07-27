package transport

import (
	"fmt"
	"gateway/homekit/config"
	"github.com/gorilla/websocket"
	"net/url"
)

//主要用于 homekit组件和gateway数据传输

type Transport struct {
	ws  *websocket.Conn
	url string
}

func NewTransport() *Transport {
	t := &Transport{}
	u := url.URL{Scheme: "ws", Host: "localhost:" + config.Port, Path: "/things"}
	t.url = u.String()
	return t
}

func (t *Transport) dial() error {
	var err error = nil
	t.ws, _, err = websocket.DefaultDialer.Dial(t.url, nil)
	if err != nil {
		fmt.Printf("dial err: %s \r\n", err.Error())
		return err
	}
	return nil
}
