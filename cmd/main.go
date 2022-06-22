package main

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter"
	"github.com/galenliu/gateway/cmd/yeelight-adapter"
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"github.com/galenliu/gateway/pkg/log"

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
	yeeAdapter := yeelight_adapter.NewVirtualAdapter(yeelightAdapterId)
	virtualAdapter := virtual_adapter.NewVirtualAdapter(virtualAdapterId)

	manager.AddAdapters(yeeAdapter, virtualAdapter)

	go yeeAdapter.StartPairing(nil)
	go virtualAdapter.StartPairing(nil)

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

	func() {
		for {
			select {
			case s := <-interruptChannel:
				log.Infof("signal %s exit", s.String())
			case <-manager.Done:
				log.Infof("Done shutting")
				return
			}
		}
	}()
}
