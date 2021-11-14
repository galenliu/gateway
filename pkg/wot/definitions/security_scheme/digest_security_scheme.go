package security_scheme

type DigestSecurityScheme struct {
	*securityScheme
	Name string `json:"name"`
	In   string `json:"in"`
	Qop  string `json:"qop"`
}
