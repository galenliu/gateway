package internal

type PropertyHandler interface {
	HandleGetValue() (interface{}, error)
}

type GetValueFunc func() (interface{}, error)

type Property struct {
	Name   string `json:"name"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	AtType string `json:"@type"`

	getValueFunc GetValueFunc
}

func NewProperty(des string, getFunc GetValueFunc) *Property {
	prop := NewPropertyFromString(des)
	prop.getValueFunc = getFunc
	return prop
}

func NewPropertyFromString(des string) *Property {

	return nil
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetValue() (interface{}, error) {
	return p.getValueFunc()
}
