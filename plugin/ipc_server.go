package plugin

import (
	"context"
	"fmt"
	"github.com/fasthttp/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsChan chan *websocket.Conn

func NewIpcServer(ctx context.Context, addr string) (chan *websocket.Conn, chan string) {

	wsChan = make(chan *websocket.Conn, 64)
	errChan := make(chan string)
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
	}()
	go func() {
		log.Println("listening at " + addr)
		err := srv.ListenAndServe()
		fmt.Println("waiting for the remaining connections to finish...")
		wg.Wait()
		if err != nil && err != http.ErrServerClosed {
			select {
			case errChan <- err.Error():
			}
		}
		log.Println("gracefully shutdown the http server...")
	}()
	return wsChan, errChan
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	select {
	case wsChan <- ws:
	}
}

func ReadLong(ws *websocket.Conn) chan []byte {
	return nil
}
