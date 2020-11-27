package models


type Property struct {
	AtType      string `json:"@type"` //引用的类型
	Type        string `json:"type"`  //数据的格式
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`

	Unit     string `json:"unit,omitempty"` //属性的单位
	ReadOnly bool   `json:"read_only"`
	Visible  bool   `json:"visible"`


	Minimum interface{} `json:"minimum,omitempty,string"`
	Maximum interface{} `json:"maximum,omitempty,string"`
	Value   interface{}
	Enum    []interface{}
}