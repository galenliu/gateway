package main

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/virtual"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var packageId = "virtual-adapter-go"
var virtualAdapterId = "virtual-adapter"
var yeelightAdapterId = "yeelight-adapter"

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	manager, err := proxy.NewAddonManager(ctx, packageId)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}
	yeeAdapter := yeelight.NewVirtualAdapter(yeelightAdapterId)
	virtualAdapter := virtual.NewVirtualAdapter(virtualAdapterId)

	manager.AddAdapters(yeeAdapter, virtualAdapter)

	yeeAdapter.StartPairing(nil)
	virtualAdapter.StartPairing(nil)

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	func() {
		for {
			select {
			case s := <-interruptChannel:
				log.Printf("signal %s exit", s.String())
			case <-manager.Done:
				return
			}
		}
	}()
}
