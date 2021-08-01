package server

import (
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/controllers"
	"github.com/galenliu/gateway/server/models"
)

type EventBus interface {
	Subscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
}

type Store interface {
	models.UsersStore
	models.ThingsStore
	models.SettingsStore
	models.JsonwebtokenStore
}


type AddonManager interface {
	controllers.AddonHandler
	controllers.ThingsHandler
}

type Options struct {
	HttpAddr    string
	HttpsAddr   string
	StaticDir   string
	TemplateDir string
	UploadDir   string
	LogDir      string
}

type WebServe struct {
	*controllers.Router
	logger  logging.Logger
	options Options
	bus     EventBus
}

func Setup(options Options,addonManager AddonManager, store Store, bus EventBus, log logging.Logger) *WebServe {
	sev := WebServe{}
	sev.options = options
	sev.logger = log
	sev.bus = bus

	sev.Router = controllers.Setup(controllers.Options{
		HttpAddr:  sev.options.HttpAddr,
		HttpsAddr: sev.options.HttpsAddr,
	}, addonManager, store, log)
	sev.bus.Subscribe(constant.GatewayStarted, sev.start)
	return &sev
}

func (serve *WebServe) start() error {

	err := serve.Router.Start()
	if err != nil {
		serve.logger.Error("Web server err: %s", err.Error())
		return err
	}
	serve.bus.Publish(constant.WebServerStarted)
	return nil
}

func (serve *WebServe) stop() error {
	err := serve.Router.Stop()
	if err != nil {
		serve.logger.Error("Web sever stop err: %s", err.Error())
		return err
	}
	return nil
}
