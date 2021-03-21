package main

import (
	"flag"
	"fmt"
	core "gateway"
	"gateway/config"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"os"
	"os/signal"
	"syscall"
)

var (
	proFile     string
	showVersion bool
)

func init() {
	flag.StringVar(&proFile, "profile", "", "Profile directory")
	flag.BoolVar(&showVersion, "version", false, "version")
}

func main() {

	var err error
	c := make(chan os.Signal)

	// 首先解析命令行参数
	flag.Parse()

	//if version command then print version
	if showVersion {
		fmt.Print(util.Version)
		return
	}

	//init config
	err = config.InitRuntime(proFile)
	CheckError(err)

	//create core instance
	gw, err := core.NewGateway()
	CheckError(err)

	//handle signal
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	var systemCall = func() {
		systemCall := <-c
		log.Info("exited system call %v", systemCall)
		gw.Close()
		os.Exit(0)
	}
	go systemCall()

	//run core
	CheckError(gw.Start())

}

func CheckError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
