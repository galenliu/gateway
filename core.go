package gateway

import (
	"context"
	"fmt"
	"gateway/addons"
	"gateway/app"
	"gateway/pkg/database"
	"gateway/pkg/log"
	"gateway/pkg/util"
	"gateway/config"
	"gorm.io/gorm/clause"
	"time"
)





//gateway strut
type HomeGateway struct {
	Preferences   *config.Preferences
	EventsBus     *EventsBus
	AddonsManager *addons.AddonsManager
	Web           *app.WebApp
	Ctx           context.Context
}

func NewGateway() (gateway *HomeGateway, err error) {


	gateway = &HomeGateway{}
	gateway.EventsBus = NewEventBus()
	gateway.Ctx = context.Background()

	//update the gateway preferences
	gateway.updatePreferences()
	return gateway, err
}

type Event struct {
	EventType string
	thingsID  string
}

func (gateway *HomeGateway) Start() error {

	log.Info("gateway start")
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

func (gateway *HomeGateway) updatePreferences() {

	//open database and create table
	db := database.GetDB()
	_ = db.AutoMigrate(&config.Preferences{})
	_ = db.AutoMigrate(&config.Units{})

	var p config.Preferences
	result := db.First(&p)
	if result.Error != nil {
		u1 := config.Units{Temperature: util.PrefUnitsTempCelsius}
		p1 := config.Preferences{Language: util.PrefLangCn}
		p1.Units = u1
		db.Debug().Create(&p1)
	}
	var p2 config.Preferences
	db.Preload(clause.Associations).Debug().First(&p2)
	gateway.Preferences = &p2

}

func (gateway *HomeGateway) Close() {
	gateway.Ctx.Done()
}

