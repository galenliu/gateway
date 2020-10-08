package gateway

import (
	"gateway/db"
	"go.uber.org/zap"
)

func CreateGateway(rc *RuntimeConfig) (gateway *HomeGateway, err error) {

	//fist : ensure runtime,check the dir etc
	gateway = NewHomeGateway(rc.cfgDir)
	//set the logger database
	var logDir = gateway.UserProfile.LogDir
	var dbDir = gateway.UserProfile.ConfigDir

	//设置logger
	InitLogger(logDir, true, rc.logRotateDays)

	//init database
	if rc.reset {
		err = db.ResetDB(dbDir)
		if err != nil {
			Log.Error("remove database err", zap.Error(err))
			return
		}
	}
	err = db.InitDB(dbDir)
	if err != nil {
		Log.Error("open data base err", zap.Error(err))
		return nil, err
	}

	//update the gateway preferences
	err = gateway.updatePreferences()
	if err != nil {
		Log.Error("update preferences err", zap.Error(err))
		return
	}

	err = gateway.addonManagerLoadAndRun()
	if err != nil {
		Log.Error("addon manager load err", zap.Error(err))
		return nil, err
	}

	err = InitRouter("./static/", "../static/template/*", gateway.UserProfile.UploadDir, gateway.UserProfile.LogDir)
	if err != nil {
		return nil, err
	}

	return
}
