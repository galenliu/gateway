package gateway

import (
	"context"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/db"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server"
	"path"
	"time"
)

type Component interface {
	Start() error
	Stop() error
}

type EvBus interface {
	Subscribe(topic string, fn interface{})
	Unsubscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
	SubscribeOnce(topic string, fn interface{})
	SubscribeAsync(topic string, fn interface{})
}

type Config struct {
	BaseDir          string
	AddonDirs        []string
	RemoveBeforeOpen bool
	Verbosity        string
	AddonUrls        []string
	IPCPort          string
	RPCPort          string
	HttpAddr         string
	HttpsAddr        string
	LogRotateDays    int
	HomeKitPin       string
	HomeKitEnable    bool
}

type Gateway struct {
	config       Config
	storage      *db.Storage
	bus          EvBus
	logger       logging.Logger
	addonManager *plugin.Manager
	sever        *server.WebServe
}

func NewGateway(config Config, logger logging.Logger) (*Gateway, error) {
	g := &Gateway{}
	g.logger = logger
	g.config = config

	var e error = nil
	g.storage, e = db.NewStorage(path.Join(g.config.BaseDir, constant.ConfigDirName), logger, db.Config{
		Reset: config.RemoveBeforeOpen,
	})
	if e != nil {
		return nil, e
	}

	g.bus, e = bus.NewEventBus(g.logger)
	if e != nil {
		return nil, e
	}

	g.addonManager = plugin.NewAddonsManager(plugin.Config{
		UserProfile: plugin.UserProfile{
			BaseDir:        g.config.BaseDir,
			DataDir:        path.Join(g.config.BaseDir, "data"),
			AddonsDir:      path.Join(g.config.BaseDir, "addons"),
			ConfigDir:      path.Join(g.config.BaseDir, "config"),
			UploadDir:      path.Join(g.config.BaseDir, "upload"),
			MediaDir:       path.Join(g.config.BaseDir, "media"),
			LogDir:         path.Join(g.config.BaseDir, "log"),
			GatewayVersion: Version,
		},
		AddonDirs: g.config.AddonDirs,
		IPCPort:   config.IPCPort,
		RPCPort:   config.IPCPort,
	}, g.storage, g.bus, g.logger)

	g.sever = server.Setup(server.Config{
		HttpAddr:    g.config.HttpAddr,
		HttpsAddr:   g.config.HttpsAddr,
		StaticDir:   path.Join(g.config.BaseDir, "static"),
		TemplateDir: path.Join(g.config.BaseDir, "template"),
		UploadDir:   path.Join(g.config.BaseDir, "upload"),
		LogDir:      path.Join(g.config.BaseDir, "log"),
	}, g.addonManager, g.storage, g.bus, g.logger)

	return g, nil
}

func (g *Gateway) Start() error {
	// 首先启动plugin
	err := g.addonManager.Start()
	if err != nil {
		return err
	}

	g.bus.Publish(constant.GatewayStarted)
	return nil
}

func (g *Gateway) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := g.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	g.bus.Publish(constant.GatewayStopped)
	<-ctx.Done()
	return nil
}
