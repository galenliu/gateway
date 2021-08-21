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

type Storage interface {
	models.UsersStore
	models.ThingsStorage
	models.SettingsStore
	models.JsonwebtokenStore
}

type Config struct {
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
	options Config
	bus     EventBus
}

func Setup(config Config, addonManager controllers.AddonManagerHandler, store Storage, bus EventBus, log logging.Logger) *WebServe {
	sev := WebServe{}
	sev.options = config
	sev.logger = log
	sev.bus = bus

	sev.Router = controllers.Setup(controllers.Config{
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
