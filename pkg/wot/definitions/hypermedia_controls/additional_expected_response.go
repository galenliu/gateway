package hypermedia_controls

type AdditionalExpectedResponse struct {
	ContentType string `json:"contentType,omitempty"`
	Success     bool   `json:"success,omitempty"`
	Schema      string `json:"schema,omitempty"`
}
