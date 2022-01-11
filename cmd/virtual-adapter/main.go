package main

import (
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	adapter := yeelight.NewVirtualAdapter("virtual-adapter-golang", "virtual-adapter", "virtual-adapter")
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case _ = <-interruptChannel:
				adapter.CloseProxy()
			}
		}
	}()

	for {
		if adapter.ProxyRunning() {
			time.Sleep(2 * time.Second)
		} else {
			return
		}
	}
}
