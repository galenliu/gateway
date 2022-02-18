package api

import (
	"context"
	"github.com/galenliu/gateway/api/controllers"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin"
)

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
}

func NewServe(ctx context.Context, config Config, addonManager *plugin.Manager, store controllers.Storage, log logging.Logger) *WebServe {
	sev := &WebServe{}
	sev.options = config
	sev.logger = log
	sev.Router = controllers.NewRouter(config.AddonUrls, addonManager, store, log)
	sev.Start(ctx)
	return sev
}

func (app *WebServe) Start(ctx context.Context) {
	go func() {
		c, cancelFunc := context.WithCancel(ctx)
		select {
		case <-c.Done():
			cancelFunc()
			_ = app.Shutdown()
		default:
			err := app.Listen(app.options.HttpAddr)
			if err != nil {
				app.logger.Errorf("http api err:%s", err.Error())
				cancelFunc()
				return
			}
		}
		cancelFunc()
	}()

	go func() {
		c, cancelFunc := context.WithCancel(ctx)
		select {
		case <-c.Done():
			cancelFunc()
			_ = app.Shutdown()
		default:
			err := app.Listen(app.options.HttpsAddr)
			if err != nil {
				app.logger.Errorf("https api err:%s", err.Error())
				cancelFunc()
				return
			}
		}
		cancelFunc()
	}()
	app.logger.Infof("api running at adders: %v, %v \n", app.options.HttpAddr, app.options.HttpsAddr)
}
