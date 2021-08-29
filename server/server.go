package server

import (
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/controllers"
)

type EventBus interface {
	Subscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
	SubscribeAsync(topic string, fn interface{})
}

type Config struct {
	HttpAddr    string
	HttpsAddr   string
	AddonUrls   []string
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

func NewServe(config Config, addonManager controllers.Manager, store controllers.Storage, bus EventBus, log logging.Logger) *WebServe {
	sev := &WebServe{}
	sev.options = config
	sev.logger = log
	sev.bus = bus
	sev.Router = controllers.NewRouter(controllers.Config{
		HttpAddr:  sev.options.HttpAddr,
		HttpsAddr: sev.options.HttpsAddr,
		AddonUrls: config.AddonUrls,
	}, addonManager, store, log)
	bus.SubscribeAsync(constant.AddonManagerStarted, sev.Start)
	return sev
}

func (serve *WebServe) Start() error {
	go func() {
		_ = serve.Router.Start()
	}()
	serve.bus.Publish(constant.WebServerStarted)
	return nil
}

func (serve *WebServe) Stop() error {
	err := serve.Router.Stop()
	if err != nil {
		serve.logger.Error("Web sever stop err: %s", err.Error())
		return err
	}
	return nil
}
