package security_scheme

import "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type ComboSecurityScheme struct {
	*securityScheme
	OnOf  hypermedia_controls.ArrayOfString `json:"onOf"`
	AllOf hypermedia_controls.ArrayOfString `json:"allOf"`
}
