package server

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/server/controllers"
	"github.com/galenliu/gateway/server/models"
)

type eventBus interface {
	Subscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
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
	things  *models.Things
	users   *models.Users

	bus eventBus
}

func NewWebServe(options Options, bus eventBus, log logging.Logger) *WebServe {
	sev := WebServe{}
	sev.options = options
	sev.logger = log
	sev.bus = bus
	sev.things = models.NewThingsModel(log)
	sev.users = models.NewUsersModel(log)

	sev.Router = controllers.NewAPP(controllers.Options{
		HttpAddr:  sev.options.HttpAddr,
		HttpsAddr: sev.options.HttpsAddr,
	}, sev.things, log)

	sev.bus.Subscribe(util.GatewayStarted, sev.start)
	return &sev
}

func (serve *WebServe) start() error {

	err := serve.Router.Start()
	if err != nil {
		serve.logger.Error("Web server err: %s", err.Error())
		return err
	}
	serve.bus.Publish(util.WebServerStarted)
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
