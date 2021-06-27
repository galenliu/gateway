package models

type PIN struct {
	Required bool        `json:"required,omitempty"`
	Pattern  interface{} `json:"pattern,omitempty"`
}
