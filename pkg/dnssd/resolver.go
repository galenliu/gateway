package dnssd

import (
	"github.com/galenliu/gateway/pkg/inet"
	"sync"
)

type Resolver struct {
}

func (r Resolver) Init(manager inet.UDPEndpointManager) {

}

var insResolver *Resolver
var onceResolver sync.Once

func ResolverInstance() *Resolver {
	onceResolver.Do(func() {
		insResolver = newResolver()
	})
	return insResolver
}

func newResolver() *Resolver {
	return &Resolver{}
}
