package main

import (
	"flag"
	"fmt"
	"github.com/galenliu/gateway/configs"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/controllers"
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

type Runner interface {
	Start() error
	Stop()
}

// TermFunc defines the function which is executed on termination.
type TermFunc func(sig os.Signal)

// OnTermination calls a function when the app receives an interrupt of kill signal.
func OnTermination(fn TermFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGTERM)

	func() {
		select {
		case sig := <-c:
			if fn != nil {
				fn(sig)
			}
		}
	}()
}

func main() {

	var err error

	// 首先解析命令行参数
	flag.Parse()

	//if version command then print version
	if showVersion {
		fmt.Print(util.Version)
		return
	}

	//init config
	conf := configs.NewConfig(proFile)
	if conf == nil {
		log.Info("config is bad")
		return
	}

	//create core instance
	runner, e := NewGateway()
	CheckError(e)

	//handle signal
	var systemCall = func(sig os.Signal) {
		log.Info("exited system call %v", sig.String())
		runner.Stop()
		os.Exit(0)
	}

	OnTermination(systemCall)

	//start
	err = runner.Start()
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

type HomeGateway struct {
	Preferences *configs.Preferences
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

func (gateway *HomeGateway) Start() error {
	log.Info("gateway start .....")
	for _, task := range gateway.Tasks {
		err := task.Start()
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

func (gateway *HomeGateway) Stop() {
	gateway.closeChan <- struct{}{}
	go func() {
		for _, task := range gateway.Tasks {
			task.Stop()
		}
	}()
}
