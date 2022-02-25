package hypermedia_controls

type Link struct {
	Href   URI    `json:"href,omitempty"`
	Type   string `json:"type,omitempty"`
	Rel    URI    `json:"rel,omitempty"`
	Anchor any    `json:"anchor,omitempty"`
	Sizes  string `json:"sizes,omitempty"`
}
