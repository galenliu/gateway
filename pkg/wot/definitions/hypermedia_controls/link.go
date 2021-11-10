package hypermedia_controls

type Link struct {
	Href   string      `json:"href,omitempty"`
	Type   string      `json:"type,omitempty"`
	Rel    URI         `json:"rel,omitempty"`
	Anchor interface{} `json:"anchor,omitempty"`
	Sizes  string      `json:"sizes,omitempty"`
}
