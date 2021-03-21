package gateway

import (
	"fmt"
	"gateway/config"
	"gateway/pkg/log"
	"gateway/plugin"
	"gateway/server/controllers"
	"time"
)

//gateway strut
type HomeGateway struct {
	Preferences   *config.Preferences
	AddonsManager *plugin.AddonManager
	Web           *controllers.Web
	closeChan     chan struct{}
}

func NewGateway() (gateway *HomeGateway, err error) {

	gateway = &HomeGateway{}
	gateway.AddonsManager = plugin.NewAddonsManager()
	gateway.Web = controllers.NewWebAPP(controllers.NewDefaultWebConfig())
	gateway.closeChan = make(chan struct{})
	//update the gateway preferences
	return gateway, err
}

func (gateway *HomeGateway) Start() error {

	log.Info("gateway start.....")
	go gateway.Web.Start()
	go gateway.AddonsManager.Start()
	for {
		select {
		case <-gateway.closeChan:
			return fmt.Errorf("gateway exit")
		default:
			time.Sleep(2 * time.Second)
		}
	}
}

func (gateway *HomeGateway) Close() {
	gateway.closeChan <- struct{}{}
	gateway.AddonsManager.Close()
	err := gateway.Web.Close()
	if err != nil {
		log.Error(err.Error())
	}
}
