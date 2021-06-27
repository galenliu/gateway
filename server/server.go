package server

import (
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/server/controllers"
	"github.com/galenliu/gateway/server/models"
)

type Options struct {
	HttpPort    int
	HttpsPort   int
	StaticDir   string
	TemplateDir string
	UploadDir   string
	LogDir      string
}

type WebServe struct {
	*controllers.Router
	logger  logging.Logger
	options Options
	things  *models.Things
	users   *models.Users
}

func NewWebServe(options Options, log logging.Logger) *WebServe {
	sev := WebServe{}
	sev.options = options
	sev.logger = log
	sev.things = models.NewThingsModel(log)
	sev.users = models.NewUsersModel(log)

	sev.Router = controllers.NewAPP(sev.things, log)
	return &sev
}
