package plugin

type Outlet struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (o *Outlet) getId() string {
	return o.Id
}
