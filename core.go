package gateway

import (
	"context"
	"fmt"
	"gateway/addons"
	"gateway/app"
	"gateway/pkg/database"
	"gateway/pkg/logger"
	"gateway/pkg/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var log logger.Logger

type Units struct {
	gorm.Model
	Temperature string `gorm:"default: degree_celsius"`
}

type Preferences struct {
	gorm.Model
	Language string `gorm:"default: zh-cn"`
	Units    Units
	UnitsID  int
}

//gateway strut
//
type HomeGateway struct {
	Preferences   *Preferences
	AddonsManager *addons.AddonsManager
	Web           *app.WebApp
	Ctx           context.Context
}

func NewGateway() (gateway *HomeGateway, err error) {

	log = logger.GetLog()
	gateway = &HomeGateway{}
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
	_ = db.AutoMigrate(&Preferences{})
	_ = db.AutoMigrate(&Units{})

	var p Preferences
	result := db.First(&p)
	if result.Error != nil {
		u1 := Units{Temperature: util.PrefUnitsTempCelsius}
		p1 := Preferences{Language: util.PrefLangCn}
		p1.Units = u1
		db.Debug().Create(&p1)
	}
	var p2 Preferences
	db.Preload(clause.Associations).Debug().First(&p2)
	gateway.Preferences = &p2

}


func (gateway *HomeGateway) Close() {
	  gateway.Ctx.Done()
}
