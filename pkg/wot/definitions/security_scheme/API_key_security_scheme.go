package security_scheme

type APIKeySecurityScheme struct {
	*securityScheme
	Name string `json:"name"`
	In   string `json:"in"`
}
