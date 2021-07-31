package gateway

import (
	"context"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/database"
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
	DataDir            string
	AddonDirs          []string
	DBRemoveBeforeOpen bool
	Verbosity          string
	AddonUrls          []string
	IpcAddr            string
	HttpAddr           string
	HttpsAddr          string
	LogRotateDays      int

	HomeKitPin    string
	HomeKitEnable bool
}

type Gateway struct {
	options Options
	store   db.Store
	bus     eventBus
	logger       logging.Logger
	addonManager *plugin.Manager
	sever        *server.WebServe
}

func NewGateway(o Options, logger logging.Logger) (*Gateway, error) {
	g := &Gateway{}
	g.logger = logger
	g.options = o

	var e error = nil
	g.store, e = db.NewStore(path.Join(g.options.DataDir, constant.ConfigDirName), g.options.DBRemoveBeforeOpen)
	if e != nil {
		return nil, e
	}

	g.bus, e = bus.NewEventBus(g.logger)
	if e != nil {
		return nil, e
	}

	g.addonManager = plugin.NewAddonsManager(plugin.Options{
		DataDir:   g.options.DataDir,
		AddonDirs: g.options.AddonDirs,
	}, g.bus, g.logger)

	g.sever = server.Setup(server.Options{
		HttpAddr:    g.options.HttpAddr,
		HttpsAddr:   g.options.HttpsAddr,
		StaticDir:   path.Join(g.options.DataDir, "static"),
		TemplateDir: path.Join(g.options.DataDir, "template"),
		UploadDir:   path.Join(g.options.DataDir, "upload"),
		LogDir:      path.Join(g.options.DataDir, "log"),
	},g.addonManager, g.store, g.bus, g.logger)

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
