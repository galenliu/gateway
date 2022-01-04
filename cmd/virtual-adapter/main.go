package main

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/virtual-adapter/pkg"
	"github.com/galenliu/gateway/pkg/addon"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	manager, err := addon.NewAddonManager("virtual-adapter-golang")
	if err != nil {
		fmt.Printf("addon manager error: %s", err.Error())
		return
	}
	adapter := pkg.NewVirtualAdapter(manager, "virtual-adapter", "virtual-adapter")
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
