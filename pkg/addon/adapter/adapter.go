package adapter

type AdapterProxy interface {
	GetId() string
	GetName() string
}

type Adapter struct {
	Id   string
	Name string
}

func (a Adapter) GetId() string {
	return a.Id
}

func (a Adapter) GetName() string {
	return a.Name
}
