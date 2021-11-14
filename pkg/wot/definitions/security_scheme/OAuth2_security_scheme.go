package security_scheme

import controls "github.com/galenliu/gateway/pkg/wot/definitions/hypermedia_controls"

type OAuth2SecurityScheme struct {
	*securityScheme
	Authorization controls.URI `json:"authorization"`
	Token         controls.URI `json:"token"`
	Refresh       controls.URI `json:"refresh"`
	Scopes        controls.ArrayOfString
	Flow          string `json:"flow"`
}

type oAuth2 struct {
	scheme        string
	flow          string
	authorization string
	token         string
	scopes        []string
}
