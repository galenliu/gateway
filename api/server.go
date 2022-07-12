package api

import (
	"context"
	"github.com/galenliu/gateway/addon"
	"github.com/galenliu/gateway/api/controllers"
	"github.com/galenliu/gateway/pkg/log"
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
	options Config
}

func NewServe(ctx context.Context, config Config, addonManager *addon.Manager, store controllers.Storage) *WebServe {
	sev := &WebServe{}
	sev.options = config
	sev.Router = controllers.NewRouter(config.AddonUrls, addonManager, store)
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
				log.Errorf("http api err:%s", err.Error())
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
				log.Errorf("https api err:%s", err.Error())
				cancelFunc()
				return
			}
		}
		cancelFunc()
	}()
	log.Infof("api running at adders: %v, %v \n", app.options.HttpAddr, app.options.HttpsAddr)
}
