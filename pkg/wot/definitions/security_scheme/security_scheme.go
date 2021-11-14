package security_scheme

import "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type SecurityScheme interface {
}

type securityScheme struct {
	AtType       string            `json:"@type"`
	Description  string            `json:"description"`
	Descriptions map[string]string `json:"descriptions"`
	Proxy        hypermedia_controls.URI
	Scheme       interface{}
}
