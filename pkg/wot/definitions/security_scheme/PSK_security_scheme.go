package security_scheme

type PSKSecurityScheme struct {
	*securityScheme
	Identity string `json:"identity"`
}
