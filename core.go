package gateway

import (
	"context"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/database"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/wot"
	"github.com/galenliu/gateway/wot/models"
	"path"
	"time"
)

type Program interface {
	Start() error
	Stop() error
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
	options         Options
	store           database.Store
	eventBus        bus.EventBusController
	logger          logging.Logger
	addonManager    plugin.AddonManager
	thingsContainer wot.ThingsContainer
}

func NewGateway(o Options, logger logging.Logger) (*Gateway, error) {
	g := &Gateway{}
	g.logger = logger
	g.options = o

	var e error = nil
	g.store, e = database.NewStore(path.Join(g.options.DataDir, util.ConfigDirName), g.options.DBRemoveBeforeOpen)
	if e != nil {
		return nil, e
	}

	g.eventBus, e = bus.NewEventBus(g.logger)
	if e != nil {
		return nil, e
	}

	return g, nil
}

func (g *Gateway) Start() error {

	// 首先启动plugin
	g.addonManager = plugin.NewAddonsManager(plugin.Options{
		DataDir:   g.options.DataDir,
		AddonDirs: g.options.AddonDirs,
	}, g.logger)
	err := g.addonManager.Start()
	if err != nil {
		return err
	}

	g.thingsContainer = wot.NewThingsContainer(wot.Options{}, g.store, g.eventBus, g.logger)

	return nil
}

func (g *Gateway) FindNewThings() (ts []*models.Thing) {
	storedThings := g.thingsContainer.GetThings()
	connectedDevices := g.addonManager.GetDevices()
	for _, dev := range connectedDevices {
		var isExit = false
		for _, th := range storedThings {
			if dev.GetID() == th.GetID() {
				isExit = true
			}
		}
		if !isExit {
			t, err := models.NewThingFromString(dev.ToJson())
			if err != nil {
				ts = append(ts, t)
			}
		}
	}
	return
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
	return nil
}
