package thing

type Property struct {
	Name        string `json:"name"`
	AtType      string `json:"@type"` //引用的类型
	Type        string `json:"type"`  //数据的格式
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Unit     string `json:"unit,omitempty"` //属性的单位
	ReadOnly bool   `json:"readOnly"`
	Visible  bool   `json:"visible"`

	Minimum interface{} `json:"minimum,omitempty"`
	Maximum interface{} `json:"maximum,omitempty"`
	Value   interface{} `json:"value"`
	Enum    []string    `json:"enum,omitempty"`

	Links []string `json:"links"`
	Href  string   `json:"href"`

	ThingId string `json:"-"`
}

func (property *Property) SetValue(value interface{}) (interface{}, error) {
	return nil, nil
}
