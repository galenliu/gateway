package gateway

import (
	"context"
	"fmt"
	"gateway/addons"
	"gateway/app"
	"gateway/config"
	"gateway/event"
	"gateway/pkg/log"
	"time"
)

//gateway strut
type HomeGateway struct {
	Preferences   *config.Preferences
	EventsBus     *event.EventBus
	AddonsManager *addons.AddonManager
	Web           *app.WebApp
	Ctx           context.Context
}

func NewGateway() (gateway *HomeGateway, err error) {


	gateway = &HomeGateway{}
	gateway.Ctx = context.Background()

	//update the gateway preferences
	return gateway, err
}


func (gateway *HomeGateway) Start() error {
	event.InitEventBus()
	log.Info("gateway start.....")
	go gateway.Web.Start()
	go gateway.AddonsManager.Start()
	for {
		select {
		case <-gateway.Ctx.Done():
			return fmt.Errorf("application exit")
		default:
			time.Sleep(2 * time.Second)
		}
	}
}


func (gateway *HomeGateway) Close() {
	gateway.Ctx.Done()
}

