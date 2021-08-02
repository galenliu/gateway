package internal

type PropertyHandler interface {
	HandleGetValue() (interface{}, error)
}

type Property struct {
	Name   string `json:"name"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	AtType string `json:"@type"`

	handler PropertyHandler
}

func NewPropertyFormString(des string) *Property {

	return nil
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetValue() (interface{}, error) {
	return p.handler.HandleGetValue()
}
