package controllers

import (
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/gofiber/websocket/v2"
)

type wsContainer struct {
	things map[string]*things.Thing
}

func (c *wsContainer) addThing() {

}

func handleWebsocket(model things.Container, log logging.Logger) func(conn *websocket.Conn) {
	handler := func(c *websocket.Conn) {
		defer c.Close()
		thingId := c.Params("thingId")
		log.Infof("websocket connection")
		sendChannel := newSendChanel(c)
		container := wsContainer{}

		if thingId == "" {
			for _, thing := range model.GetThings() {
				container.addThing(thing)
			}
		}

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Infof("read:", err)
				break
			}
			log.Infof("rev: %s", msg)
		}

	}
	return handler
}

func newSendChanel(conn *websocket.Conn) chan<- any {
	sendChan := make(chan any, 10)
	go func() {
		for {
			select {

			case data, ok := <-sendChan:
				if ok {
					err := conn.WriteJSON(data)
					if err != nil {
						return
					}
				}
			}

		}
	}()
	return sendChan
}
