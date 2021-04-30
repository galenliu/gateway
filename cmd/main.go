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

type Runner interface {
	Start() error
	Stop()
}

func main() {

	var err error
	sig := make(chan os.Signal, 1)

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
	runner, err := NewGateway()
	CheckError(err)

	//handle signal
	signal.Notify(sig)
	var systemCall = func() {
		callMessage := <-sig
		log.Info("exited system call %v", callMessage)
		runner.Stop()
		os.Exit(0)
	}
	//start
	runner.Start()
	systemCall()
	log.Info("exited")
}

func CheckError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

type HomeGateway struct {
	Preferences *config.Preferences
	closeChan   chan struct{}
	Tasks       []Runner
}

func NewGateway() (gateway *HomeGateway, err error) {
	gateway = &HomeGateway{}
	gateway.Tasks = append(gateway.Tasks, plugin.NewAddonsManager())
	gateway.Tasks = append(gateway.Tasks, controllers.NewWebAPP())
	gateway.closeChan = make(chan struct{})
	//update the gateway preferences
	return gateway, err
}

func (gateway *HomeGateway) Start() {
	log.Info("gateway start.....")
	func() {
		for _, task := range gateway.Tasks {
			task := task
			go func() {
				err := task.Start()
				if err != nil {
					log.Error(err.Error())
				}
			}()

		}
	}()
}

func (gateway *HomeGateway) Stop() {
	gateway.closeChan <- struct{}{}
	go func() {
		for _, task := range gateway.Tasks {
			task.Stop()
		}
	}()
}
