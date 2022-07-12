package gateway

import (
	"context"
	"github.com/galenliu/gateway/addon"
	"github.com/galenliu/gateway/api"
	"github.com/galenliu/gateway/pkg/db"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"path"
	"time"
)

type Config struct {
	BaseDir          string
	RemoveBeforeOpen bool
	Verbosity        string
	AddonUrls        []string
	IPCPort          string
	RPCPort          string
	HttpAddr         string
	HttpsAddr        string
	LogRotateDays    int
}

type Gateway struct {
	addonManager *addon.Manager
	sever        *api.WebServe
	logger       log.Logger
}

func NewGateway(ctx context.Context, config Config) (*Gateway, error) {
	g := &Gateway{}
	u := &messages.PluginRegisterResponseJsonDataUserProfile{
		BaseDir:    config.BaseDir,
		DataDir:    path.Join(config.BaseDir, "data"),
		AddonsDir:  path.Join(config.BaseDir, "addons"),
		ConfigDir:  path.Join(config.BaseDir, "config"),
		MediaDir:   path.Join(config.BaseDir, "media"),
		LogDir:     path.Join(config.BaseDir, "log"),
		GatewayDir: config.BaseDir,
	}
	log.Infof("userprofile: %v ", util.JsonIndent(u))

	//检查Gateway运行需要的文件目录
	util.EnsureDir(u.BaseDir, u.DataDir, u.ConfigDir, u.AddonsDir, u.ConfigDir, u.MediaDir, u.LogDir)

	// 数据初始化
	storage, err := db.NewStorage(u.ConfigDir, db.Config{
		Reset: config.RemoveBeforeOpen,
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("database init.")

	//Addon manager init
	g.addonManager = addon.NewAddonsManager(ctx, addon.Config{
		UserProfile: u,
		AddonsDir:   u.AddonsDir,
		IPCPort:     config.IPCPort,
		RPCPort:     config.RPCPort,
	}, storage)
	log.Infof("addon manager init.")

	// Web service init
	g.sever = api.NewServe(ctx, api.Config{
		HttpAddr:    config.HttpAddr,
		HttpsAddr:   config.HttpsAddr,
		AddonUrls:   config.AddonUrls,
		StaticDir:   path.Join(config.BaseDir, "static"),
		TemplateDir: path.Join(config.BaseDir, "template"),
		UploadDir:   path.Join(config.BaseDir, "upload"),
		LogDir:      path.Join(config.BaseDir, "log"),
	}, g.addonManager, storage)
	return g, nil
}

func (g *Gateway) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	log.Info("shutting down, wait 5 second ...")
	defer cancel()
	err := g.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	time.Sleep(1 * time.Second)
	<-ctx.Done()
	return nil
}
