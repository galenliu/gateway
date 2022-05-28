package server

import "github.com/galenliu/gateway/pkg/matter/transport"

type GroupDataProviderListener struct {
	mTransports transport.TransportManager
}

func IntGroupDataProviderListener(transport transport.TransportManager) (*GroupDataProviderListener, error) {
	ins := &GroupDataProviderListener{mTransports: transport}
	return ins, nil
}
