package controllers

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/models/model"
	"github.com/gofiber/websocket/v2"
)

type wsClint struct {
	ws          *websocket.Conn
	thingsModel model.Container
	thingId     string
	readChan    chan []byte
	logger      logging.Logger
}

func NewWsClint(ws *websocket.Conn, thingId string, thing model.Container, log logging.Logger) *wsClint {
	c := &wsClint{}
	c.ws = ws
	c.thingsModel = thing
	c.logger = log
	c.readChan = make(chan []byte, 10)
	c.thingId = thingId
	go c.readLoop()
	return c
}

func (c *wsClint) handler() error {
	for {
		select {
		case data := <-c.readChan:
			c.handleMessage(data)
		}
	}
}

func (c wsClint) handleMessage(data []byte) {

}

func (c *wsClint) readLoop() {
	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		c.readChan <- data
	}
}
