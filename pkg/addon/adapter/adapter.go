package adapter

type AdapterProxy interface {
	GetId() string
}

type Adapter struct {
	Id string
}

func (a Adapter) GetId() string {
	return a.Id
}
