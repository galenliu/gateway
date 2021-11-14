package security_scheme

type BasicSecurityScheme struct {
	*securityScheme
	Name string `json:"name"`
	In   string `json:"in"`
}
