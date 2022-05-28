package credentials

import "github.com/galenliu/gateway/pkg/matter/lib"

type FabricTable struct {
}

func (f FabricTable) FabricCount() int {
	//TODO implement me
	panic("implement me")
}

func (f FabricTable) Init(storage lib.PersistentStorageDelegate) (err error) {
	return
}

func (f FabricTable) DeleteAllFabrics() {

}
