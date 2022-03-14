package main

import (
	"context"
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
	ctx, cancelFunc := context.WithCancel(context.Background())
	manager, err := proxy.NewAddonManager("virtual-adapter-go")
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}
	yeeAdapter := yeelight.NewVirtualAdapter("yeelight-adapter")
	virtualAdapter := virtual.NewVirtualAdapter("virtual-adapter")

	manager.AddAdapters(yeeAdapter, virtualAdapter)

	yeeAdapter.StartPairing(nil)
	virtualAdapter.StartPairing(nil)

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case _ = <-interruptChannel:
				cancelFunc()
				manager.Close()
			}
		}
	}()

	go func() {
		c, f := context.WithCancel(ctx)
		for {
			fmt.Printf("automatically discover deivces every 30 minute\n")
			yeeAdapter.StartPairing(nil)
			select {
			case <-c.Done():
				f()
				return
			default:
				time.Sleep(30 * time.Minute)
			}
		}
	}()

	for {
		if manager.IsRunning() {
			time.Sleep(5 * time.Second)
		} else {
			cancelFunc()
			return
		}
	}
}
