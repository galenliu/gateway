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

type eventBus interface {
	Subscribe(topic string, fn interface{})
	Unsubscribe(topic string, fn interface{})
	Publish(topic string, args ...interface{})
	SubscribeOnce(topic string, fn interface{})
	SubscribeAsync(topic string, fn interface{})
}

type Options struct {
	BaseDir   string
	AddonDirs []string

	DBRemoveBeforeOpen bool
	Verbosity          string
	AddonUrls          []string
	IPCPort            string
	RPCPort            string
	HttpAddr           string
	HttpsAddr          string
	LogRotateDays      int
	HomeKitPin         string
	HomeKitEnable      bool
}

type Gateway struct {
	options      Options
	store        *db.Store
	bus          eventBus
	logger       logging.Logger
	addonManager *plugin.Manager
	sever        *server.WebServe
}

func NewGateway(o Options, logger logging.Logger) (*Gateway, error) {
	g := &Gateway{}
	g.logger = logger
	g.options = o

	var e error = nil
	g.store, e = db.NewStore(path.Join(g.options.BaseDir, constant.ConfigDirName), g.options.DBRemoveBeforeOpen, logger)
	if e != nil {
		return nil, e
	}

	g.bus, e = bus.NewEventBus(g.logger)
	if e != nil {
		return nil, e
	}

	g.addonManager = plugin.NewAddonsManager(plugin.Config{
		UserProfile: plugin.UserProfile{
			BaseDir:        g.options.BaseDir,
			DataDir:        path.Join(g.options.BaseDir, "data"),
			AddonsDir:      path.Join(g.options.BaseDir, "addons"),
			ConfigDir:      path.Join(g.options.BaseDir, "config"),
			UploadDir:      path.Join(g.options.BaseDir, "upload"),
			MediaDir:       path.Join(g.options.BaseDir, "media"),
			LogDir:         path.Join(g.options.BaseDir, "log"),
			GatewayVersion: Version,
		},
		AddonDirs: g.options.AddonDirs,
		IPCPort:   o.IPCPort,
		RPCPort:   o.IPCPort,
	}, g.store, g.bus, g.logger)

	g.sever = server.Setup(server.Options{
		HttpAddr:    g.options.HttpAddr,
		HttpsAddr:   g.options.HttpsAddr,
		StaticDir:   path.Join(g.options.BaseDir, "static"),
		TemplateDir: path.Join(g.options.BaseDir, "template"),
		UploadDir:   path.Join(g.options.BaseDir, "upload"),
		LogDir:      path.Join(g.options.BaseDir, "log"),
	}, g.addonManager, g.store, g.bus, g.logger)

	return g, nil
}

//func (g *Gateway) FindNewThings() (ts []*models.Thing) {
//	storedThings := g.thingsContainer.GetThings()
//	connectedDevices := g.addonManager.GetDevices()
//	for _, dev := range connectedDevices {
//		var isExit = false
//		for _, th := range storedThings {
//			if dev.GetID() == th.GetID() {
//				isExit = true
//			}
//		}
//		if !isExit {
//			t, err := models.NewThingFromString(dev.ToJson())
//			if err != nil {
//				ts = append(ts, t)
//			}
//		}
//	}
//	return
//}

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
