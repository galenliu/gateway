package gateway

import (
	"context"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/container"
	"github.com/galenliu/gateway/pkg/db"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server"
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
	storage      *db.Storage
	bus          bus.Controller
	logger       logging.Logger
	addonManager *plugin.Manager
	sever        *server.WebServe
	container    container.Container
}

func NewGateway(ctx context.Context,config Config, logger logging.Logger) (*Gateway, error) {

	g := &Gateway{}
	g.logger = logger
	g.config = config
	u := &rpc.UsrProfile{
		BaseDir:    g.config.BaseDir,
		DataDir:    path.Join(g.config.BaseDir, "data"),
		AddonsDir:  path.Join(g.config.BaseDir, "addons"),
		ConfigDir:  path.Join(g.config.BaseDir, "config"),
		MediaDir:   path.Join(g.config.BaseDir, "media"),
		LogDir:     path.Join(g.config.BaseDir, "log"),
		GatewayDir: g.config.BaseDir,
	}
	s, _ := json.MarshalIndent(u,"","   ")
	g.logger.Infof("userprofile: %v ", string(s))

	//检查Gateway运行需要的文件目录
	util.EnsureDir(logger, u.BaseDir, u.DataDir, u.ConfigDir, u.AddonsDir, u.ConfigDir, u.MediaDir, u.LogDir)

	// 数据化初始化
	var err error
	g.storage, err = db.NewStorage(u.ConfigDir, logger, db.Config{
		Reset: config.RemoveBeforeOpen,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	g.logger.Infof("database init.")

	//  Things container init
	g.container = container.NewThingsContainerModel(g.storage, g.logger)
	g.logger.Infof("things container init.")

	//  EventBus init
	newBus := bus.NewBusController(g.logger)
	g.logger.Infof("event bus init.")

	//Addon manager init
	g.addonManager = plugin.NewAddonsManager(ctx,plugin.Config{
		UserProfile: u,
		AddonsDir:       u.AddonsDir,
		AttachAddonsDir: g.config.AttachAddonsDir,
		IPCPort:         config.IPCPort,
		RPCPort:         config.RPCPort,
	}, g.storage, newBus, g.logger)
	g.logger.Infof("addon manager init.")

	// Web service init
	g.sever = server.NewServe(ctx,server.Config{
		HttpAddr:    g.config.HttpAddr,
		HttpsAddr:   g.config.HttpsAddr,
		AddonUrls:   g.config.AddonUrls,
		StaticDir:   path.Join(g.config.BaseDir, "static"),
		TemplateDir: path.Join(g.config.BaseDir, "template"),
		UploadDir:   path.Join(g.config.BaseDir, "upload"),
		LogDir:      path.Join(g.config.BaseDir, "log"),
	}, g.addonManager, g.addonManager, g.container, g.storage, newBus, g.logger)
	g.bus = newBus
	g.logger.Infof("gateway web server running.")
	return g, nil
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
