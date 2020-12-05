package main

import (
	"flag"
	"fmt"
	core "gateway"
	"gateway/addons"
	"gateway/app"
	"gateway/pkg/runtime"
	"gateway/pkg/util"
	"os"
	"os/signal"
	"syscall"
)

var (
	proFile  string
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

	//init runtime
	err = runtime.InitRuntime(proFile)
	CheckError(err)

	//create core instance
	gw, err := core.NewGateway()
	CheckError(err)

	gw.AddonsManager,err=addons.NewAddonsManager(gw.Ctx)
	CheckError(err)


	gw.Web =app.NewWebAPP(app.NewDefaultWebConfig())



	//handle signal
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	var systemCall = func() {
		<-c
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
