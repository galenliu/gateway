package main

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/virtual"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	manager, err := proxy.NewAddonManager("virtual-adapter-go")
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}
	yeeAdapter := yeelight.NewVirtualAdapter("yeelight-adapter")
	virtualAdapter := virtual.NewVirtualAdapter("virtual-adapter")
	manager.AddAdapters(yeeAdapter, virtualAdapter)

	yeeAdapter.StartPairing(1 * time.Duration(1))
	virtualAdapter.StartPairing(1 * time.Duration(1))
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case _ = <-interruptChannel:
				yeeAdapter.CloseProxy()
				virtualAdapter.CloseProxy()
			}
		}
	}()

	for {
		if yeeAdapter.ProxyRunning() {
			time.Sleep(2 * time.Second)
		} else {
			return
		}
	}
}
