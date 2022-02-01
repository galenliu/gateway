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
	sev.Router = controllers.NewRouter(ctx, controllers.Config{
		HttpAddr:  sev.options.HttpAddr,
		HttpsAddr: sev.options.HttpsAddr,
		AddonUrls: config.AddonUrls,
	}, addonManager, store, log)
	return sev
}
