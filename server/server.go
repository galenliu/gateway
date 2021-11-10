package server

import (
	"context"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/container"
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
	logger    logging.Logger
	options   Config
	bus       bus.Controller
	container container.Container
}

func NewServe(ctx context.Context, config Config, addonManager controllers.Manager, container container.Container, store controllers.Storage, bus bus.Controller, log logging.Logger) *WebServe {
	sev := &WebServe{}
	sev.container = container
	sev.options = config
	sev.logger = log
	sev.bus = bus
	sev.Router = controllers.NewRouter(ctx, controllers.Config{
		HttpAddr:  sev.options.HttpAddr,
		HttpsAddr: sev.options.HttpsAddr,
		AddonUrls: config.AddonUrls,
	}, addonManager, container, store, bus, log)
	return sev
}
