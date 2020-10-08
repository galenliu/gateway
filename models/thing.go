package models

type Thing struct {
	Tittle     string `json:"tittle"`
	Properties map[string]interface{}
}

func (t *Thing) SetIcon(data string, updateDatabase bool) {

}
