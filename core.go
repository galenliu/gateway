package gateway

import (
	"context"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/db"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server"
	"path"
	"time"
)

type Component interface {
	Start() error
	Stop() error
}

type BusController interface {
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
	bus          BusController
	logger       logging.Logger
	addonManager *plugin.Manager
	sever        *server.WebServe
}

func NewGateway(config Config, logger logging.Logger) (*Gateway, error) {

	g := &Gateway{}
	g.logger = logger
	g.config = config
	u := &rpc.PluginRegisterResponseMessage_Data_UsrProfile{
		BaseDir:    g.config.BaseDir,
		DataDir:    path.Join(g.config.BaseDir, "data"),
		AddonsDir:  path.Join(g.config.BaseDir, "addons"),
		ConfigDir:  path.Join(g.config.BaseDir, "config"),
		MediaDir:   path.Join(g.config.BaseDir, "media"),
		LogDir:     path.Join(g.config.BaseDir, "log"),
		GatewayDir: "",
	}

	//检查Gateway运行需要的文件目录
	err := util.EnsureDir(u.BaseDir, u.DataDir, u.ConfigDir, u.AddonsDir, u.ConfigDir, u.MediaDir, u.LogDir)

	// 数据化初始化
	g.storage, err = db.NewStorage(u.ConfigDir, logger, db.Config{
		Reset: config.RemoveBeforeOpen,
	})
	if err != nil {
		return nil, err
	}

	//  EventBus init
	newBus, err := bus.NewBus(g.logger)
	if err != nil {
		return nil, err
	}

	//Addon manager init
	g.addonManager = plugin.NewAddonsManager(plugin.Config{
		UserProfile: u,
		Preferences: &rpc.PluginRegisterResponseMessage_Data_Preferences{
			Language: "zh-cn",
			Units:    &rpc.PluginRegisterResponseMessage_Data_Preferences_Units{Temperature: ""},
		},
		AddonDirs: g.config.AddonDirs,
		IPCPort:   config.IPCPort,
		RPCPort:   config.RPCPort,
	}, g.storage, newBus, g.logger)

	// Web service init
	g.sever = server.NewServe(server.Config{
		HttpAddr:    g.config.HttpAddr,
		HttpsAddr:   g.config.HttpsAddr,
		StaticDir:   path.Join(g.config.BaseDir, "static"),
		TemplateDir: path.Join(g.config.BaseDir, "template"),
		UploadDir:   path.Join(g.config.BaseDir, "upload"),
		LogDir:      path.Join(g.config.BaseDir, "log"),
	}, g.addonManager, g.storage, newBus, g.logger)

	g.bus = newBus
	return g, nil
}

func (g *Gateway) Start() error {
	// 向总线发送启运信号
	go g.bus.Publish(constant.GatewayStart)
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
	go g.bus.Publish(constant.GatewayStop)
	time.Sleep(1 * time.Second)
	<-ctx.Done()
	return nil
}
