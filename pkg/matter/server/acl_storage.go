package server

import (
	"github.com/galenliu/gateway/pkg/matter/credentials"
	"github.com/galenliu/gateway/pkg/matter/lib"
)

type AclStorage struct {
}

func (s AclStorage) Init(storage lib.PersistentStorageDelegate, fabrics *credentials.FabricTable) error {
	return nil
}
