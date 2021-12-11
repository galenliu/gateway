package security_scheme

import "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type ComboSecurityScheme struct {
	*securityScheme
	OnOf  hypermedia_controls.ArrayOrString `json:"onOf"`
	AllOf hypermedia_controls.ArrayOrString `json:"allOf"`
}
