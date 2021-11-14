package security_scheme

type BearerSecurityScheme struct {
	*securityScheme
	Authorization string `json:"authorization"`
	Name          string `json:"name"`
	Alg           string `json:"alg"`
	Format        string `json:"format"`
	In            string `json:"in"`
}
