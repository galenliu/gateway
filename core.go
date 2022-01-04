package gateway

import (
	"context"
	"github.com/galenliu/gateway/api"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/db"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	json "github.com/json-iterator/go"
	"path"
	"time"
)

type Config struct {
	BaseDir          string
	AttachAddonsDir  string
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
	bus          *bus.Bus
	addonManager *plugin.Manager
	sever        *api.WebServe
	logger       logging.Logger
}

func NewGateway(ctx context.Context, config Config, logger logging.Logger) (*Gateway, error) {

	g := &Gateway{}
	g.logger = logger
	g.config = config
	u := &messages.PluginRegisterResponseJsonDataUserProfile{
		BaseDir:    g.config.BaseDir,
		DataDir:    path.Join(g.config.BaseDir, "data"),
		AddonsDir:  path.Join(g.config.BaseDir, "addons"),
		ConfigDir:  path.Join(g.config.BaseDir, "config"),
		MediaDir:   path.Join(g.config.BaseDir, "media"),
		LogDir:     path.Join(g.config.BaseDir, "log"),
		GatewayDir: g.config.BaseDir,
	}
	s, _ := json.MarshalIndent(u, "", "   ")
	logger.Infof("userprofile: %v ", string(s))

	//检查Gateway运行需要的文件目录
	util.EnsureDir(logger, u.BaseDir, u.DataDir, u.ConfigDir, u.AddonsDir, u.ConfigDir, u.MediaDir, u.LogDir)

	// 数据初始化
	storage, err := db.NewStorage(u.ConfigDir, logger, db.Config{
		Reset: config.RemoveBeforeOpen,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Infof("database init.")

	//  EventBus init
	newBus := bus.NewBusController(logger)
	logger.Infof("events bus init.")

	//Addon manager init
	g.addonManager = plugin.NewAddonsManager(ctx, plugin.Config{
		UserProfile:     u,
		AddonsDir:       u.AddonsDir,
		AttachAddonsDir: g.config.AttachAddonsDir,
		IPCPort:         config.IPCPort,
		RPCPort:         config.RPCPort,
	}, storage, newBus, logger)
	logger.Infof("addon manager init.")

	// Web service init
	g.sever = api.NewServe(ctx, api.Config{
		HttpAddr:    g.config.HttpAddr,
		HttpsAddr:   g.config.HttpsAddr,
		AddonUrls:   g.config.AddonUrls,
		StaticDir:   path.Join(g.config.BaseDir, "static"),
		TemplateDir: path.Join(g.config.BaseDir, "template"),
		UploadDir:   path.Join(g.config.BaseDir, "upload"),
		LogDir:      path.Join(g.config.BaseDir, "log"),
	}, g.addonManager, storage, newBus, logger)
	g.bus = newBus
	logger.Infof("web api running.")
	return g, nil
}

func (g *Gateway) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	g.logger.Info("shutting down, wait 5 second ...")
	defer cancel()
	err := g.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	go g.bus.Publish(constant.GatewayStop)
	time.Sleep(1 * time.Second)
	<-ctx.Done()
	return nil
}
