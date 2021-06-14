package gateway

import (
	"context"
	"github.com/galenliu/gateway/pkg/logging"
)

type Options struct {
	DataDir            string
	DBRemoveBeforeOpen bool
}

type Gateway struct {
}

func NewGateway(logger logging.Logger, o Options) (*Gateway, error) {
	g := &Gateway{}

	return g, nil
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	return nil
}
