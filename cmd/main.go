package main

import (
	"flag"
	"fmt"
	gateway "gateway"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgDir        string
	logRotateDays int
	debug         bool

	showVersion bool
	reset       bool
)

func init() {
	flag.StringVar(&cfgDir, "config", gateway.GetDefaultConfigDir(), "Profile directory")
	flag.IntVar(&logRotateDays, "log_rotate_days", 0, "Enables daily log rotation and keeps up to the specified days")
	flag.BoolVar(&reset, "reset", false, "reset")
	flag.BoolVar(&showVersion, "version", false, "version")
	flag.BoolVar(&debug, "debug", true, "enable debug log to file")
}

func main() {

	c := make(chan os.Signal)

	// 首先解析命令行参数
	flag.Parse()

	//if version command then print version
	if showVersion {
		fmt.Print(gateway.Version)
		return
	}

	runtimeConfig := gateway.NewRuntimeConfig(cfgDir, logRotateDays, debug, reset)

	gw, err := gateway.CreateGateway(runtimeConfig)
	if err != nil {
		return
	}

	signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	var systemCall = func() {
		<-c
		gw.Close()
		os.Exit(0)
	}

	go systemCall()

	CheckError(gw.Run())

}

func CheckError(err error) {
	if err != nil {
		os.Exit(1)
	}

}
