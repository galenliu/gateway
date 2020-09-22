package zeroconf

import "gateway/config"

type Service interface {
	RegisterService(c *config.Config)
	Stop()
}
