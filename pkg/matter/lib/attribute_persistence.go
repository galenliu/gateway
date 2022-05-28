package lib

type AttributePersistenceProvider interface {
	Init(storage PersistentStorageDelegate) error
}

type AttributePersistence struct {
	mStorage PersistentStorageDelegate
}

func (p AttributePersistence) Init(storage PersistentStorageDelegate) (err error) {
	p.mStorage = storage
	return
}
