package gateway

import (
	"context"
	"gateway/addons"
	"gateway/util/database"
	"gateway/util/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Runtime       *RuntimeConfig
	Preferences   *Preferences
	AddonsManager *addons.Manager
	ctx           context.Context
}

type Event struct {
	EventType string
	thingsID  string
}

func (gateway *HomeGateway) Start() error {
	//gateway.AddonsManager.Start()
	return nil
}

func (gateway *HomeGateway) updatePreferences() {

	//open database and create table
	db := database.GetDB()
	_ = db.AutoMigrate(&Preferences{})
	_ = db.AutoMigrate(&Units{})

	var p Preferences
	result := db.First(&p)
	if result.Error != nil {
		u1 := Units{Temperature: PrefUnitsTempCelsius}
		p1 := Preferences{Language: PrefLangCn}
		p1.Units = u1
		db.Debug().Create(&p1)
	}
	var p2 Preferences
	db.Preload(clause.Associations).Debug().First(&p2)
	gateway.Preferences = &p2

}

func (gateway *HomeGateway) addonManagerLoadAndRun() error {

	addonManagerConfig := addons.ManagerConfig{
		AddonsDir: gateway.Runtime.AddonsDir,
	}
	addonManager := addons.NewAddonsManager(addonManagerConfig)
	gateway.AddonsManager = addonManager
	addonManager.LoadAddons()
	return nil
}

func (gateway *HomeGateway) Close() {

}
