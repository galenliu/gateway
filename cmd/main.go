package main

import (
	"flag"
	"fmt"
	"github.com/galenliu/gateway/config"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/controllers"
	"os"
	"os/signal"
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
	conf := config.NewConfig(proFile)
	if conf == nil {
		log.Info("config is bad")
		return
	}

	//create core instance
	gw, err := NewGateway()
	CheckError(err)

	//handle signal
	signal.Notify(c)
	var systemCall = func() {
		systemCall := <-c
		log.Info("exited system call %v", systemCall)
		gw.Close()
		os.Exit(0)
	}
	//start
	gw.Start()
	systemCall()

}

func CheckError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

type HomeGateway struct {
	Preferences   *config.Preferences
	AddonsManager *plugin.AddonManager
	Web           *controllers.Web
	closeChan     chan struct{}
}

func NewGateway() (gateway *HomeGateway, err error) {

	gateway = &HomeGateway{}
	gateway.AddonsManager = plugin.NewAddonsManager()
	gateway.Web = controllers.NewWebAPP()
	gateway.closeChan = make(chan struct{})
	//update the gateway preferences
	return gateway, err
}

func (gateway *HomeGateway) Start() {
	log.Info("gateway start.....")
	go func() {
		err := gateway.Web.Start()
		if err != nil {

		}
	}()
	go gateway.AddonsManager.Start()
}

func (gateway *HomeGateway) Close() {
	gateway.closeChan <- struct{}{}
	gateway.AddonsManager.Close()
	err := gateway.Web.Close()
	if err != nil {
		log.Error(err.Error())
	}
}
