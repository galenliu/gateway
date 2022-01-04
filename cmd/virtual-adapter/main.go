package main

import (
	"github.com/galenliu/gateway/cmd/virtual-adapter/pkg"
	"github.com/galenliu/gateway/pkg/addon"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	manager := addon.NewAddonManager("virtual-adapter-golang")
	adapter := pkg.NewVirtualAdapter("virtual-adapter", "virtual-adapter")
	manager.AddAdapters(adapter)

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
